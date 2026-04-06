package discord

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/config"
	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/bwmarrin/discordgo"
)

const replyTimeout = 30 * time.Second

type Responder interface {
	Reply(context.Context, gateway.Envelope) (string, error)
}

type Bot struct {
	cfg       config.DiscordConfig
	agentName string
	token     string
	responder Responder
	logger    io.Writer

	mu                sync.RWMutex
	botUserID         string
	limitMu           sync.Mutex
	userCooldownUntil map[string]time.Time
	inflight          chan struct{}
	now               func() time.Time
}

func TokenFromEnv(cfg config.DiscordConfig) (string, error) {
	name := strings.TrimSpace(cfg.BotTokenEnv)
	if name == "" {
		return "", errors.New("discord bot token env is not configured")
	}

	token := strings.TrimSpace(os.Getenv(name))
	if token == "" {
		return "", fmt.Errorf("%s is not set", name)
	}

	return token, nil
}

func NewBot(cfg config.DiscordConfig, agentName string, token string, responder Responder, logger io.Writer) (*Bot, error) {
	if strings.TrimSpace(token) == "" {
		return nil, errors.New("discord token is empty")
	}
	if responder == nil {
		return nil, errors.New("discord responder is required")
	}

	resolvedName := strings.TrimSpace(agentName)
	if resolvedName == "" {
		resolvedName = strings.TrimSpace(cfg.DisplayName)
	}
	if resolvedName == "" {
		resolvedName = "Savitar"
	}

	return &Bot{
		cfg:               cfg,
		agentName:         resolvedName,
		token:             token,
		responder:         responder,
		logger:            logger,
		userCooldownUntil: map[string]time.Time{},
		inflight:          makeReplySlots(cfg.MaxConcurrentReplies),
		now:               time.Now,
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	session, err := discordgo.New("Bot " + b.token)
	if err != nil {
		return fmt.Errorf("create discord session: %w", err)
	}

	session.Identify.Intents = b.intents()
	session.ShouldReconnectOnError = true
	session.ShouldRetryOnRateLimit = true
	session.AddHandler(b.onReady)
	session.AddHandler(b.onMessageCreate)

	if err := session.Open(); err != nil {
		return fmt.Errorf("open discord session: %w", err)
	}
	defer func() {
		_ = session.Close()
	}()

	b.logf("discord transport connected")
	<-ctx.Done()
	if errors.Is(ctx.Err(), context.Canceled) {
		return nil
	}

	return ctx.Err()
}

func (b *Bot) onReady(s *discordgo.Session, ready *discordgo.Ready) {
	if ready == nil || ready.User == nil {
		return
	}

	b.setBotUserID(ready.User.ID)
	b.logf("discord ready as %s in %d guilds", ready.User.String(), len(ready.Guilds))

	if text := strings.TrimSpace(b.cfg.PresenceText); text != "" {
		if err := s.UpdateWatchStatus(0, text); err != nil {
			b.logf("failed to update discord presence: %v", err)
		}
	}
}

func (b *Bot) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if !b.shouldRespond(m) {
		return
	}

	release, denial := b.reserveReplyBudget(authorID(m))
	if denial != "" {
		b.sendStatusMessage(s, m, denial)
		return
	}
	defer release()

	_ = s.ChannelTyping(m.ChannelID)
	ackID := ""
	if b.cfg.ImmediateAck {
		ack, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
			Content:         b.ackText(),
			AllowedMentions: allowedMentions(),
			Reference:       m.Reference(),
			Flags:           discordgo.MessageFlagsSuppressNotifications,
		})
		if err != nil {
			b.logf("failed to send discord ack: %v", err)
		} else {
			ackID = ack.ID
		}
	}

	replyCtx, cancel := context.WithTimeout(context.Background(), replyTimeout)
	defer cancel()

	response, err := b.responder.Reply(replyCtx, normalizeEnvelope(m))
	if err != nil {
		b.logf("failed to generate discord reply: %v", err)
		response = fmt.Sprintf("%s hit an internal error while handling that message. Try help or status again.", b.agentName)
	}

	response = clampResponse(response, b.cfg.MaxResponseChars)
	if response == "" {
		response = fmt.Sprintf("%s did not have a useful reply for that. Try help.", b.agentName)
	}

	if ackID != "" {
		edit := discordgo.NewMessageEdit(m.ChannelID, ackID).SetContent(response)
		edit.AllowedMentions = allowedMentions()
		if _, err := s.ChannelMessageEditComplex(edit); err == nil {
			return
		} else {
			b.logf("failed to edit discord ack: %v", err)
		}
	}

	if _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:         response,
		AllowedMentions: allowedMentions(),
		Reference:       m.Reference(),
		Flags:           discordgo.MessageFlagsSuppressNotifications,
	}); err != nil {
		b.logf("failed to send discord reply: %v", err)
	}
}

func (b *Bot) shouldRespond(m *discordgo.MessageCreate) bool {
	if m == nil || m.Author == nil {
		return false
	}
	if m.Author.Bot {
		return false
	}
	if m.Type != discordgo.MessageTypeDefault && m.Type != discordgo.MessageTypeReply {
		return false
	}
	if isDirectMessage(m) {
		return b.cfg.RespondInDirectMessages
	}
	if len(b.cfg.AllowedChannelIDs) > 0 && !containsString(b.cfg.AllowedChannelIDs, m.ChannelID) {
		return false
	}
	if b.cfg.RequireMention && !b.mentionsBot(m) && !b.referencesBot(m) {
		return false
	}
	return true
}

func (b *Bot) intents() discordgo.Intent {
	intents := discordgo.IntentsGuildMessages
	if b.cfg.RespondInDirectMessages {
		intents |= discordgo.IntentsDirectMessages
	}
	if b.cfg.UseMessageContentIntent {
		intents |= discordgo.IntentsMessageContent
	}
	return intents
}

func (b *Bot) mentionsBot(m *discordgo.MessageCreate) bool {
	botUserID := b.currentBotUserID()
	if botUserID == "" {
		return false
	}
	for _, mention := range m.Mentions {
		if mention != nil && mention.ID == botUserID {
			return true
		}
	}
	return false
}

func (b *Bot) referencesBot(m *discordgo.MessageCreate) bool {
	botUserID := b.currentBotUserID()
	if botUserID == "" || m.ReferencedMessage == nil || m.ReferencedMessage.Author == nil {
		return false
	}
	return m.ReferencedMessage.Author.ID == botUserID
}

func (b *Bot) ackText() string {
	return fmt.Sprintf("%s is on it...", b.agentName)
}

func (b *Bot) reserveReplyBudget(senderID string) (func(), string) {
	release := b.acquireReplySlot()
	if release == nil {
		return nil, fmt.Sprintf("%s is handling too many requests right now. Try again shortly.", b.agentName)
	}

	if wait := b.checkAndMarkCooldown(senderID); wait > 0 {
		release()
		return nil, fmt.Sprintf("%s just handled a request from you. Try again in %ds.", b.agentName, int(wait.Seconds())+1)
	}

	return release, ""
}

func (b *Bot) acquireReplySlot() func() {
	if b.inflight == nil {
		return func() {}
	}

	select {
	case b.inflight <- struct{}{}:
		return func() {
			<-b.inflight
		}
	default:
		return nil
	}
}

func (b *Bot) checkAndMarkCooldown(senderID string) time.Duration {
	if b.cfg.PerUserCooldownSeconds <= 0 || senderID == "" {
		return 0
	}

	now := time.Now
	if b.now != nil {
		now = b.now
	}

	b.limitMu.Lock()
	defer b.limitMu.Unlock()
	current := now()
	if until, ok := b.userCooldownUntil[senderID]; ok && current.Before(until) {
		return until.Sub(current)
	}
	b.userCooldownUntil[senderID] = current.Add(time.Duration(b.cfg.PerUserCooldownSeconds) * time.Second)
	return 0
}

func (b *Bot) sendStatusMessage(s *discordgo.Session, m *discordgo.MessageCreate, content string) {
	if s == nil || m == nil || strings.TrimSpace(content) == "" {
		return
	}
	if _, err := s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
		Content:         content,
		AllowedMentions: allowedMentions(),
		Reference:       m.Reference(),
		Flags:           discordgo.MessageFlagsSuppressNotifications,
	}); err != nil {
		b.logf("failed to send discord status reply: %v", err)
	}
}

func makeReplySlots(limit int) chan struct{} {
	if limit <= 0 {
		return nil
	}
	return make(chan struct{}, limit)
}

func (b *Bot) currentBotUserID() string {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.botUserID
}

func (b *Bot) setBotUserID(userID string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.botUserID = userID
}

func (b *Bot) logf(format string, args ...any) {
	if b.logger == nil {
		return
	}
	_, _ = fmt.Fprintf(b.logger, format+"\n", args...)
}

func normalizeEnvelope(m *discordgo.MessageCreate) gateway.Envelope {
	if m == nil || m.Message == nil {
		return gateway.Envelope{Surface: gateway.SurfaceDiscord}
	}

	attachments := make([]gateway.Attachment, 0, len(m.Attachments))
	for _, attachment := range m.Attachments {
		attachments = append(attachments, gateway.Attachment{
			Name:     attachment.Filename,
			MIMEType: attachment.ContentType,
			Path:     attachment.URL,
		})
	}

	mentions := make([]string, 0, len(m.Mentions))
	for _, mention := range m.Mentions {
		if mention != nil {
			mentions = append(mentions, mention.DisplayName())
		}
	}

	replyToID := ""
	if m.ReferencedMessage != nil {
		replyToID = m.ReferencedMessage.ID
	}

	metadata := map[string]string{
		"channelID": m.ChannelID,
		"messageID": m.ID,
	}
	if isDirectMessage(m) {
		metadata["dm"] = "true"
	}
	if m.GuildID != "" {
		metadata["guildID"] = m.GuildID
	}

	return gateway.Envelope{
		Surface:           gateway.SurfaceDiscord,
		ConversationID:    m.ChannelID,
		SenderID:          authorID(m),
		SenderDisplayName: authorDisplayName(m),
		ReplyToID:         replyToID,
		Body:              strings.TrimSpace(m.ContentWithMentionsReplaced()),
		Mentions:          mentions,
		Attachments:       attachments,
		Metadata:          metadata,
	}
}

func authorID(m *discordgo.MessageCreate) string {
	if m == nil || m.Author == nil {
		return ""
	}
	return m.Author.ID
}

func authorDisplayName(m *discordgo.MessageCreate) string {
	if m != nil && m.Member != nil {
		if name := strings.TrimSpace(m.Member.DisplayName()); name != "" {
			return name
		}
	}
	if m != nil && m.Author != nil {
		if name := strings.TrimSpace(m.Author.DisplayName()); name != "" {
			return name
		}
		return m.Author.Username
	}
	return "unknown"
}

func allowedMentions() *discordgo.MessageAllowedMentions {
	return &discordgo.MessageAllowedMentions{
		Parse:       []discordgo.AllowedMentionType{},
		RepliedUser: false,
	}
}

func isDirectMessage(m *discordgo.MessageCreate) bool {
	return m != nil && m.GuildID == ""
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

func clampResponse(content string, limit int) string {
	trimmed := strings.TrimSpace(content)
	if limit <= 0 {
		return trimmed
	}

	runes := []rune(trimmed)
	if len(runes) <= limit {
		return trimmed
	}
	if limit <= 3 {
		return string(runes[:limit])
	}
	return string(runes[:limit-3]) + "..."
}

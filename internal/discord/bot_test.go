package discord

import (
	"strings"
	"testing"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/config"
	"github.com/bwmarrin/discordgo"
)

func newMessage(channelID string, guildID string) *discordgo.MessageCreate {
	return &discordgo.MessageCreate{
		Message: &discordgo.Message{
			ChannelID: channelID,
			GuildID:   guildID,
			Type:      discordgo.MessageTypeDefault,
			Author: &discordgo.User{
				ID:       "user-1",
				Username: "operator",
			},
		},
	}
}

func TestShouldRespondToDirectMessages(t *testing.T) {
	cfg := config.Default().Transports.Discord
	bot := &Bot{cfg: cfg}

	if !bot.shouldRespond(newMessage("dm-channel", "")) {
		t.Fatal("expected direct message to trigger a reply")
	}
}

func TestShouldRequireMentionInGuildChannels(t *testing.T) {
	cfg := config.Default().Transports.Discord
	bot := &Bot{cfg: cfg}
	bot.setBotUserID("bot-1")

	message := newMessage("guild-channel", "guild-1")
	if bot.shouldRespond(message) {
		t.Fatal("expected guild message without a mention to be ignored")
	}

	message.Mentions = []*discordgo.User{{ID: "bot-1", Username: "Savitar"}}
	if !bot.shouldRespond(message) {
		t.Fatal("expected guild message with a mention to trigger a reply")
	}
}

func TestShouldHonorChannelAllowlist(t *testing.T) {
	cfg := config.Default().Transports.Discord
	cfg.AllowedChannelIDs = []string{"allowed-channel"}
	bot := &Bot{cfg: cfg}
	bot.setBotUserID("bot-1")

	message := newMessage("other-channel", "guild-1")
	message.Mentions = []*discordgo.User{{ID: "bot-1", Username: "Savitar"}}
	if bot.shouldRespond(message) {
		t.Fatal("expected message outside the allowlist to be ignored")
	}

	message.ChannelID = "allowed-channel"
	if !bot.shouldRespond(message) {
		t.Fatal("expected allowlisted message to trigger a reply")
	}
}

func TestClampResponseTruncatesLongContent(t *testing.T) {
	if got := clampResponse("abcdef", 5); got != "ab..." {
		t.Fatalf("unexpected truncated response: %q", got)
	}
}

func TestReserveReplyBudgetEnforcesCooldown(t *testing.T) {
	cfg := config.Default().Transports.Discord
	cfg.PerUserCooldownSeconds = 10
	bot := &Bot{
		cfg:               cfg,
		agentName:         "Savitar",
		userCooldownUntil: map[string]time.Time{},
		inflight:          makeReplySlots(1),
		now: func() time.Time {
			return time.Unix(100, 0)
		},
	}

	release, denial := bot.reserveReplyBudget("user-1")
	if denial != "" {
		t.Fatalf("unexpected denial on first request: %q", denial)
	}
	release()

	_, denial = bot.reserveReplyBudget("user-1")
	if !strings.Contains(denial, "Try again in") {
		t.Fatalf("expected cooldown denial, got %q", denial)
	}
}

func TestReserveReplyBudgetHonorsConcurrencyLimit(t *testing.T) {
	cfg := config.Default().Transports.Discord
	cfg.PerUserCooldownSeconds = 0
	bot := &Bot{
		cfg:               cfg,
		agentName:         "Savitar",
		userCooldownUntil: map[string]time.Time{},
		inflight:          makeReplySlots(1),
		now:               time.Now,
	}

	release, denial := bot.reserveReplyBudget("user-1")
	if denial != "" {
		t.Fatalf("unexpected denial on first request: %q", denial)
	}
	defer release()

	_, denial = bot.reserveReplyBudget("user-2")
	if !strings.Contains(denial, "too many requests") {
		t.Fatalf("expected concurrency denial, got %q", denial)
	}
}

package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"text/tabwriter"
	"time"

	"github.com/alexiosbluffmara/savitar/internal/config"
	discordtransport "github.com/alexiosbluffmara/savitar/internal/discord"
	"github.com/alexiosbluffmara/savitar/internal/gateway"
	"github.com/alexiosbluffmara/savitar/internal/memory"
	"github.com/alexiosbluffmara/savitar/internal/models"
	"github.com/alexiosbluffmara/savitar/internal/respond"
	savitarruntime "github.com/alexiosbluffmara/savitar/internal/runtime"
	"github.com/alexiosbluffmara/savitar/internal/session"
	"github.com/alexiosbluffmara/savitar/internal/webui"
)

func Run(stdout io.Writer, stderr io.Writer, version string, args []string) int {
	loaded, err := config.Load(config.DefaultPath())
	if err != nil {
		fmt.Fprintf(stderr, "failed to load config: %v\n", err)
		return 1
	}

	rt := savitarruntime.New(loaded)
	command := "help"
	if len(args) > 0 {
		command = strings.ToLower(args[0])
	}

	switch command {
	case "help", "-h", "--help":
		printHelp(stdout)
		return 0
	case "status":
		return printStatus(stdout, stderr, rt)
	case "doctor":
		return printDoctor(stdout, rt)
	case "agents":
		return printAgents(stdout, stderr, rt)
	case "skills":
		return printSkills(stdout, stderr, rt)
	case "integrations":
		printIntegrations(stdout, rt)
		return 0
	case "gateway":
		printGateway(stdout, rt)
		return 0
	case "persona":
		printPersona(stdout, rt)
		return 0
	case "session":
		return handleSession(stdout, stderr, rt, args[1:])
	case "discord":
		return handleDiscord(stdout, stderr, rt, args[1:])
	case "plan":
		printPlan(stdout, rt)
		return 0
	case "contracts":
		printContracts(stdout, rt)
		return 0
	case "models":
		printModels(stdout, rt)
		return 0
	case "mcp":
		return handleMCP(stdout, stderr, rt, args[1:])
	case "repo":
		return handleRepo(stdout, stderr, rt, args[1:])
	case "memory":
		return handleMemory(stdout, stderr, rt, args[1:])
	case "webui":
		return handleWebUI(stdout, stderr, rt, args[1:])
	case "version":
		fmt.Fprintln(stdout, version)
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command %q\n\n", command)
		printHelp(stderr)
		return 1
	}
}

func printHelp(out io.Writer) {
	fmt.Fprintln(out, "Savitar CLI")
	fmt.Fprintln(out)
	fmt.Fprintln(out, "Usage:")
	fmt.Fprintln(out, "  savitar status")
	fmt.Fprintln(out, "  savitar doctor")
	fmt.Fprintln(out, "  savitar agents")
	fmt.Fprintln(out, "  savitar skills")
	fmt.Fprintln(out, "  savitar integrations")
	fmt.Fprintln(out, "  savitar gateway")
	fmt.Fprintln(out, "  savitar persona")
	fmt.Fprintln(out, "  savitar session [show|init|list]")
	fmt.Fprintln(out, "  savitar discord [status|preview|run]")
	fmt.Fprintln(out, "  savitar mcp [status]")
	fmt.Fprintln(out, "  savitar repo analyze <url>")
	fmt.Fprintln(out, "  savitar memory [list|show <subject> <name>|write <subject> <name> <body>|search <query>|graph <query>]")
	fmt.Fprintln(out, "  savitar webui [serve --demo --addr :8080]")
	fmt.Fprintln(out, "  savitar plan")
	fmt.Fprintln(out, "  savitar contracts")
	fmt.Fprintln(out, "  savitar models")
	fmt.Fprintln(out, "  savitar version")
}

func printStatus(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime) int {
	status, err := rt.Status()
	if err != nil {
		fmt.Fprintf(errOut, "failed to build status: %v\n", err)
		return 1
	}

	fmt.Fprintln(out, "Savitar status")
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "field	value")
	fmt.Fprintf(tw, "root\t%s\n", status.RootDir)
	fmt.Fprintf(tw, "config\t%s\n", status.ConfigPath)
	fmt.Fprintf(tw, "config mode\t%s\n", ternary(status.ConfigLoaded, "loaded", "defaults"))
	fmt.Fprintf(tw, "agent\t%s\n", status.AgentName)
	fmt.Fprintf(tw, "model profiles\t%d\n", status.ModelProfiles)
	fmt.Fprintf(tw, "enabled mcp servers\t%d\n", status.EnabledMCPServers)
	fmt.Fprintf(tw, "workspace agents\t%d\n", status.AgentCount)
	fmt.Fprintf(tw, "workspace skills\t%d\n", status.SkillCount)
	fmt.Fprintf(tw, "enabled public surfaces\t%d/%d\n", status.EnabledSurfaces, status.TotalSurfaces)
	fmt.Fprintf(tw, "session\t%s\n", ternary(status.SessionInitialized, status.SessionPath, status.SessionPath+" (not initialized)"))
	_ = tw.Flush()

	return 0
}

func printDoctor(out io.Writer, rt *savitarruntime.Runtime) int {
	loaded := rt.Config()
	if loaded.Exists {
		fmt.Fprintf(out, "config\tloaded\t%s\n", loaded.Path)
	} else {
		fmt.Fprintf(out, "config\tdefaults\t%s not found\n", loaded.Path)
	}

	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "status\tname\tdetails")

	missingRequired := 0
	for _, check := range rt.Doctor() {
		status := "ok"
		if !check.Available && check.Required {
			status = "missing"
			missingRequired++
		} else if !check.Available {
			status = "optional-missing"
		}

		fmt.Fprintf(tw, "%s\t%s\t%s\n", status, check.Name, check.Details)
	}

	_ = tw.Flush()

	if missingRequired > 0 {
		return 1
	}

	return 0
}

func printPlan(out io.Writer, rt *savitarruntime.Runtime) {
	for index, step := range rt.Plan() {
		fmt.Fprintf(out, "%d. %s\n", index+1, step)
	}
}

func printAgents(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime) int {
	agents, err := rt.Agents()
	if err != nil {
		fmt.Fprintf(errOut, "failed to discover agents: %v\n", err)
		return 1
	}

	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "name	invocable	tools	description")
	for _, agent := range agents {
		fmt.Fprintf(tw, "%s	%t	%s	%s\n", agent.Name, agent.UserInvocable, strings.Join(agent.Tools, ", "), agent.Description)
	}
	_ = tw.Flush()

	return 0
}

func printSkills(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime) int {
	skills, err := rt.Skills()
	if err != nil {
		fmt.Fprintf(errOut, "failed to discover skills: %v\n", err)
		return 1
	}

	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "name	argument hint	description")
	for _, skill := range skills {
		fmt.Fprintf(tw, "%s	%s	%s\n", skill.Name, skill.ArgumentHint, skill.Description)
	}
	_ = tw.Flush()

	return 0
}

func printIntegrations(out io.Writer, rt *savitarruntime.Runtime) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "name\tenabled\tauth\tcredential\tenv\tcli\tdetails")
	for _, status := range rt.IntegrationStatuses() {
		credential := ternary(status.CredentialPresent, "ready", "missing")
		if status.AuthSource == "local-http" || status.AuthSource == "none" {
			credential = "n/a"
		}
		fmt.Fprintf(
			tw,
			"%s\t%t\t%s\t%s\t%s\t%s\t%s\n",
			status.Name,
			status.Enabled,
			status.AuthSource,
			credential,
			blankIfEmpty(status.TokenEnv),
			blankIfEmpty(status.CLIPath),
			blankIfEmpty(status.Details),
		)
	}
	_ = tw.Flush()
}

func printGateway(out io.Writer, rt *savitarruntime.Runtime) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "surface	enabled	kind	identity	details")
	for _, surface := range rt.GatewayPlan().Surfaces {
		fmt.Fprintf(tw, "%s	%t	%s	%s	%s\n", surface.Name, surface.Enabled, surface.Kind, surface.Identity, surface.Details)
	}
	_ = tw.Flush()
}

func printPersona(out io.Writer, rt *savitarruntime.Runtime) {
	profile := rt.Persona()
	fmt.Fprintln(out, profile.Name)
	fmt.Fprintf(out, "style: %s\n", profile.Style)
	fmt.Fprintf(out, "tone: %s\n", profile.Tone)
	fmt.Fprintf(out, "design bias: %s\n", profile.DesignBias)
	fmt.Fprintf(out, "commentary density: %s\n", profile.CommentaryDensity)
	fmt.Fprintf(out, "public bio: %s\n", profile.PublicBio)
	fmt.Fprintf(out, "disclosure policy: %s\n", profile.DisclosurePolicy)
}

func handleSession(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, args []string) int {
	action := "show"
	if len(args) > 0 {
		action = strings.ToLower(args[0])
	}

	switch action {
	case "show", "status":
		report, err := rt.Session()
		if err != nil {
			fmt.Fprintf(errOut, "failed to read session state: %v\n", err)
			return 1
		}

		printSession(out, report)
		return 0
	case "init":
		report, err := rt.InitSession()
		if err != nil {
			fmt.Fprintf(errOut, "failed to initialize session state: %v\n", err)
			return 1
		}

		printSession(out, report)
		return 0
	case "list":
		summaries, err := rt.SessionList()
		if err != nil {
			fmt.Fprintf(errOut, "failed to list sessions: %v\n", err)
			return 1
		}
		if len(summaries) == 0 {
			fmt.Fprintln(out, "no session files found")
			return 0
		}
		tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "name\tsurface\tupdated\tpath")
		for _, s := range summaries {
			updated := "<unknown>"
			if !s.UpdatedAt.IsZero() {
				updated = s.UpdatedAt.Format(time.RFC3339)
			}
			fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n", s.Name, blankIfEmpty(s.Surface), updated, s.Path)
		}
		_ = tw.Flush()
		return 0
	default:
		fmt.Fprintf(errOut, "unknown session action %q\n", action)
		fmt.Fprintln(errOut, "usage: savitar session [show|init|list]")
		return 1
	}
}

func handleDiscord(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, args []string) int {
	action := "status"
	if len(args) > 0 {
		action = strings.ToLower(args[0])
	}

	switch action {
	case "status":
		printDiscordStatus(out, rt)
		return 0
	case "preview":
		if len(args) < 2 {
			fmt.Fprintln(errOut, "usage: savitar discord preview <message>")
			return 1
		}

		return previewDiscordReply(out, errOut, rt, strings.Join(args[1:], " "))
	case "run":
		return runDiscord(out, errOut, rt)
	default:
		fmt.Fprintf(errOut, "unknown discord action %q\n", action)
		fmt.Fprintln(errOut, "usage: savitar discord [status|preview|run]")
		return 1
	}
}

func handleWebUI(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, args []string) int {
	action := "serve"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		action = strings.ToLower(args[0])
		args = args[1:]
	}
	if action != "serve" {
		fmt.Fprintf(errOut, "unknown webui action %q\n", action)
		fmt.Fprintln(errOut, "usage: savitar webui [serve --demo --addr :8080]")
		return 1
	}

	demo := false
	addr := ":8080"
	for index := 0; index < len(args); index++ {
		switch args[index] {
		case "--demo":
			demo = true
		case "--addr":
			if index+1 >= len(args) {
				fmt.Fprintln(errOut, "usage: savitar webui [serve --demo --addr :8080]")
				return 1
			}
			index++
			addr = args[index]
		default:
			fmt.Fprintln(errOut, "usage: savitar webui [serve --demo --addr :8080]")
			return 1
		}
	}
	if !demo {
		fmt.Fprintln(errOut, "web UI auth is not implemented yet; use savitar webui serve --demo for the local prototype")
		return 1
	}

	return webui.RunDemo(out, errOut, rt, addr)
}

func printDiscordStatus(out io.Writer, rt *savitarruntime.Runtime) {
	status := rt.DiscordStatus()
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "field\tvalue")
	fmt.Fprintf(tw, "enabled\t%t\n", status.Enabled)
	fmt.Fprintf(tw, "token env\t%s\n", status.TokenEnv)
	fmt.Fprintf(tw, "token present\t%t\n", status.TokenPresent)
	fmt.Fprintf(tw, "display name\t%s\n", blankIfEmpty(status.DisplayName))
	fmt.Fprintf(tw, "operator users\t%d\n", status.OperatorUserCount)
	fmt.Fprintf(tw, "trigger mode\t%s\n", status.TriggerMode)
	fmt.Fprintf(tw, "live web guild channels\t%t\n", status.AllowLiveWebLookupInGuilds)
	fmt.Fprintf(tw, "live web DMs\t%t\n", status.AllowLiveWebLookupInDMs)
	fmt.Fprintf(tw, "message content intent\t%t\n", status.UseMessageContentIntent)
	fmt.Fprintf(tw, "allowed channels\t%s\n", blankIfEmpty(strings.Join(status.AllowedChannelIDs, ", ")))
	fmt.Fprintf(tw, "per-user cooldown\t%ds\n", status.PerUserCooldownSeconds)
	fmt.Fprintf(tw, "max concurrent replies\t%d\n", status.MaxConcurrentReplies)
	fmt.Fprintf(tw, "max response chars\t%d\n", status.MaxResponseChars)
	fmt.Fprintf(tw, "presence\t%s\n", blankIfEmpty(status.PresenceText))
	fmt.Fprintf(tw, "immediate ack\t%t\n", status.ImmediateAck)
	_ = tw.Flush()

	if len(status.AllowedChannelIDs) == 0 {
		fmt.Fprintln(out, "note: no channel allowlist is configured")
	}
	if !status.RequireMention && !status.UseMessageContentIntent {
		fmt.Fprintln(out, "warning: guild-wide replies require useMessageContentIntent=true and the privileged intent enabled in the Discord developer portal")
	}
	if !status.RequireMention && len(status.AllowedChannelIDs) == 0 {
		fmt.Fprintln(out, "warning: disabling requireMention without allowedChannelIDs would capture all joined guild messages")
	}
	fmt.Fprintln(out, "note: conversational replies stay local-only and fail closed unless a local Ollama target is configured")
	if status.RespondInDirectMessages && !status.AllowLiveWebLookupInDMs {
		fmt.Fprintln(out, "note: direct-message live web lookup is disabled unless explicitly opted in locally")
	}
	if !status.AllowLiveWebLookupInGuilds {
		fmt.Fprintln(out, "note: guild live web lookup is disabled unless explicitly opted in locally")
	}
}

func previewDiscordReply(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, body string) int {
	engine := respond.NewEngine(rt)
	reply, err := engine.Reply(context.Background(), gateway.Envelope{
		Surface:           gateway.SurfaceDiscord,
		ConversationID:    "preview",
		SenderID:          "preview-user",
		SenderDisplayName: "preview-user",
		Body:              body,
		Metadata: map[string]string{
			"mode": "preview",
		},
	})
	if err != nil {
		fmt.Fprintf(errOut, "failed to generate preview reply: %v\n", err)
		return 1
	}

	fmt.Fprintln(out, reply)
	return 0
}

func runDiscord(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime) int {
	cfg := rt.Config().Config.Transports.Discord
	if cfg.PerUserCooldownSeconds <= 0 {
		fmt.Fprintln(errOut, "discord transport requires perUserCooldownSeconds to be greater than zero")
		return 1
	}
	if cfg.MaxConcurrentReplies <= 0 {
		fmt.Fprintln(errOut, "discord transport requires maxConcurrentReplies to be greater than zero")
		return 1
	}
	if !cfg.RequireMention && !cfg.UseMessageContentIntent {
		fmt.Fprintln(errOut, "discord transport requires useMessageContentIntent=true when requireMention=false")
		return 1
	}
	if !cfg.RequireMention && len(cfg.AllowedChannelIDs) == 0 {
		fmt.Fprintln(errOut, "discord transport requires allowedChannelIDs to be set when requireMention=false")
		return 1
	}

	token, err := discordtransport.TokenFromEnv(cfg)
	if err != nil {
		fmt.Fprintf(errOut, "failed to load discord bot token: %v\n", err)
		return 1
	}

	if !cfg.Enabled {
		fmt.Fprintln(errOut, "warning: discord transport is disabled in config; running because the command was invoked explicitly")
	}

	bot, err := discordtransport.NewBot(cfg, rt.Persona().Name, token, respond.NewEngine(rt), errOut)
	if err != nil {
		fmt.Fprintf(errOut, "failed to create discord transport: %v\n", err)
		return 1
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	fmt.Fprintln(out, "starting Discord transport; press Ctrl+C to stop")
	if err := bot.Run(ctx); err != nil {
		fmt.Fprintf(errOut, "discord transport exited with error: %v\n", err)
		return 1
	}
	fmt.Fprintln(out, "Discord transport stopped")
	return 0
}

func printSession(out io.Writer, report session.Report) {
	fmt.Fprintln(out, "Savitar session")
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "field	value")
	fmt.Fprintf(tw, "path	%s\n", report.Path)
	fmt.Fprintf(tw, "initialized	%t\n", report.Exists)
	fmt.Fprintf(tw, "surface	%s\n", report.State.CurrentSurface)
	fmt.Fprintf(tw, "model profile	%s\n", report.State.CurrentModelProfile)
	fmt.Fprintf(tw, "conversation	%s\n", blankIfEmpty(report.State.ActiveConversationID))
	fmt.Fprintf(tw, "last command	%s\n", blankIfEmpty(report.State.LastCommand))
	_ = tw.Flush()
}

func ternary(condition bool, left string, right string) string {
	if condition {
		return left
	}

	return right
}

func blankIfEmpty(value string) string {
	if value == "" {
		return "<none>"
	}

	return value
}

func printContracts(out io.Writer, rt *savitarruntime.Runtime) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "name\tstage\tgoal")
	for _, contract := range rt.Contracts() {
		fmt.Fprintf(tw, "%s\t%s\t%s\n", contract.Name, contract.Stage, contract.Goal)
	}
	_ = tw.Flush()
}

func handleMCP(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, args []string) int {
	action := "status"
	if len(args) > 0 {
		action = strings.ToLower(args[0])
	}
	switch action {
	case "status":
		return printMCPStatus(out, errOut, rt)
	default:
		fmt.Fprintf(errOut, "unknown mcp action %q\n", action)
		fmt.Fprintln(errOut, "usage: savitar mcp [status]")
		return 1
	}
}

func printMCPStatus(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime) int {
	fmt.Fprintln(out, "probing MCP servers (15s timeout per server)...")
	statuses := rt.MCPStatus(15 * time.Second)
	if len(statuses) == 0 {
		fmt.Fprintln(out, "no MCP servers configured (check .vscode/mcp.json)")
		return 0
	}
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "name\tmode\tenabled\treachable\ttools\tserver\terror")
	for _, s := range statuses {
		toolCount := fmt.Sprintf("%d", len(s.Tools))
		serverName := s.Info.Name
		errText := s.Err
		fmt.Fprintf(tw, "%s\t%s\t%t\t%t\t%s\t%s\t%s\n",
			s.Name, s.Mode, s.Enabled, s.Reachable,
			toolCount, blankIfEmpty(serverName), blankIfEmpty(errText))
	}
	_ = tw.Flush()

	// Print tool lists for reachable servers.
	for _, s := range statuses {
		if !s.Reachable || len(s.Tools) == 0 {
			continue
		}
		fmt.Fprintf(out, "\n%s tools:\n", s.Name)
		for _, t := range s.Tools {
			desc := t.Description
			if len(desc) > 80 {
				desc = desc[:77] + "..."
			}
			fmt.Fprintf(out, "  - %s: %s\n", t.Name, desc)
		}
	}
	_ = errOut
	return 0
}

func handleRepo(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, args []string) int {
	action := ""
	if len(args) > 0 {
		action = strings.ToLower(args[0])
	}
	switch action {
	case "analyze":
		if len(args) < 2 {
			fmt.Fprintln(errOut, "usage: savitar repo analyze <url>")
			return 1
		}
		repoURL := args[1]
		fmt.Fprintf(out, "analyzing %s...\n", repoURL)
		summary, err := rt.RepoAnalyze(context.Background(), repoURL)
		if err != nil {
			fmt.Fprintf(errOut, "repo analysis failed: %v\n", err)
			return 1
		}
		tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "field\tvalue")
		fmt.Fprintf(tw, "url\t%s\n", summary.URL)
		fmt.Fprintf(tw, "branch\t%s\n", summary.Branch)
		fmt.Fprintf(tw, "commit\t%s\n", summary.CommitHash)
		fmt.Fprintf(tw, "commits\t%d\n", summary.Stats.CommitCount)
		fmt.Fprintf(tw, "files\t%d\n", summary.Stats.FileCount)
		fmt.Fprintf(tw, "analyzed at\t%s\n", summary.Provenance.AnalyzedAt)
		_ = tw.Flush()
		fmt.Fprintln(out, "\nrecent commits:")
		for _, line := range strings.Split(summary.LogSnippet, "\n")[:min(10, len(strings.Split(summary.LogSnippet, "\n")))] {
			if line != "" {
				fmt.Fprintf(out, "  %s\n", line)
			}
		}
		return 0
	default:
		fmt.Fprintf(errOut, "unknown repo action %q\n", action)
		fmt.Fprintln(errOut, "usage: savitar repo analyze <url>")
		return 1
	}
}

func handleMemory(out io.Writer, errOut io.Writer, rt *savitarruntime.Runtime, args []string) int {
	action := "list"
	if len(args) > 0 {
		action = strings.ToLower(args[0])
	}
	switch action {
	case "list":
		metas, err := rt.MemoryList()
		if err != nil {
			fmt.Fprintf(errOut, "failed to list memory packs: %v\n", err)
			return 1
		}
		if len(metas) == 0 {
			fmt.Fprintln(out, "no memory packs found")
			return 0
		}
		tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
		fmt.Fprintln(tw, "name\tsubject\tupdated\tsize")
		for _, m := range metas {
			fmt.Fprintf(tw, "%s\t%s\t%s\t%d bytes\n",
				m.Name, m.Subject, m.UpdatedAt.Format(time.RFC3339), m.SizeBytes)
		}
		_ = tw.Flush()
		return 0
	case "show":
		if len(args) < 3 {
			fmt.Fprintln(errOut, "usage: savitar memory show <subject> <name>")
			return 1
		}
		subject, name := args[1], args[2]
		pack, err := rt.MemoryStore().Get(subject, name)
		if err != nil {
			fmt.Fprintf(errOut, "failed to read memory pack: %v\n", err)
			return 1
		}
		fmt.Fprintf(out, "name: %s\n", pack.Name)
		fmt.Fprintf(out, "subject: %s\n", pack.Subject)
		fmt.Fprintf(out, "updated: %s\n", pack.UpdatedAt.Format(time.RFC3339))
		fmt.Fprintln(out)
		fmt.Fprintln(out, pack.Body)
		return 0
	case "write":
		// savitar memory write <subject> <name> <body> — for testing from the CLI
		if len(args) < 4 {
			fmt.Fprintln(errOut, "usage: savitar memory write <subject> <name> <body>")
			return 1
		}
		pack := memory.New(args[2], args[1], strings.Join(args[3:], " "))
		if err := rt.MemoryWrite(pack); err != nil {
			fmt.Fprintf(errOut, "failed to write memory pack: %v\n", err)
			return 1
		}
		fmt.Fprintf(out, "wrote pack %q to subject %q\n", pack.Name, pack.Subject)
		return 0
	case "search":
		if len(args) < 2 {
			fmt.Fprintln(errOut, "usage: savitar memory search <query>")
			return 1
		}
		result, err := rt.RepoMarkdownSearch(strings.Join(args[1:], " "))
		if err != nil {
			fmt.Fprintf(errOut, "failed to search repo markdown: %v\n", err)
			return 1
		}
		if len(result.Results) == 0 {
			fmt.Fprintln(out, "no repo markdown results found")
			return 0
		}
		for _, match := range result.Results {
			fmt.Fprintf(out, "%s\n", match.SourceLabel())
			fmt.Fprintf(out, "score: %d\n", match.Score)
			fmt.Fprintf(out, "%s\n\n", match.Snippet)
		}
		return 0
	case "graph":
		if len(args) < 2 {
			fmt.Fprintln(errOut, "usage: savitar memory graph <query>")
			return 1
		}
		result, err := rt.RepoMarkdownSearch(strings.Join(args[1:], " "))
		if err != nil {
			fmt.Fprintf(errOut, "failed to build repo markdown graph: %v\n", err)
			return 1
		}
		graph := result.GraphString()
		if graph == "" {
			fmt.Fprintln(out, "no repo markdown graph found")
			return 0
		}
		fmt.Fprintln(out, graph)
		return 0
	default:
		fmt.Fprintf(errOut, "unknown memory action %q\n", action)
		fmt.Fprintln(errOut, "usage: savitar memory [list|show <subject> <name>|write <subject> <name> <body>|search <query>|graph <query>]")
		return 1
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printModels(out io.Writer, rt *savitarruntime.Runtime) {
	tw := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "name\tprovider\tmodel\tusage\tpurpose")
	for _, profile := range rt.ModelProfiles() {
		fmt.Fprintf(tw, "%s\t%s\t%s\t%.2fx\t%s\n", profile.Name, profile.Provider, profile.Model, profile.UsageMultiplier, profile.Purpose)
	}
	_ = tw.Flush()

	decision := rt.Router().Route(models.Task{Complexity: models.ComplexityComplex})
	fmt.Fprintf(out, "\nexample complex-task route: %s (%s)\n", decision.Profile.Name, decision.Reason)
}

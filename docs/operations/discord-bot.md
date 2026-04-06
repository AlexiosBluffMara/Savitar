# Discord Bot Runbook

This runbook covers the current Savitar Discord surface.

## Required now

- A Discord application with a bot user.
- A local bot token exposed as `SAVITAR_DISCORD_BOT_TOKEN`.
- A local `config/savitar.local.json` override for Discord behavior.

## Recommended local config

Use the Discord block below inside `config/savitar.local.json`.

```json
{
  "transports": {
    "discord": {
      "enabled": true,
      "botTokenEnv": "SAVITAR_DISCORD_BOT_TOKEN",
      "displayName": "Savitar",
      "operatorUserIDs": ["123456789012345678"],
      "requireMention": true,
      "respondInDirectMessages": true,
      "allowCloudRepliesInGuildChannels": false,
      "allowCloudRepliesInDirectMessages": false,
      "useMessageContentIntent": false,
      "allowedChannelIDs": ["123456789012345678"],
      "perUserCooldownSeconds": 5,
      "maxConcurrentReplies": 2,
      "maxResponseChars": 1800,
      "presenceText": "the repo",
      "immediateAck": true
    }
  }
}
```

Keep real channel IDs local. Do not commit tokens or guild-specific identifiers.

`operatorUserIDs` controls which Discord user IDs can run operator diagnostics like `status`, `doctor`, `gateway`, and `session`. If the list is empty, remote users only get the public-safe command set.

## Intent policy

- Default mode: mention-or-DM. This keeps the transport usable without channel-wide message capture.
- Cloud-backed replies for guild channels are disabled by default. Enable `allowCloudRepliesInGuildChannels` locally only if you explicitly want normal guild traffic to leave the host.
- Cloud-backed replies for direct messages are disabled by default. Enable `allowCloudRepliesInDirectMessages` locally only if you explicitly want DM content to leave the host.
- If `requireMention` is `false`, you should also set `useMessageContentIntent` to `true` locally.
- If `requireMention` is `false`, you should keep `allowedChannelIDs` non-empty so the bot does not capture every joined guild channel.
- When `useMessageContentIntent` is `true`, enable the privileged `MESSAGE_CONTENT` intent in the Discord developer portal.
- If the gateway closes with `4014`, the app is asking for a privileged intent that is not enabled or approved.
- `perUserCooldownSeconds` and `maxConcurrentReplies` are the first backpressure controls for the public reply path.
- `savitar discord run` rejects zero or negative values for those backpressure settings.

## Operator commands

- `savitar discord status` shows token presence, trigger mode, allowlist state, and response limits.
- `savitar discord preview status` shows the deterministic operator reply output without connecting to Discord.
- `savitar discord preview <message>` exercises the conversational Ollama-backed reply path without connecting to Discord.
- `savitar discord run` starts the live bot and keeps running until interrupted.

## Current behavior

- Ignores bot-authored messages.
- Replies to DMs by default.
- Requires a mention or reply in guild channels by default.
- Keeps `help` and `ping` public-safe, while operator diagnostics require `operatorUserIDs` membership.
- Sends normal non-command messages through the Ollama-backed reply path when a reviewed Ollama target is configured.
- Guild-channel cloud replies require explicit local opt-in when no local Ollama target is available.
- Keeps a short per-user-per-channel conversation history in memory for follow-up turns during the current bot process.
- Falls back to a public-safe unavailable message if the model-backed path is not configured or fails during live handling.
- Applies a per-user cooldown and a small concurrent reply cap before starting model work.
- Sends an immediate placeholder reply, then edits it with the final response.
- Supports public-safe replies for `help` and `ping`.
- Supports operator-only replies for `status`, `doctor`, `plan`, `models`, `contracts`, `gateway`, `persona`, and `session`.

## Troubleshooting

- Missing token: export `SAVITAR_DISCORD_BOT_TOKEN` before running the bot.
- Model reply preview failed: confirm `integrations.ollama.enabled` is true in local config and `OLLAMA_API_KEY` is loaded into the reviewed shell.
- Denied diagnostics: add the sender to `operatorUserIDs` if remote operator access is intended.
- Empty guild message content: keep mention mode on, or enable the privileged message-content intent.
- Guild-wide capture blocked on startup: set `allowedChannelIDs` before disabling `requireMention`.
- No reply in a channel: confirm the channel is in `allowedChannelIDs` if you use a channel allowlist.
- Oversized replies: lower `maxResponseChars` in local config if your target channel needs shorter responses.
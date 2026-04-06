#!/usr/bin/env node

const chunks = [];

for await (const chunk of process.stdin) {
  chunks.push(chunk);
}

const payload = Buffer.concat(chunks).toString("utf8");
const blockedPatterns = [
  /git\s+reset\s+--hard/i,
  /git\s+checkout\s+--/i,
  /\brm\s+-rf\b/i,
  /\bsudo\b/i,
  /\bdiskutil\b/i,
  /\bshutdown\b/i,
  /\breboot\b/i
];

const blockedPattern = blockedPatterns.find((pattern) => pattern.test(payload));

if (blockedPattern) {
  process.stdout.write(
    JSON.stringify({
      hookSpecificOutput: {
        hookEventName: "PreToolUse",
        permissionDecision: "deny",
        permissionDecisionReason: "Blocked by Savitar workspace safety hook"
      },
      systemMessage: "Use a non-destructive alternative or get explicit approval before retrying."
    })
  );
  process.exit(0);
}

process.stdout.write(
  JSON.stringify({
    hookSpecificOutput: {
      hookEventName: "PreToolUse",
      permissionDecision: "allow"
    }
  })
);
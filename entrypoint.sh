#!/bin/sh
set -eu

# If no explicit config JSON was provided, generate a default stdio config
# that points the bot to the embedded MCP server binary.
if [ -z "${MCP_CONFIG_JSON:-}" ] && [ -n "${DISCORD_BOT_TOKEN:-}" ]; then
  MCP_CONFIG_JSON=$(cat <<EOF
{"discord":{"token":"${DISCORD_BOT_TOKEN}","guildId":"${DISCORD_GUILD_ID:-}"},"mcp":{"command":"/app/mcp-server","args":[],"env":{},"transport":"stdio"}}
EOF
)
  export MCP_CONFIG_JSON
fi

mcp_pid=""
bot_pid=""

if [ "${START_MCP_SIDE_PROCESS:-false}" = "true" ]; then
  /app/mcp-server ${MCP_SERVER_ARGS:-} &
  mcp_pid=$!
fi

/app/mcpdiscord "$@" &
bot_pid=$!

if [ -z "$mcp_pid" ]; then
  wait "$bot_pid"
  exit $?
fi

# If either process dies, exit non-zero so orchestrators can restart the container.
while true; do
  if ! kill -0 "$bot_pid" 2>/dev/null; then
    wait "$bot_pid" || true
    kill "$mcp_pid" 2>/dev/null || true
    wait "$mcp_pid" 2>/dev/null || true
    exit 1
  fi

  if ! kill -0 "$mcp_pid" 2>/dev/null; then
    wait "$mcp_pid" || true
    kill "$bot_pid" 2>/dev/null || true
    wait "$bot_pid" 2>/dev/null || true
    exit 1
  fi

  sleep 1
done

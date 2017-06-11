#!/bin/zsh

__history::history::add()
{
    local status_code="$status"
    local last_command="$(fc -ln -1)"

    command history add \
        --command "$last_command" \
        --dir     "$PWD" \
        --status  "$status_code" \
        --branch  "$(git rev-parse --abbrev-ref HEAD 2>/dev/null)"
}

__history::history::sync()
{
    local status_code
    local before after
    before=$SECONDS
    command history sync \
        --ask \
        --diff=100 \
        --interval=${ZSH_HISTORY_AUTO_SYNC_INTERVAL:-"1h"} \
        2>/dev/null
    status_code=$status
    after=$SECONDS
    if (( $status_code == 0 && (after - before) > 1 )); then
        printf "[$(date)] Synced successfully\n"
    fi
}

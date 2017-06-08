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
    command history sync --interval=1h 2>/dev/null
    if (( $status == 0 )); then
        echo "$(date): synced successfully ~ $(command history config --get history.path)"
    fi
}

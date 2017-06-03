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

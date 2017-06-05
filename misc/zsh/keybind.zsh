#!/bin/zsh

__history::keybind::get_by_dir()
{
    local buf
    buf="$(command history search --dir --branch --query "$LBUFFER" 2>/dev/null)"
    if [[ -n $buf ]]; then
        BUFFER="$buf"
        CURSOR=$#BUFFER
        zle reset-prompt
    fi
}

__history::keybind::get_all()
{
    local buf
    buf="$(command history search --query "$LBUFFER" 2>/dev/null)"
    if [[ -n $buf ]]; then
        BUFFER="$buf"
        CURSOR=$#BUFFER
        zle reset-prompt
    fi
}

__history::keybind::arrow_up()
{
    __history::substring::search_begin
    __history::substring::history_up
    __history::substring::search_end
}

__history::keybind::arrow_down()
{
    __history::substring::search_begin
    __history::substring::history_down
    __history::substring::search_end
}

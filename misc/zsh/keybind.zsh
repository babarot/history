#!/bin/zsh

__history::keybind::get_by_dir()
{
    local buf
    buf="$(command history search --dir --branch --query "$LBUFFER")"
    if [[ -n $buf ]]; then
        BUFFER="$buf"
        CURSOR=$#BUFFER
        zle reset-prompt
    fi
}

__history::keybind::get_all()
{
    local buf col opt
    col="$(command history config --get "history.record.columns" | sed 's/\[//;s/\]//;s/ /,/g')"
    if [[ ! $col =~ "{{.Base}}" ]]; then
        opt="--columns $col,{{.Base}}"
    fi
    buf="$(command history search $opt --query "$LBUFFER")"
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

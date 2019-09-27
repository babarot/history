#!/bin/zsh

__history::keybind::get()
{
    local buf opt
    # by default, equals to __history::keybind::get_by_dir behavior
    cmd="command history search $ZSH_HISTORY_FILTER_OPTIONS"
    if [[ -n "$LBUFFER" ]]; then
        cmd="$cmd --query "$LBUFFER""
    fi
    buf="$(eval $cmd)"
    if [[ -n $buf ]]; then
        BUFFER="$buf"
        CURSOR=$#BUFFER
    fi
    zle reset-prompt
}

__history::keybind::get_by_dir()
{
    local buf
    cmd="command history search $ZSH_HISTORY_FILTER_OPTIONS_BY_DIR"
    if [[ -n "$LBUFFER" ]]; then
        cmd="$cmd --query "$LBUFFER""
    fi
    buf="$(eval $cmd)"
    if [[ -n $buf ]]; then
        BUFFER="$buf"
        CURSOR=$#BUFFER
    fi
    zle reset-prompt
}

__history::keybind::get_all()
{
    local buf opt
    if [[ -n $ZSH_HISTORY_COLUMNS_GET_ALL ]]; then
        opt="--columns $ZSH_HISTORY_COLUMNS_GET_ALL"
    fi
    buf="$(command history search $opt --query "$LBUFFER")"
    if [[ -n $buf ]]; then
        BUFFER="$buf"
        CURSOR=$#BUFFER
    fi
    zle reset-prompt
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

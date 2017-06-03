#!/bin/zsh

__history::keybind::get_by_dir()
{
    BUFFER="$(command history search --dir --branch 2>/dev/null)"
    CURSOR=$#BUFFER
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

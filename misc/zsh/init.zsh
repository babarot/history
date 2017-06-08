#!/bin/zsh

ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_FOUND="bg=magenta,fg=white,bold"
ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_NOT_FOUND="bg=red,fg=white,bold"
ZSH_HISTORY_SUBSTRING_SEARCH_GLOBBING_FLAGS="i"

#
# Keybindings
#

if [[ -n $ZSH_HISTORY_KEYBIND_GET_BY_DIR ]]; then
    zle -N "__history::keybind::get_by_dir"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_BY_DIR" "__history::keybind::get_by_dir"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_GET_ALL ]]; then
    zle -N "__history::keybind::get_all"
    bindkey "$ZSH_HISTORY_KEYBIND_GET_ALL" "__history::keybind::get_all"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_ARROW_UP ]]; then
    zle -N "__history::keybind::arrow_up"
    bindkey "$ZSH_HISTORY_KEYBIND_ARROW_UP" "__history::keybind::arrow_up"
fi

if [[ -n $ZSH_HISTORY_KEYBIND_ARROW_DOWN ]]; then
    zle -N "__history::keybind::arrow_down"
    bindkey "$ZSH_HISTORY_KEYBIND_ARROW_DOWN" "__history::keybind::arrow_down"
fi

if [[ -z $ZSH_HISTORY_COLUMNS_GET_ALL ]]; then
    export ZSH_HISTORY_COLUMNS_GET_ALL="{{.Time}},{{.Status}},{{.Command}},({{.Base}})"
fi

#
# Configurations
#

if [[ $ZSH_HISTORY_CASE_SENSITIVE == true ]]; then
    unset ZSH_HISTORY_SUBSTRING_SEARCH_GLOBBING_FLAGS
fi

if [[ $ZSH_HISTORY_DISABLE_COLOR == true ]]; then
    unset ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_FOUND
    unset ZSH_HISTORY_SUBSTRING_SEARCH_HIGHLIGHT_NOT_FOUND
fi

#
# Loading
#

for f in "${0:A:h}"/*.zsh(N-.)
do
    source "$f" 2>/dev/null
done
unset f

autoload -Uz add-zsh-hook
add-zsh-hook precmd  "__history::history::add"
add-zsh-hook preexec "__history::substring::reset"

if [[ ${ZSH_HISTORY_AUTO_SYNC:-true} == true ]]; then
    add-zsh-hook precmd  "__history::history::sync"
fi

#
# Configurations
#

if test -z "$fish_history_cmd_name"
    set -g fish_history_cmd_name history
end

if test -z "$fish_history_auto_sync"
    set -g fish_history_auto_sync true
end

if test -z "$fish_history_auto_sync_interval"
    set -g fish_history_auto_sync_interval 1h
end

if test -z "$fish_history_columns_get_all"
    set -g fish_history_columns_get_all "{{.Time}}, {{.Status}},({{.Base}}:{{.Branch}})"
end

if test -z "$fish_history_filter_options"
    set -g fish_history_filter_options "--filter-dir --filter-branch" 
end

#
# Alias
#

function $fish_history_cmd_name -d "enhanced history for your shell"
    command history $argv
end

#
# Completions
#

## erase old completions
complete -ec $fish_history_cmd_name

## subcommands
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'add' -d 'Add new history'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'config' -d 'Config the setting file'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'edit' -d 'Edit your history file directly'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'list' -d 'List the history'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'search' -d 'Search the command from the history file'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'delete' -d 'Delete the record from history file'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'sync' -d 'Sync the history file with gist'
complete -xc $fish_history_cmd_name -n '__fish_use_subcommand' -a 'help' -d 'Show help for any command'

## global options
complete -xc $fish_history_cmd_name -n '__fish_no_arguments' -s h -l help -d 'Show the help message'
complete -xc $fish_history_cmd_name -n '__fish_no_arguments' -s v -l version -d 'Show the version and exit'

## options for add
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from add' -s h -l help -d 'Show the help message'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from add' -l branch -d 'Set branch'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from add' -l command -d 'Set command'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from add' -l dir -d 'Set dir'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from add' -l status -d 'Set status'

## options for search/list/delete
for cmd in search list delete
    set -l Cmd (string sub -l 1 $cmd | tr '[:lower:]' '[:upper:]')(string sub -s 2 $cmd)

    complete -xc $fish_history_cmd_name -n "__fish_seen_subcommand_from $cmd" -s h -l help -d "Show the help and exit"
    complete -xc $fish_history_cmd_name -n "__fish_seen_subcommand_from $cmd" -s b -l filter-branch -d "$Cmd with branch"
    complete -xc $fish_history_cmd_name -n "__fish_seen_subcommand_from $cmd" -s d -l filter-dir -d "$Cmd with dir"
    complete -xc $fish_history_cmd_name -n "__fish_seen_subcommand_from $cmd" -s p -l filter-hostname -d "$Cmd with hostname"
    complete -xc $fish_history_cmd_name -n "__fish_seen_subcommand_from $cmd" -s q -l query -d "$Cmd with query"
    complete -xc $fish_history_cmd_name -n "__fish_seen_subcommand_from $cmd" -s c -l filter-branch -d "$Cmd columns with options"
end

## options for sync 
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from sync' -s h -l help -d 'Show the help message'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from sync' -l interval -d 'Sync with the interval'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from sync' -l diff -d 'Sync if the diff exceeds a certain number'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from sync' -l ask -d 'Sync after the confirmation'

## options for edit
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from edit' -s h -l help -d 'Show the help message'

## options for config
and set -l keys (command history config --keys)
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from config' -s h -l help -d 'Show the help message'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from config' -l get -a "$keys" -d 'Get the config value'
complete -xc $fish_history_cmd_name -n '__fish_seen_subcommand_from config' -l keys -d 'Get the config keys'

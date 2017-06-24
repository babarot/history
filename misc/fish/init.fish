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

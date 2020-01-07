## Installation

You should get binaries in advance (see [README.md](https://github.com/b4b4r07/history/blob/master/README.md#installation)). Then you can get fish-shell support with your plugin manager. 

- Install with [fundle](https://github.com/tuvistavie/fundle)

  Add the following in your config.fish.

  ```fish
  fundle plugin 'b4b4r07/history' --path 'misc/fish'
  fundle init
  ```

- Install with [fisherman](https://github.com/fisherman/fisherman)

  You cannot choose a specific directory in a repository to load as a fish plugin.
  But you can load the local directory misc/fish as a plugin.
  If you get the binaries with `go get`, run the following command.

  ```
  fisher $GOPATH/src/github.com/b4b4r07/history/misc/fish
  ```

## Usage

You should specify some enviroment variables for using this tool.

<details>
<summary><strong><code>fish_history_cmd_name</code></strong></summary>



It should be used as an alias  of `command history`. Completions are genereted for this alias. 

</details>

<details>
<summary><strong><code>fish_history_filter_options</code></strong></summary>



It should be set `history search` option. See also `command history help search`.

</details>

<details>
<summary><strong><code>fish_history_auto_sync</code></strong></summary>



Example:

```fish
set -U fish_history_auto_sync true
```

If you set sync option (for more datail, see and run `history config`)

</details>

<details>
<summary><strong><code>fish_history_auto_sync_interval</code></strong></summary>



Example:

```zsh
set -U fish_history_auto_sync_intareval "1h"
```

</details>

## Keybindings

These are functions to use for user specified keybindings.
To save custom keybindings, put the bind statements into your `fish_user_key_bindings` function.

<details>
<summary><strong><code>__history_keybind_get</code></strong></summary>



You can set keybind for getting history.

Example:

```fish
bind \cr __history_keybind_get
```

</details>

<details>
<summary><strong><code>__history_keybind_get_all</code></strong></summary>

Ignore `fish_history_filter_options` and search all history.

Example:

```fish
bind \cr\ca __fish_history_keybind_get_all
```


</details>

<details>
<summary><strong><code>__history_keybind_get_by_dir</code></strong></summary>

It's equals to `__fish_history_keybind_get` with `fish_history_filter_options="--filter-branch --filter-dir"`.

</details>

<details>
<summary><strong><code>__history_keybind_arrow_up</code></strong></summary>

Example:

```fish
bind \cp __history_keybind_arrow_up
```

</details>

<details>
<summary><strong><code>__history_keybind_arrow_down</code></strong></summary>

Example:

```fish
bind \cn __history_keybind_arrow_down
```

</details>

---

Anyway, if you want to use it immediately please run the following commands:

```fish
set -U fish_history_cmd_name hs    # as you like
set -U fish_history_auto_sync true
set -U fish_history_filter_options "--filter-branch --filter-dir"
```

Then, add the following statements into the definition of `fish_user_key_bindings` function.

(You can edit and save `fish_user_key_bindings` by `funced fish_user_key_bindings; and funcsave fish_user_key_bindings`.)

```fish
function fish_user_key_bindings

  bind \cr __history_keybind_get
  bind \cp __history_keybind_arrow_up
  bind \cn __history_keybind_arrow_down

end
```



# snips

_CLI to help with snippets and scripts._

## Installation

Download the standalone binary from the [releases](https://github.com/JanMalch/snips/releases) page,
and put it somewhere it can be picked up as a global CLI.

## Configuration

snips requires a `config.yaml` to define one or more sources for your snippets.
You can set the location with the `SNIPS_CONFIG` environment variable.
If not set, it will look in the default [user config dir](https://pkg.go.dev/os#UserConfigDir) of your platform,
e.g. `$XDG_CONFIG_HOME/snips/config.yaml`.

```yaml
# Define one or more sources, which are directories acting as a root for your snippets
sources:
    - ./relative/
    - /absolute/
# fzf doesn't have to be installed, because it's directly used as a Go library
# All fzf configuration is optional. Displayed values are defaults.
# Make sure you use {1} to get the absolute file to the path.
fzf:
    # Whether to load fzf defaults ($FZF_DEFAULT_OPTS_FILE and $FZF_DEFAULT_OPTS).
    use_env: true
    # shell command to preview the file.
    preview: "cat {1}"
    # shell command for the 'focus:transform-preview-label' bind.
    # Set to an empty string to disable.
    preview_label: "[[ -n {} ]] && printf \" %s \" {1}"
    # shell command for the 'result:transform-list-label' bind.
    # Set to an empty string to disable.
    list_label: |-
        if [[ -z $FZF_QUERY ]]; then
          echo " $FZF_MATCH_COUNT snippets "
        else
          echo " $FZF_MATCH_COUNT snippets for '$FZF_QUERY' "
        fi
```

Feel free to customize, e.g. with [`bat`](https://github.com/sharkdp/bat) for preview via `bat --color=always --style=plain {1}`.

### Sources

Optionally, you can put a `snips.yaml` in the root of each source. 
It allows you to configure which files to actually consider as snippets and scripts.

```yaml
include:
    - "scripts/**/*.ts"
```

This way you can exclude utility files which are only used by the actual scripts.

> See the [`_example`](./_example/) directory for a working setup.

## Usage

When snips is invoked, the [`fzf`](https://github.com/junegunn/fzf) fuzzy finder displays all available snippets.
Select a snippet by pressing `Enter`. You can also set an initial query via as an argument: `snips foo`.
When only one snippet matches, it is selected automatically.

When invoked without additional options, `snips` will simply print the selected file and exit.
Using the `--copy/-c` flag will copy the file content to your system clipboard.

When invoked with the `--exec/-x` flag, it will try to run the file instead.
`snips` displays one or more options on how to run the file, which you have to confirm manually.
To do so, `snips` will analyze shebangs and file extensions, even providing options for some well-known tools.

Run `snips -h` for more details and complimentary actions.

> For ad-hoc usage, you can run `snips -w` to use the current working directory as a source.

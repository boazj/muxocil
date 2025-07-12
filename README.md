<p align="center">
    <a href="https://github.com/boazj/muxocil">
        Muxocil
    </a>
    <br />
    Muxocil is a simple tool used to create and maintain window layouts in your Terminal Emulators and Multiplexers
</p>

# Supports

## Multiplexers

| Multiplexer | Supported |
| :---------- | :-------: |
| tmux        |    ✅     |
| Zellij      |    ⏰     |

## Terminal Emulators

| Multiplexer | Supported |          OS           |
| :---------- | :-------: | :-------------------: |
| iTerm2      |    ✅     |         macOS         |
| Kitty       |    ⏰     |     Linux, macOS      |
| WezTerm     |    ⏰     | Linux, macOS, Windows |

# Installation

# Usage

### List available layouts

```
    > muxocil layouts

    List all layouts available

    Usage:
      muxocil layouts [flags]

    Flags:
      -h, --help   help for layouts

```

### List available providers

```
    > muxocil providers

    List all available multiplexing providers in muxocil configuration

    Usage:
      muxocil providers [flags]

    Flags:
      -h, --help   help for providers

```

### Use a specific layout

```
    muxocil use

    Use a layout, specified by name or path

    Usage:
      muxocil use [flags]

    Flags:
      -h, --help   help for use

```

### teamocil & itermocil backward compatible command

```
    > muxocil bc

    Command compatible with iTermocil and teamocil, can be used in an alias to replace both tools

    Usage:
      muxocil bc [flags]

    Flags:
          --edit            Edit the layout file in either $EDITOR or your preferred GUI editor
      -h, --help            help for bc
          --here            Uses the current window as the layout’s first window
          --layout string   Takes a custom file path to a YAML layout file instead of [layout-name]
          --list            Lists all available layouts in ~/.itermocil, ~/.teamocil & locations per configuration file
          --show            Shows the layout content instead of executing it

```

## Options

## Configuration

## Lay

# Backward Compatability

# Credits

A lot of the work in this project is based on 2 great open source tools that are, unfortunatley, years out of date (7 and 3 years respectively)

- [Teamocil](https://github.com/remi/teamocil) - By remi. The OG, The tool for preconfigured tmux layouts
- [iTermocil](https://github.com/TomAnthony/itermocil) - By TomAnthony. Uses Apple Script to control iTerm2 for macOS, Backward compatible with Teamocil

# Milestone 1

- Project
  - [x] Build
    - [x] Makefile
  - [x] Lint
    - [x] .golangci.yaml
  - [ ] Test
  - [ ] Publish
    - [ ] .goreleaser.yml
  - [w] Logging
  - [w] Errors
- Interface
  - [w] Command layout
  - [w] Command list
    - [x] Providers
    - [x] Layouts
  - [w] Force Provider
  - [x] Configuration
    - [x] Config options
    - [x] Config file
    - [x] Config processing
- Layout Files
  - [x] Data structure
  - [w] Marshalling
- Mux Provider Infrastructure
  - [x] Provider execution logic
  - [x] Infer Provider from env
  - [x] Layout validator
  - [x] Layout normalizer
- Support tmux
  - [x] tmux driver
    - [x] NewSession
    - [x] NewWindow
    - [x] RenameWindow
    - [x] SetWindowOption
    - [x] ListPanes
    - [x] ListWindows
    - [x] RenameSession
    - [x] SelectLayout
    - [x] SelectWindow
    - [x] SelectPane
    - [x] SendKeys
    - [x] SendKeysToWindow
    - [x] SendKeysToPane
    - [x] ShowOptions
    - [x] ShowWindowOptions
    - [x] SplitWindow
  - [w] tmux provider
    - [x] Create Session
    - [x] Create Window
    - [x] Window Commands
    - [x] Create Pane
    - [x] Pane Layouts
    - [w] Custom layouts mechanism
    - [w] Pane commands
    - [ ] Windows Script
- [ ] Support Iterm2
- Verify BC with itermocil & teamocil

# Milestone 2

- Support [Kitty](https://sw.kovidgoyal.net/kitty/)
- Support [WezTerm](https://wezterm.org/)
- Support [Zellij](https://zellij.dev/)

# Potential future

- Terminal Emulators
  - [Alacritty](https://alacritty.org/)
  - [Warp](https://www.warp.dev/)
    - Not OSS
  - [termius](https://termius.com)
    - Not OSS
  - [Tabby](https://tabby.sh/)
  - [Hyper](https://hyper.is/)
  - [Windows Terminal](https://github.com/microsoft/terminal)

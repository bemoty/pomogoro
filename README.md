# pomogoro

![goro](.github/goro.jpg)

_Goro my beloved_
Art by [wroniec](https://x.com/wrroniec/status/1340335840971657216)

## Build

CGo is only required on macOS (Cocoa systray). Linux builds without it.

```
# macOS
CGO_ENABLED=1 go build -o pomogoro .

# Linux
CGO_ENABLED=0 go build -o pomogoro .
```

## Install

**Arch Linux (AUR)**
```
yay -S pomogoro-bin
```

**macOS (Homebrew)**
```
brew install --cask bemoty/tap/pomogoro
```

**macOS (dmg)** 
Download from the [releases page](https://github.com/bemoty/pomogoro/releases).

> The app is ad-hoc signed but not notarized, so on first launch macOS will
> warn that the developer cannot be verified. Right-click (or Control-click)
> the app in Finder and choose **Open**, then confirm in the dialog that
> appears — you only need to do this once.

## CLI

The IPC subcommands work on Linux and macOS while the daemon is running:

```
pomogoro status   # print current state, e.g. "W 22:14 2/4 standing"
pomogoro pause    # toggle pause/resume
pomogoro skip     # skip to next phase
pomogoro reset    # restart from pomodoro 1
pomogoro stop     # quit the daemon
```

`status` exits 1 and prints `stopped` if the daemon is not running, so it is safe to use in your favorite status bar:

```
# tmux
set -g status-right "#(pomogoro status) | %H:%M"

# skhd (macOS)
cmd + shift - p : pomogoro pause
```

## Notes

- Linux: confirmed to work with KDE, your mileage may vary with other DEs

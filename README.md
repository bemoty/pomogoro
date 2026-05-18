# pomogoro

![goro](.github/goro.jpg)
_Goro my beloved_

Art by [wroniec](https://x.com/wrroniec/status/1340335840971657216)

## Build

Needs CGo unfortunately

```
CGO_ENABLED=1 go build -o pomogoro .
```

## Install

```
cp pomogoro ~/.local/bin/pomogoro
```

Create `~/.local/share/applications/pomogoro.desktop`:

```ini
[Desktop Entry]
Type=Application
Name=pomogoro
Comment=Simple Pomodoro timer
Exec=/home/josh/.local/bin/pomogoro -d
Icon=chronometer
Terminal=false
Categories=Utility;Clock;
```

## CLI

While the daemon is running, subcommands talk to it over a unix socket:

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

## Linux only

Confirmed to work with KDE, your mileage may vary with other DEs

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

## Linux only

Confirmed to work with KDE, your mileage may vary with other DEs

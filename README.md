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

Create `~/.config/autostart/pomogoro.desktop`:

```ini
[Desktop Entry]
Type=Application
Name=pomogoro
Exec=%h/.local/bin/pomogoro
X-KDE-autostart-after=panel
```

## Linux only

Confirmed to work with KDE, your mileage may vary with other DEs

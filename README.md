# aksara

A terminal desktop.
![aksara demo](assets/demo.gif)

Aksara is a fullscreen TUI that works like a tiny desktop environment.
Instead of opening a launcher, file manager and terminal separately, I wanted
one place where I could just type commands.

press `e` and start typing. it matches your saved commands, arrow keys to pick one, enter to run it. esc gets you out, `q` quits.

## commands

everything lsited here happens in the search bar:

- `:w name : "some command"` — save a command. e.g. `:w oc : "opencode"`
- `:r name` — remove it
- `:bg "~/.config/aksara/wall.png"` — change wallpaper
- `:! whatever` — run shell directly, output shows up on the desktop plane (scroll with up/down or j/k, esc to go back)

`:` also focuses the bar if you're not in it.

## config

first run creates everything for you:

`~/.config/aksara/config.json`

```json
{
  "commands": {
    "oc": "opencode",
    "nvim": "nvim"
  },
  "settings": {
    "wallpaper": "/home/you/.config/aksara/wall.png"
  }
}
```

edit it by hand if you want, the app picks it up next time you focus the bar. your wallpaper also gets installed there automatically so a fresh machine just works.

## installation

easiest way (linux and macos):

```sh
curl -fsSL https://raw.githubusercontent.com/Dat-one-dev/Aksara/main/install.sh | sh
```

then run `aksara`. windows folks grab `Aksara_Windows_amd64.tar.gz` from the latest github release.

from source, you need go:

```sh
go build -o aksara .
./aksara
```

or run straight with `go run .`

flag if you need it: `-w /path/to/image.png` overrides the wallpaper for that session.

## built with

- go, 
- bubbletea — the whole tui
- lipgloss — styling and boxes
- bubbles — the search input
- prana — my own library lol. i literally made a whole ass library just for this project: https://github.com/Dat-one-dev/Prana.
## author

Dat-One-Dev. if you like it give it a star or whatever, if its broken open an issue.

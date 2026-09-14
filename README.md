# aksara

terminal desktop lwk

a fullscreen terminal thing that looks like a tiny desktop. it has a wallpaper, a search bar on top and a little info panel at the bottom. instead of clicking icons you just type.

i made it because opening a file manager and a launcher and a terminal separately felt dumb when i basically live in the terminal anyway. so now the terminal IS the desktop. kinda.

![aksara demo](assets/demo.gif)

## what it does

when you start it you get:

- a top search bar (looks like a browser bar)
- your wallpaper rendered as ascii art in the middle
- a bottom panel with battery, ram and time

press `e` and start typing. it matches your saved commands, arrow keys to pick one, enter to run it. esc gets you out, `q` quits.

## commands

everything happens in the search bar:

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

you need go. then:

```sh
go build -o aksara .
./aksara
```

or run straight with `go run .`

flag if you need it: `-w /path/to/image.png` overrides the wallpaper for that session.

## built with

- go, obviously
- bubbletea — the whole tui
- lipgloss — styling and boxes
- bubbles — the search input
- prana — my own library lol. i literally made a whole ass library just for this project: https://github.com/Dat-one-dev/Prana. it does the boxes and terminal sizing stuff. felt cool to depend on something i wrote myself instead of only other people's code

wallpaper ascii part uses `nfnt/resize` under the hood.

## what i learned

state-driven tuis are easy until they arent. the search popup seemed simple and then focus handling ate like three sessions. also gui apps print garbage to your terminal if you let them keep stderr (xdg-open im looking at you), redirect that stuff.

biggest lesson is probably the same as always: i keep adding panels and features and then deleting them. the side panels lasted one whole session before i ripped them out. less is more lwk.

## author

Dat-One-Dev. if you like it give it a star or whatever, if its broken open an issue.

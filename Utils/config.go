package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Entry struct {
	Cmd      string `json:"cmd"`
	Shortcut string `json:"shortcut,omitempty"`
	Color    string `json:"color,omitempty"`
	Fav      bool   `json:"fav,omitempty"`
}

type Settings struct {
	Wallpaper string   `json:"wallpaper"`
	Recents   []string `json:"recents,omitempty"`
}

type Config struct {
	Commands map[string]Entry `json:"commands"`
	Settings Settings         `json:"settings"`
}

func Path() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/aksara/config.json"
	}
	return filepath.Join(home, ".config", "aksara", "config.json")
}

func WallPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/aksara/wall.png"
	}
	return filepath.Join(home, ".config", "aksara", "wall.png")
}

func EnsureWallpaper(data []byte) string {
	p := WallPath()
	if _, err := os.Stat(p); err == nil {
		return p
	}
	if len(data) == 0 {
		return p
	}
	_ = os.MkdirAll(filepath.Dir(p), 0755)
	if err := os.WriteFile(p, data, 0644); err != nil {
		return p
	}
	return p
}

func defaults() map[string]Entry {
	return map[string]Entry{
		"oc":       {Cmd: "opencode", Shortcut: "O", Color: "205"},
		"zed":      {Cmd: "zed", Shortcut: "Z", Color: "81", Fav: true},
		"aseprite": {Cmd: "aseprite", Shortcut: "A", Color: "208"},
		"files":    {Cmd: "xdg-open .", Shortcut: "F", Color: "245"},
		"minitone": {Cmd: "minitone", Shortcut: "M", Color: "135"},
		"codex":    {Cmd: "codex", Shortcut: "C", Color: "120"},
		"nvim":     {Cmd: "nvim", Shortcut: "N", Color: "118", Fav: true},
		"firefox":  {Cmd: "firefox", Shortcut: "W", Color: "75"},
		"vlc":      {Cmd: "vlc", Shortcut: "V", Color: "201"},
	}
}

func fixEntry(name string, e Entry) Entry {
	if e.Shortcut == "" && name != "" {
		e.Shortcut = strings.ToUpper(name[:1])
	}
	e.Shortcut = strings.ToUpper(e.Shortcut)
	if e.Color == "" {
		e.Color = "205"
	}
	return e
}

func Load() (Config, error) {
	p := Path()
	data, err := os.ReadFile(p)
	if err != nil {
		def := Config{Commands: defaults(), Settings: Settings{Wallpaper: WallPath()}}
		if mkErr := Save(def); mkErr != nil {
			return def, mkErr
		}
		return def, nil
	}
	var raw struct {
		Commands map[string]json.RawMessage `json:"commands"`
		Settings Settings                   `json:"settings"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Config{Commands: map[string]Entry{}}, err
	}
	c := Config{Commands: map[string]Entry{}, Settings: raw.Settings}
	migrated := false
	for k, v := range raw.Commands {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			if d, ok := defaults()[k]; ok && d.Cmd == s {
				c.Commands[k] = d
			} else {
				c.Commands[k] = fixEntry(k, Entry{Cmd: s})
			}
			migrated = true
			continue
		}
		var e Entry
		if err := json.Unmarshal(v, &e); err == nil {
			c.Commands[k] = fixEntry(k, e)
			continue
		}
	}
	added := migrated
	for k, v := range defaults() {
		if _, ok := c.Commands[k]; !ok {
			c.Commands[k] = v
			added = true
		}
	}
	if c.Settings.Wallpaper == "" {
		c.Settings.Wallpaper = WallPath()
		added = true
	}
	taken := map[string]string{}
	var keys []string
	for k := range c.Commands {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		e := c.Commands[k]
		up := strings.ToUpper(e.Shortcut)
		if other, dup := taken[up]; dup {
			_ = other
			fixed := false
			for _, r := range strings.ToUpper(k) {
				s := string(r)
				if _, dup2 := taken[s]; !dup2 {
					e.Shortcut = s
					fixed = true
					break
				}
			}
			if !fixed {
				for ch := 'A'; ch <= 'Z'; ch++ {
					s := string(ch)
					if _, dup2 := taken[s]; !dup2 {
						e.Shortcut = s
						fixed = true
						break
					}
				}
			}
			if !fixed {
				e.Shortcut = up + "*"
			}
			c.Commands[k] = e
			added = true
		}
		taken[strings.ToUpper(e.Shortcut)] = k
	}
	var recents []string
	for _, r := range c.Settings.Recents {
		if _, ok := c.Commands[r]; ok {
			recents = append(recents, r)
		}
	}
	if len(recents) != len(c.Settings.Recents) {
		c.Settings.Recents = recents
		added = true
	}
	if added {
		_ = Save(c)
	}
	return c, nil
}

func PushRecent(c *Config, name string) {
	var out []string
	out = append(out, name)
	for _, r := range c.Settings.Recents {
		if r != name {
			out = append(out, r)
		}
	}
	if len(out) > 8 {
		out = out[:8]
	}
	c.Settings.Recents = out
}

func Save(c Config) error {
	p := Path()
	if err := os.MkdirAll(filepath.Dir(p), 0755); err != nil {
		return err
	}
	for k, e := range c.Commands {
		c.Commands[k] = fixEntry(k, e)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

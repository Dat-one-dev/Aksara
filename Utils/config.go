package utils

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Settings struct {
	Wallpaper string   `json:"wallpaper"`
	Recents   []string `json:"recents,omitempty"`
}

type Config struct {
	Commands map[string]string `json:"commands"`
	Settings Settings          `json:"settings"`
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

func defaults() map[string]string {
	return map[string]string{
		"oc":       "opencode",
		"zed":      "zed",
		"aseprite": "aseprite",
		"files":    "xdg-open .",
		"minitone": "minitone",
		"codex":    "codex",
		"nvim":     "nvim",
		"firefox":  "firefox",
		"vlc":      "vlc",
	}
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
		return Config{Commands: map[string]string{}}, err
	}
	c := Config{Commands: map[string]string{}, Settings: raw.Settings}
	migrated := false
	for k, v := range raw.Commands {
		var s string
		if err := json.Unmarshal(v, &s); err == nil {
			c.Commands[k] = s
			continue
		}
		var e struct {
			Cmd string `json:"cmd"`
		}
		if err := json.Unmarshal(v, &e); err == nil {
			c.Commands[k] = e.Cmd
			migrated = true
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
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(p, data, 0644)
}

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	desktop "github.com/Dat-one-dev/Aksara/Desktop"
	utils "github.com/Dat-one-dev/Aksara/Utils"
	prana "github.com/Dat-one-dev/Prana"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

//go:embed wall.png
var defaultWall []byte

type doneMsg struct {
	err error
}

type cmdOutMsg struct {
	out string
	err error
}

type model struct {
	w       int
	h       int
	wp      *desktop.Wallpaper
	focused bool
	input   textinput.Model
	cmds    map[string]string
	status  string
	sel     int
	wall    string
	recents []string
	showOut bool
	outLine []string
	outOff  int
}

func saveCfg(m *model) error {
	return utils.Save(utils.Config{
		Commands: m.cmds,
		Settings: utils.Settings{Wallpaper: m.wall, Recents: m.recents},
	})
}

func pushRecent(m *model, name string) {
	var out []string
	out = append(out, name)
	for _, r := range m.recents {
		if r != name {
			out = append(out, r)
		}
	}
	if len(out) > 8 {
		out = out[:8]
	}
	m.recents = out
	_ = saveCfg(m)
}

func matches(cmds map[string]string, q string) []string {
	q = strings.ToLower(strings.TrimSpace(q))
	if q == "" {
		return nil
	}
	var out []string
	for k := range cmds {
		if strings.Contains(strings.ToLower(k), q) {
			out = append(out, k)
		}
	}
	sort.Strings(out)
	if len(out) > 5 {
		out = out[:5]
	}
	return out
}

func runCmd(s string) tea.Cmd {
	c := exec.Command("sh", "-c", "("+s+") 2>/dev/null")
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return doneMsg{err}
	})
}

func runCapture(s string) tea.Cmd {
	return func() tea.Msg {
		c := exec.Command("sh", "-c", s)
		b, err := c.CombinedOutput()
		return cmdOutMsg{out: string(b), err: err}
	}
}

func parseWrite(q string) (string, string, bool) {
	rest := strings.TrimSpace(strings.TrimPrefix(q, ":w"))
	rest = strings.TrimSpace(rest)
	i := strings.Index(rest, ":")
	if i < 0 {
		return "", "", false
	}
	name := strings.TrimSpace(rest[:i])
	cmd := strings.TrimSpace(rest[i+1:])
	cmd = strings.Trim(cmd, "\"'")
	if name == "" || cmd == "" || strings.Contains(name, " ") {
		return "", "", false
	}
	return name, cmd, true
}

func parseRemove(q string) (string, bool) {
	rest := strings.TrimSpace(strings.TrimPrefix(q, ":r"))
	rest = strings.TrimSpace(rest)
	if rest == "" || strings.Contains(rest, " ") {
		return "", false
	}
	return rest, true
}

func parseBg(q string) (string, bool) {
	rest := strings.TrimSpace(strings.TrimPrefix(q, ":bg"))
	rest = strings.TrimSpace(rest)
	rest = strings.Trim(rest, "\"'")
	if rest == "" {
		return "", false
	}
	if strings.HasPrefix(rest, "~/") {
		if home, err := os.UserHomeDir(); err == nil {
			rest = home + rest[1:]
		}
	}
	return rest, true
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case doneMsg:
		if msg.err != nil {
			m.status = "error: " + msg.err.Error()
		} else {
			m.status = "done"
		}
		m.focused = false
		m.input.Blur()
		m.input.SetValue("")
		m.sel = 0
		return m, nil
	case cmdOutMsg:
		m.outLine = strings.Split(strings.TrimRight(msg.out, "\n"), "\n")
		if len(m.outLine) == 1 && m.outLine[0] == "" {
			m.outLine = []string{"(no output)"}
		}
		if msg.err != nil {
			m.outLine = append(m.outLine, "error: "+msg.err.Error())
		}
		m.showOut = true
		m.outOff = 0
		m.focused = false
		m.input.Blur()
		m.input.SetValue("")
		m.sel = 0
		m.status = ""
		return m, nil
	case tea.KeyMsg:
		if m.showOut && !m.focused {
			switch msg.Type {
			case tea.KeyCtrlC:
				return m, tea.Quit
			case tea.KeyEsc, tea.KeyEnter:
				m.showOut = false
				m.outOff = 0
				return m, nil
			case tea.KeyUp:
				if m.outOff > 0 {
					m.outOff--
				}
				return m, nil
			case tea.KeyDown:
				m.outOff++
				return m, nil
			default:
				if msg.String() == "q" {
					return m, tea.Quit
				}
				if msg.String() == "k" && m.outOff > 0 {
					m.outOff--
					return m, nil
				}
				if msg.String() == "j" {
					m.outOff++
					return m, nil
				}
				return m, nil
			}
		}
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.focused {
				m.focused = false
				m.input.Blur()
				m.input.SetValue("")
				m.sel = 0
				return m, nil
			}
		case tea.KeyUp:
			if m.focused {
				ms := matches(m.cmds, m.input.Value())
				if len(ms) > 0 {
					m.sel = (m.sel - 1 + len(ms)) % len(ms)
				}
				return m, nil
			}
		case tea.KeyDown:
			if m.focused {
				ms := matches(m.cmds, m.input.Value())
				if len(ms) > 0 {
					m.sel = (m.sel + 1) % len(ms)
				}
				return m, nil
			}
		case tea.KeyEnter:
			if m.focused {
				q := strings.TrimSpace(m.input.Value())
				if q == "" {
					return m, nil
				}
				if strings.HasPrefix(q, ":!") {
					rest := strings.TrimSpace(strings.TrimPrefix(q, ":!"))
					if rest == "" {
						m.status = "use :! <shell command>"
						return m, nil
					}
					m.status = "running " + rest
					m.showOut = false
					m.outOff = 0
					m.focused = false
					m.input.Blur()
					return m, runCapture(rest)
				}
				if strings.HasPrefix(q, ":w") {
					name, cmd, ok := parseWrite(q)
					if !ok {
						m.status = "use :w name : \"cmd\""
						return m, nil
					}
					if m.cmds == nil {
						m.cmds = map[string]string{}
					}
					m.cmds[name] = cmd
					if err := saveCfg(&m); err != nil {
						m.status = "error: " + err.Error()
						return m, nil
					}
					m.status = "saved " + name
					m.input.SetValue("")
					m.sel = 0
					return m, nil
				}
				if strings.HasPrefix(q, ":r") {
					name, ok := parseRemove(q)
					if !ok {
						m.status = "use :r name"
						return m, nil
					}
					if _, ok := m.cmds[name]; !ok {
						m.status = "no match for " + name
						return m, nil
					}
					delete(m.cmds, name)
					var kept []string
					for _, r := range m.recents {
						if r != name {
							kept = append(kept, r)
						}
					}
					m.recents = kept
					if err := saveCfg(&m); err != nil {
						m.status = "error: " + err.Error()
						return m, nil
					}
					m.status = "removed " + name
					m.input.SetValue("")
					m.sel = 0
					return m, nil
				}
				if strings.HasPrefix(q, ":bg") {
					path, ok := parseBg(q)
					if !ok {
						m.status = "use :bg \"~/.config/aksara/wall.png\""
						return m, nil
					}
					wp, err := desktop.Load(path)
					if err != nil {
						m.status = "error: " + err.Error()
						return m, nil
					}
					m.wp = wp
					m.wall = path
					if err := saveCfg(&m); err != nil {
						m.status = "error: " + err.Error()
						return m, nil
					}
					m.status = "wallpaper set"
					m.input.SetValue("")
					m.sel = 0
					return m, nil
				}
				if cmd, ok := m.cmds[q]; ok {
					m.status = "running " + q
					pushRecent(&m, q)
					return m, runCmd(cmd)
				}
				ms := matches(m.cmds, q)
				if len(ms) > 0 {
					if m.sel < 0 || m.sel >= len(ms) {
						m.sel = 0
					}
					m.status = "running " + ms[m.sel]
					pushRecent(&m, ms[m.sel])
					return m, runCmd(m.cmds[ms[m.sel]])
				}
				m.status = "no match for " + q
				return m, nil
			}
		default:
			if (msg.String() == "e" || msg.String() == "E") && !m.focused {
				m.focused = true
				m.status = ""
				m.sel = 0
				if cfg, err := utils.Load(); err == nil {
					m.cmds = cfg.Commands
					m.wall = cfg.Settings.Wallpaper
					m.recents = cfg.Settings.Recents
				}
				return m, m.input.Focus()
			}
			if msg.String() == ":" && !m.focused {
				m.focused = true
				m.status = ""
				m.sel = 0
				if cfg, err := utils.Load(); err == nil {
					m.cmds = cfg.Commands
					m.wall = cfg.Settings.Wallpaper
					m.recents = cfg.Settings.Recents
				}
				m.input.SetValue(":")
				return m, m.input.Focus()
			}
			if msg.String() == "q" && !m.focused {
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.w = msg.Width
		m.h = msg.Height
	}

	if m.focused {
		old := m.input.Value()
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		if m.input.Value() != old {
			m.sel = 0
		}
		ms := matches(m.cmds, m.input.Value())
		if len(ms) > 0 {
			if m.sel < 0 || m.sel >= len(ms) {
				m.sel = 0
			}
		} else {
			m.sel = 0
		}
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.w < 40 || m.h < 10 {
		return ""
	}

	top := desktop.SearchBar(m.w, m.input.View(), m.focused)

	if m.showOut {
		return top + "\n" + desktop.Output(m.w, m.h-3, m.outLine, m.outOff)
	}

	mid := ""
	if m.focused {
		q := strings.TrimSpace(m.input.Value())
		if !strings.HasPrefix(q, ":") {
			ms := matches(m.cmds, m.input.Value())
			lines := []string{}
			for i, k := range ms {
				mark := "  "
				if i == m.sel {
					mark = "> "
				}
				lines = append(lines, mark+k+" -> "+m.cmds[k])
			}
			if m.status != "" {
				lines = append(lines, m.status)
			}
			if len(lines) > 0 {
				mid = strings.Join(lines, "\n") + "\n"
			}
		} else if m.status != "" {
			mid = m.status + "\n"
		}
	} else if m.status != "" {
		mid = m.status + "\n"
	}

	rest := m.h - 3
	used := strings.Count(mid, "\n")
	if rest-used < 10 {
		rest = 10
	} else {
		rest = rest - used
	}

	return top + "\n" + mid + desktop.Render(m.w, rest, m.wp)
}

func main() {
	wpPath := flag.String("w", "", "path to wallpaper image")
	flag.Parse()

	ti := textinput.New()
	ti.Placeholder = "Search"
	ti.CharLimit = 156
	ti.Width = 30
	ti.Blur()

	cfg, _ := utils.Load()
	wpFile := utils.EnsureWallpaper(defaultWall)
	if cfg.Settings.Wallpaper == "" {
		cfg.Settings.Wallpaper = wpFile
	}

	var wp *desktop.Wallpaper
	if *wpPath != "" {
		var err error
		wp, err = desktop.Load(*wpPath)
		if err != nil {
			fmt.Printf("Error loading wallpaper: %v\n", err)
			os.Exit(1)
		}
	} else if cfg.Settings.Wallpaper != "" {
		if loaded, err := desktop.Load(cfg.Settings.Wallpaper); err == nil {
			wp = loaded
		}
	}

	w, h := prana.TerminalSize()
	if w < 40 || h < 10 {
		w, h = 80, 24
	}

	p := tea.NewProgram(model{
		w:       w,
		h:       h,
		wp:      wp,
		input:   ti,
		focused: false,
		cmds:    cfg.Commands,
		wall:    cfg.Settings.Wallpaper,
		recents: cfg.Settings.Recents,
	}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

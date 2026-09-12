package main

import (
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

type doneMsg struct {
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
	c := exec.Command("sh", "-c", s)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return doneMsg{err}
	})
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
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEsc:
			if m.focused {
				m.focused = false
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			}
		case tea.KeyEnter:
			if m.focused {
				q := strings.TrimSpace(m.input.Value())
				if q == "" {
					return m, nil
				}
				if cmd, ok := m.cmds[q]; ok {
					m.status = "running " + q
					return m, runCmd(cmd)
				}
				ms := matches(m.cmds, q)
				if len(ms) > 0 {
					m.status = "running " + ms[0]
					return m, runCmd(m.cmds[ms[0]])
				}
				m.status = "no match for " + q
				return m, nil
			}
		default:
			if (msg.String() == "e" || msg.String() == "E") && !m.focused {
				m.focused = true
				m.status = ""
				if cfg, err := utils.Load(); err == nil {
					m.cmds = cfg.Commands
				}
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
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.w < 40 || m.h < 10 {
		return ""
	}

	top := desktop.SearchBar(m.w, m.input.View(), m.focused)

	mid := ""
	if m.focused {
		ms := matches(m.cmds, m.input.Value())
		lines := []string{}
		for _, k := range ms {
			lines = append(lines, k+" -> "+m.cmds[k])
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

	var wp *desktop.Wallpaper
	if *wpPath != "" {
		var err error
		wp, err = desktop.Load(*wpPath)
		if err != nil {
			fmt.Printf("Error loading wallpaper: %v\n", err)
			os.Exit(1)
		}
	}

	w, h := prana.TerminalSize()
	if w < 40 || h < 10 {
		w, h = 80, 24
	}

	cfg, _ := utils.Load()

	p := tea.NewProgram(model{
		w:       w,
		h:       h,
		wp:      wp,
		input:   ti,
		focused: false,
		cmds:    cfg.Commands,
	}, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

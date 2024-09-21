package display

import (
	"time"

	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
)

type model struct {
	feed    *gofeed.Feed
	timer   timer.Model
	loading bool
	err     error
}

func NewModel(feed *gofeed.Feed) model {
	timeInterval := time.Millisecond * 50
	timeout := timeInterval * time.Duration(len(feed.Items))
	return model{
		feed:    feed,
		timer:   timer.NewWithInterval(timeout, timeInterval),
		loading: true,
	}
}

func (m model) Init() tea.Cmd {
	return m.timer.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case timer.TickMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		m.loading = false
	case tea.KeyMsg:
		if m.loading {
			//return m, nil
		}
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case "m":
			return m, m.timer.Toggle()
		}
	}

	return m, nil
}

func (m model) View() string {
	return ""
}

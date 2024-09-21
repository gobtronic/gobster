package load

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/mmcdole/gofeed"
)

type model struct {
	spinner spinner.Model
	err     error
}

func NewModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{
		spinner: s,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, parseFeed)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		default:
			return m, nil
		}
	case *gofeed.Feed:
		return m, nil
	case error:
		m.err = msg
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		log.Error("An error occured while retrieving posts", "err", m.err)
		return ""
	}
	return fmt.Sprintf("%s Retrieving latest Lobsters posts...\n", m.spinner.View())
}

func parseFeed() tea.Msg {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://lobste.rs/rss")
	if err != nil {
		return err
	}
	return feed
}

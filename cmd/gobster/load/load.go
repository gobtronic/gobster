package load

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gobtronic/gobster/cmd/gobster/display"
	"github.com/mmcdole/gofeed"
)

type model struct {
	spinner  spinner.Model
	err      error
	termSize [2]int
}

func NewModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, parseFeed)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.termSize = [2]int{msg.Width, msg.Height}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		}

	case *gofeed.Feed:
		displayModel := display.NewModel(msg, m.termSize)
		return displayModel, displayModel.Init()
	case error:
		m.err = msg
		return m, tea.Quit

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		log.Error("An error occured while retrieving discussions", "err", m.err)
		return ""
	}
	style := lipgloss.NewStyle().Bold(true)
	return fmt.Sprintf("\n  %s%s\n", m.spinner.View(), style.Render("retrieving latest discussions..."))
}

func parseFeed() tea.Msg {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL("https://lobste.rs/rss")
	if err != nil {
		return err
	}
	return feed
}

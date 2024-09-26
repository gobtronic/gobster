package load

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	"github.com/gobtronic/gobster/cmd/gobster/display"
	"github.com/gobtronic/gobster/cmd/gobster/feed"
)

type model struct {
	// The load spinner model
	spinner spinner.Model
	// Any error that could occur during loading
	err error
	// The current terminal window size {width, height}
	termSize [2]int
}

// Returns a new load.model instance
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

	case *feed.LobsterFeed:
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

// Retrieves and parses the lobste.rs rss feed
// Returns the feed if it succeed, returns an error otherwise
func parseFeed() tea.Msg {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	feed, err := feed.FetchFeed(feed.Active)
	if err != nil {
		f.WriteString(err.Error() + "\n")
		return err
	}
	return &feed
}

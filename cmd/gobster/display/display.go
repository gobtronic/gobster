package display

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	listdef "github.com/gobtronic/gobster/cmd/gobster/display/list"
	"github.com/mmcdole/gofeed"
)

type model struct {
	feed            *gofeed.Feed
	initialTermSize [2]int
	list            list.Model
	err             error
}

func NewModel(feed *gofeed.Feed, initialTermSize [2]int) model {
	l := listdef.NewList(feed, initialTermSize)
	return model{
		feed:            feed,
		initialTermSize: initialTermSize,
		list:            l,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		m.list.SetHeight(msg.Height)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

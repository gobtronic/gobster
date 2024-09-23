package display

import (
	"os/exec"
	"runtime"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	listdef "github.com/gobtronic/gobster/cmd/gobster/display/list"
	"github.com/mmcdole/gofeed"
)

type model struct {
	feed *gofeed.Feed
	list list.Model
	err  error
}

// Returns a new display.model
func NewModel(feed *gofeed.Feed, initialTermSize [2]int) model {
	l := listdef.NewList(feed, initialTermSize)
	return model{
		feed: feed,
		list: l,
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
		case " ", "enter":
			item := m.list.SelectedItem().(listdef.Item)
			openInBrowser(item.Url)
			return m, nil
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

// Opens an url in the default user's browser
func openInBrowser(url string) {
	if url == "" {
		return
	}

	switch runtime.GOOS {
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		_ = exec.Command("open", url).Start()
	}
}

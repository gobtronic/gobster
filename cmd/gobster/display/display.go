package display

import (
	"fmt"
	"io"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/timer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

type model struct {
	feed    *gofeed.Feed
	list    list.Model
	timer   timer.Model
	loading bool
	err     error
}

const listHeight = 20

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	paginationStyle   = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item string

func (i item) FilterValue() string { return string(i) }

type itemDelegate struct{}

func (d itemDelegate) Height() int                             { return 1 }
func (d itemDelegate) Spacing() int                            { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle.Render
	if index == m.Index() {
		fn = func(s ...string) string {
			return selectedItemStyle.Render("➜ " + strings.Join(s, " "))
		}
	}

	fmt.Fprint(w, fn(str))
}

func NewModel(feed *gofeed.Feed) model {
	l := list.New([]list.Item{}, itemDelegate{}, 50, listHeight)
	l.Title = "lobste.rs - active discussions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	timeInterval := time.Millisecond * 30
	maxItemsPerPage := listHeight - 6
	timeout := timeInterval * time.Duration(math.Min(float64(len(feed.Items)-1), float64(maxItemsPerPage)))
	return model{
		feed:    feed,
		list:    l,
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
		if msg.Timeout {
			return m, nil
		}
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		curListItemsLen := len(m.list.Items())
		updatedItems := append(m.list.Items(), item(m.feed.Items[curListItemsLen].Title))
		m.list.SetItems(updatedItems)
		return m, cmd
	case timer.StartStopMsg:
		var cmd tea.Cmd
		m.timer, cmd = m.timer.Update(msg)
		return m, cmd
	case timer.TimeoutMsg:
		curItems := m.list.Items()
		endItems := []list.Item{}
		for i := len(curItems); i < len(m.feed.Items); i++ {
			endItems = append(endItems, item(m.feed.Items[i].Title))
		}
		curItems = append(curItems, endItems...)

		m.list.SetItems(curItems)
		m.loading = false

	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		if m.loading {
			return m, nil
		}
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

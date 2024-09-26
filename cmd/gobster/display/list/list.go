package list

import (
	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/gobtronic/gobster/cmd/gobster/feed"
)

var (
	titleStyle      = lipgloss.NewStyle().MarginLeft(2).MarginTop(1).Bold(true)
	paginationStyle = bubblelist.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle       = bubblelist.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
)

var tagDefaultBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#3b320d", Dark: "#f0e8c7"}
var tagRedBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#3b1719", Dark: "#cf8488"}
var tagGreyBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#c2c2c2", Dark: "#d4d4d4"}
var tagBlueBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#15293d", Dark: "#9fbfde"}
var dimForeground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#343434", Dark: "#727272"}

var catBackgrounds map[string]lipgloss.AdaptiveColor = map[string]lipgloss.AdaptiveColor{
	"ask":        tagRedBackground,
	"show":       tagRedBackground,
	"announce":   tagRedBackground,
	"interview":  tagRedBackground,
	"audio":      tagBlueBackground,
	"book":       tagBlueBackground,
	"pdf":        tagBlueBackground,
	"slides":     tagBlueBackground,
	"transcript": tagBlueBackground,
	"video":      tagBlueBackground,
	"meta":       tagGreyBackground,
}

type itemDelegate struct{}

func (d itemDelegate) Height() int                                   { return 2 }
func (d itemDelegate) Spacing() int                                  { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *bubblelist.Model) tea.Cmd { return nil }

func NewList(feed *feed.LobsterFeed, initialTermSize [2]int) bubblelist.Model {
	l := bubblelist.New([]bubblelist.Item{}, itemDelegate{}, initialTermSize[0], initialTermSize[1])
	l.Title = "lobste.rs - active discussions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(" ", "enter"),
				key.WithHelp("space/enter", "open in browser"),
			),
		}
	}

	items := []bubblelist.Item{}
	for _, v := range *feed {
		items = append(items, v)
	}
	l.SetItems(items)

	return l
}

package list

import (
	"fmt"
	"io"
	"strings"

	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
)

var (
	titleStyle        = lipgloss.NewStyle().MarginLeft(2).Bold(true)
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Underline(true).Foreground(lipgloss.Color("170"))
	paginationStyle   = bubblelist.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle         = bubblelist.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle     = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

var tagDefaultBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#3b320d", Dark: "#f0e8c7"}
var tagRedBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#3b1719", Dark: "#cf8488"}
var tagGreyBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#c2c2c2", Dark: "#d4d4d4"}
var tagBlueBackground lipgloss.AdaptiveColor = lipgloss.AdaptiveColor{Light: "#15293d", Dark: "#9fbfde"}

var catBackgrounds map[string]lipgloss.AdaptiveColor = map[string]lipgloss.AdaptiveColor{
	"ask":       tagRedBackground,
	"show":      tagRedBackground,
	"announce":  tagRedBackground,
	"interview": tagRedBackground,

	"audio":      tagBlueBackground,
	"book":       tagBlueBackground,
	"pdf":        tagBlueBackground,
	"slides":     tagBlueBackground,
	"transcript": tagBlueBackground,
	"video":      tagBlueBackground,

	"meta": tagGreyBackground,
}

type item struct {
	title      string
	categories []string
}

func (i item) FilterValue() string { return i.title + strings.Join(i.categories, " ") }

type itemDelegate struct{}

func (d itemDelegate) Height() int                                   { return 1 }
func (d itemDelegate) Spacing() int                                  { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *bubblelist.Model) tea.Cmd { return nil }
func (d itemDelegate) renderTitle(i item, index int, selected bool) string {
	str := fmt.Sprintf("%d. %s", index+1, i.title)
	if selected {
		str = selectedItemStyle.Render("➜ " + str)
	} else {
		str = itemStyle.Render(str)
	}
	return str
}
func (d itemDelegate) renderCategories(i item) string {
	style := lipgloss.NewStyle().PaddingLeft(1)
	fmtCategories := []string{}
	for _, c := range i.categories {
		fmtCategories = append(fmtCategories, d.renderCategory(c))
	}
	return "" + style.Render(strings.Join(fmtCategories, " "))
}
func (d itemDelegate) renderCategory(cat string) string {
	bgColor := tagDefaultBackground
	if color, ok := catBackgrounds[cat]; ok {
		bgColor = color
	}
	style := lipgloss.NewStyle().Foreground(bgColor).Italic(true)
	return style.Render("<" + cat + ">")
}
func (d itemDelegate) Render(w io.Writer, m bubblelist.Model, index int, listItem bubblelist.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := d.renderTitle(i, index, index == m.Index())
	str += d.renderCategories(i)
	fmt.Fprint(w, str)
}

func NewList(feed *gofeed.Feed, initialTermSize [2]int) bubblelist.Model {
	l := bubblelist.New([]bubblelist.Item{}, itemDelegate{}, initialTermSize[0], initialTermSize[1])
	l.Title = "lobste.rs - active discussions"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	items := []bubblelist.Item{}
	for _, v := range feed.Items {
		items = append(items, item{
			title:      v.Title,
			categories: v.Categories,
		})
	}
	l.SetItems(items)

	return l
}

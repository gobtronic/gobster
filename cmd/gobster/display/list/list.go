package list

import (
	"fmt"
	"io"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	bubblelist "github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mmcdole/gofeed"
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

type Item struct {
	title      string
	categories []string
	Url        string
}

func (i Item) FilterValue() string { return i.title + strings.Join(i.categories, " ") }

type itemDelegate struct{}

func (d itemDelegate) Height() int                                   { return 1 }
func (d itemDelegate) Spacing() int                                  { return 0 }
func (d itemDelegate) Update(_ tea.Msg, _ *bubblelist.Model) tea.Cmd { return nil }
func (d itemDelegate) renderMainLine(i Item, index int, selected bool) string {
	style := lipgloss.NewStyle().PaddingLeft(4)
	if selected {
		style = style.PaddingLeft(2)
	}
	indexStr := d.renderIndex(index, selected)
	titleStr := d.renderTitle(i.title, selected)
	categoriesStr := d.renderCategories(i.categories)
	str := fmt.Sprintf("%s %s %s", indexStr, titleStr, categoriesStr)
	return style.Render(str)
}
func (d itemDelegate) renderIndex(index int, selected bool) string {
	style := lipgloss.NewStyle()
	fmtIndex := fmt.Sprintf("%d.", index+1)
	if selected {
		style = style.Bold(true)
		return style.Render("➜ " + fmtIndex)
	}
	return style.Render(fmtIndex)
}
func (d itemDelegate) renderTitle(title string, selected bool) string {
	style := lipgloss.NewStyle()
	str := fmt.Sprintf("%s", title)
	if selected {
		style = style.Underline(true)
	}
	return style.Render(str)
}
func (d itemDelegate) renderCategories(categories []string) string {
	style := lipgloss.NewStyle()
	fmtCategories := []string{}
	for _, c := range categories {
		fmtCategories = append(fmtCategories, d.renderCategory(c))
	}
	return style.Render(strings.Join(fmtCategories, " "))
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
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	selected := index == m.Index()
	str := d.renderMainLine(i, index, selected)
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
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			key.NewBinding(
				key.WithKeys(" ", "enter"),
				key.WithHelp("space/enter", "open in browser"),
			),
		}
	}

	items := []bubblelist.Item{}
	for _, v := range feed.Items {
		items = append(items, Item{
			title:      v.Title,
			categories: v.Categories,
			Url:        v.Link,
		})
	}
	l.SetItems(items)

	return l
}

package list

import (
	"fmt"
	"io"
	"strings"
	"time"

	bubblelist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

const (
	itemPrefixLength = len("➜  1. ")
	dateFormat       = "Mon _2 Jan 2006 - 15:04"
)

type styleProvider struct {
	mainLine   lipgloss.Style
	index      lipgloss.Style
	title      lipgloss.Style
	categories lipgloss.Style
	category   lipgloss.Style
	date       lipgloss.Style
}

func newStyleProvider(selected bool) styleProvider {
	mainLinePadding := 4
	if selected {
		mainLinePadding = 2
	}
	mainLineStyle := lipgloss.NewStyle().PaddingLeft(mainLinePadding)

	indexStyle := lipgloss.NewStyle()

	titleStyle := lipgloss.NewStyle()
	if selected {
		titleStyle = titleStyle.Underline(true)
	}

	categoriesStyle := lipgloss.NewStyle()
	categoryStyle := lipgloss.NewStyle().Italic(true)
	dateStyle := lipgloss.NewStyle().Foreground(dimForeground)

	return styleProvider{
		mainLine:   mainLineStyle,
		index:      indexStyle,
		title:      titleStyle,
		categories: categoriesStyle,
		category:   categoryStyle,
		date:       dateStyle,
	}
}

func (d itemDelegate) Render(w io.Writer, m bubblelist.Model, index int, listItem bubblelist.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	selected := index == m.Index()
	str := d.renderMainLine(newStyleProvider(selected), i, index, selected)
	fmt.Fprint(w, str)
}

// Renders the item's main line (index, title and categories)
func (d itemDelegate) renderMainLine(styles styleProvider, i Item, index int, selected bool) string {
	style := styles.mainLine
	indexStr := d.renderIndex(styles.index, index, selected)
	titleStr := d.renderTitle(styles.title, i.title)
	categoriesStr := d.renderCategories(styles.categories, styles.category, i.categories)
	dateStr := d.renderDate(styles.date, i.date)
	str := fmt.Sprintf("%s %s %s\n%[4]*s%s", indexStr, titleStr, categoriesStr, itemPrefixLength-style.GetPaddingLeft(), "", dateStr)
	return style.Render(str)
}

// Renders the item's index
func (d itemDelegate) renderIndex(style lipgloss.Style, index int, selected bool) string {
	fmtIndex := fmt.Sprintf("%2d.", index+1)
	if selected {
		return style.Render("➜ " + fmtIndex)
	}
	return style.Render(fmtIndex)
}

// Renders the item's title
func (d itemDelegate) renderTitle(style lipgloss.Style, title string) string {
	return style.Render(title)
}

// Renders the item's categories
func (d itemDelegate) renderCategories(style lipgloss.Style, categoryStyle lipgloss.Style, categories []string) string {
	fmtCategories := []string{}
	for _, c := range categories {
		fmtCategories = append(fmtCategories, d.renderCategory(categoryStyle, c))
	}
	return style.Render(strings.Join(fmtCategories, " "))
}

// Renders a single category
func (d itemDelegate) renderCategory(style lipgloss.Style, cat string) string {
	bgColor := tagDefaultBackground
	if color, ok := catBackgrounds[cat]; ok {
		bgColor = color
	}
	style = style.Foreground(bgColor)
	return style.Render("<" + cat + ">")
}

func (d itemDelegate) renderDate(style lipgloss.Style, date *time.Time) string {
	if date == nil {
		return ""
	}
	return style.Render(date.Format(dateFormat))
}

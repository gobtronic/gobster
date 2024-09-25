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

func (d itemDelegate) Render(w io.Writer, m bubblelist.Model, index int, listItem bubblelist.Item) {
	i, ok := listItem.(Item)
	if !ok {
		return
	}

	selected := index == m.Index()
	str := d.renderMainLine(i, index, selected)
	fmt.Fprint(w, str)
}

// Renders the item's main line (index, title and categories)
func (d itemDelegate) renderMainLine(i Item, index int, selected bool) string {
	padding := 4
	if selected {
		padding = 2
	}

	style := lipgloss.NewStyle().PaddingLeft(padding)

	indexStr := d.renderIndex(index, selected)
	titleStr := d.renderTitle(i.title, selected)
	categoriesStr := d.renderCategories(i.categories)
	dateStr := d.renderDate(i.date)
	str := fmt.Sprintf("%s %s %s\n%[4]*s%s", indexStr, titleStr, categoriesStr, itemPrefixLength-padding, "", dateStr)
	return style.Render(str)
}

// Renders the item's index
func (d itemDelegate) renderIndex(index int, selected bool) string {
	style := lipgloss.NewStyle()
	fmtIndex := fmt.Sprintf("%2d.", index+1)
	if selected {
		style = style.Bold(true)
		return style.Render("➜ " + fmtIndex)
	}
	return style.Render(fmtIndex)
}

// Renders the item's title
func (d itemDelegate) renderTitle(title string, selected bool) string {
	style := lipgloss.NewStyle()
	str := fmt.Sprintf("%s", title)
	if selected {
		style = style.Underline(true)
	}
	return style.Render(str)
}

// Renders the item's categories
func (d itemDelegate) renderCategories(categories []string) string {
	style := lipgloss.NewStyle()
	fmtCategories := []string{}
	for _, c := range categories {
		fmtCategories = append(fmtCategories, d.renderCategory(c))
	}
	return style.Render(strings.Join(fmtCategories, " "))
}

// Renders a single category
func (d itemDelegate) renderCategory(cat string) string {
	bgColor := tagDefaultBackground
	if color, ok := catBackgrounds[cat]; ok {
		bgColor = color
	}
	style := lipgloss.NewStyle().Foreground(bgColor).Italic(true)
	return style.Render("<" + cat + ">")
}

func (d itemDelegate) renderDate(date *time.Time) string {
	if date == nil {
		return ""
	}

	style := lipgloss.NewStyle().Foreground(dimForeground)

	return style.Render(date.Format(dateFormat))
}

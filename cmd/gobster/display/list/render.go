package list

import (
	"fmt"
	"io"
	"strings"
	"time"

	bubblelist "github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/gobtronic/gobster/cmd/gobster/feed"
	"github.com/gobtronic/gobster/cmd/gobster/format"
)

const (
	itemPrefixLength = len("➜  1. ")
	dateFormat       = "Mon _2 Jan 2006 - 15:04"
)

type styleProvider struct {
	mainLine lipgloss.Style
	index    lipgloss.Style
	title    lipgloss.Style
	tags     lipgloss.Style
	tag      lipgloss.Style
	date     lipgloss.Style
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

	tagsStyle := lipgloss.NewStyle()
	tagStyle := lipgloss.NewStyle().Italic(true)
	dateStyle := lipgloss.NewStyle().Foreground(dimForeground)

	return styleProvider{
		mainLine: mainLineStyle,
		index:    indexStyle,
		title:    titleStyle,
		tags:     tagsStyle,
		tag:      tagStyle,
		date:     dateStyle,
	}
}

func (d itemDelegate) Render(w io.Writer, m bubblelist.Model, index int, listItem bubblelist.Item) {
	i, ok := listItem.(feed.Item)
	if !ok {
		return
	}

	selected := index == m.Index()
	str := d.renderItem(newStyleProvider(selected), i, index, selected)
	fmt.Fprint(w, str)
}

// Renders the item and all its subcomponents
func (d itemDelegate) renderItem(styles styleProvider, i feed.Item, index int, selected bool) string {
	style := styles.mainLine
	indexStr := d.renderIndex(styles.index, index, selected)
	titleStr := d.renderTitle(styles.title, i.Title)
	tagsStr := d.renderTags(styles.tags, styles.tag, i.Tags)
	dateStr := d.renderDate(styles.date, &i.CreatedAt.Time)
	str := fmt.Sprintf("%s %s\n%[3]*s%s %s", indexStr, titleStr, itemPrefixLength-style.GetPaddingLeft(), "", dateStr, tagsStr)
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

func (d itemDelegate) renderDate(style lipgloss.Style, date *time.Time) string {
	if date == nil {
		return ""
	}
	return style.Render(format.FmtRelativeDateToNow(date))
}

// Renders the item's tags
func (d itemDelegate) renderTags(style lipgloss.Style, tagStyle lipgloss.Style, tags []string) string {
	fmtTags := []string{}
	for _, t := range tags {
		fmtTags = append(fmtTags, d.renderTag(tagStyle, t))
	}
	return style.Render(strings.Join(fmtTags, " "))
}

// Renders a single tag
func (d itemDelegate) renderTag(style lipgloss.Style, cat string) string {
	bgColor := tagDefaultBackground
	if color, ok := catBackgrounds[cat]; ok {
		bgColor = color
	}
	style = style.Foreground(bgColor)
	return style.Render("<" + cat + ">")
}

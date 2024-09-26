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
	dimmed   lipgloss.Style
	score    lipgloss.Style
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
	dimmedStyle := lipgloss.NewStyle().Foreground(dimForeground)
	scoreStyle := lipgloss.NewStyle()

	return styleProvider{
		mainLine: mainLineStyle,
		index:    indexStyle,
		title:    titleStyle,
		tags:     tagsStyle,
		tag:      tagStyle,
		dimmed:   dimmedStyle,
		score:    scoreStyle,
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
	scoreStr := d.renderScore(styles.dimmed, i.Score)
	dateStr := d.renderDate(styles.dimmed, &i.CreatedAt.Time)
	bottomComponents := []string{
		scoreStr,
		dateStr,
		tagsStr,
	}
	str := fmt.Sprintf("%s %s\n%[3]*s%s", indexStr, titleStr, itemPrefixLength-style.GetPaddingLeft(), "", strings.Join(bottomComponents, styles.dimmed.Render(" · ")))
	return style.Render(str)
}

func (d itemDelegate) renderIndex(style lipgloss.Style, index int, selected bool) string {
	fmtIndex := fmt.Sprintf("%2d.", index+1)
	if selected {
		return style.Render("➜ " + fmtIndex)
	}
	return style.Render(fmtIndex)
}

func (d itemDelegate) renderTitle(style lipgloss.Style, title string) string {
	return style.Render(title)
}

func (d itemDelegate) renderScore(style lipgloss.Style, score int) string {
	return style.Render(fmt.Sprintf("▲ %d", score))
}

func (d itemDelegate) renderDate(style lipgloss.Style, date *time.Time) string {
	if date == nil {
		return ""
	}
	return style.Render(format.FmtRelativeDateToNow(date))
}

func (d itemDelegate) renderTags(style lipgloss.Style, tagStyle lipgloss.Style, tags []string) string {
	fmtTags := []string{}
	for _, t := range tags {
		fmtTags = append(fmtTags, d.renderTag(tagStyle, t))
	}
	return style.Render(strings.Join(fmtTags, " "))
}

func (d itemDelegate) renderTag(style lipgloss.Style, cat string) string {
	bgColor := tagDefaultBackground
	if color, ok := catBackgrounds[cat]; ok {
		bgColor = color
	}
	style = style.Foreground(bgColor)
	return style.Render("<" + cat + ">")
}

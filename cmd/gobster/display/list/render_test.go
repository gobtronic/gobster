package list

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/gobtronic/gobster/cmd/gobster/feed"
	"github.com/gobtronic/gobster/cmd/gobster/format"
	"github.com/stretchr/testify/assert"
)

var (
	testDate = time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC)
)

func TestRenderIndex(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}

	str := delegate.renderIndex(styles.index, 3, false)

	assert.Equal(t, styles.index.Render(" 4."), str)
}

func TestRenderIndexSelected(t *testing.T) {
	styles := newStyleProvider(true)
	delegate := itemDelegate{}

	str := delegate.renderIndex(styles.index, 3, true)

	assert.Equal(t, styles.index.Render("➜  4."), str)
}

func TestRenderTitle(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}

	str := delegate.renderTitle(styles.title, "Title")

	assert.Equal(t, styles.title.Render("Title"), str)
}

func TestRenderDate(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	now := time.Now()
	expected := format.FmtRelativeDateToNow(&now)

	str := delegate.renderDate(styles.date, &now)

	assert.Equal(t, styles.date.Render(expected), str)
}

func TestRenderTags(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	tags := []string{"programming", "go"}
	fmtTags := []string{}
	for _, t := range tags {
		fmtTags = append(fmtTags, delegate.renderTag(styles.tag, t))
	}

	str := delegate.renderTags(styles.tags, styles.tag, []string{"programming", "go"})

	assert.Equal(t, styles.tags.Render(strings.Join(fmtTags, " ")), str)
}

func TestRenderTag(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	cat := "programming"
	catSpecificForeground, exists := catBackgrounds[cat]
	if !exists {
		catSpecificForeground = tagDefaultBackground
	}

	str := delegate.renderTag(styles.tag, cat)
	expectedCatStyle := styles.tag.Foreground(catSpecificForeground)

	assert.Equal(t, expectedCatStyle.Render(fmt.Sprintf("<%s>", cat)), str)
}

func TestRenderItem(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	item := feed.Item{
		Title:     "Title",
		CreatedAt: feed.ItemTime{Time: testDate},
		Tags: []string{
			"programming",
			"go",
		},
	}
	indexStr := delegate.renderIndex(styles.index, 1, false)
	titleStr := delegate.renderTitle(styles.title, item.Title)
	tagsStr := delegate.renderTags(styles.tags, styles.tag, item.Tags)
	dateStr := delegate.renderDate(styles.date, &item.CreatedAt.Time)
	expected := fmt.Sprintf("%s %s\n%[3]*s%s %s", indexStr, titleStr, itemPrefixLength-styles.mainLine.GetPaddingLeft(), "", dateStr, tagsStr)

	str := delegate.renderItem(styles, item, 1, false)

	assert.Equal(t, styles.mainLine.Render(expected), str)
}

func TestRenderItemSelected(t *testing.T) {
	styles := newStyleProvider(true)
	delegate := itemDelegate{}
	item := feed.Item{
		Title:     "Title",
		CreatedAt: feed.ItemTime{Time: testDate},
		Tags: []string{
			"programming",
			"go",
		},
	}
	indexStr := delegate.renderIndex(styles.index, 1, true)
	titleStr := delegate.renderTitle(styles.title, item.Title)
	tagsStr := delegate.renderTags(styles.tags, styles.tag, item.Tags)
	dateStr := delegate.renderDate(styles.date, &item.CreatedAt.Time)
	expected := fmt.Sprintf("%s %s\n%[3]*s%s %s", indexStr, titleStr, itemPrefixLength-styles.mainLine.GetPaddingLeft(), "", dateStr, tagsStr)

	str := delegate.renderItem(styles, item, 1, true)

	assert.Equal(t, styles.mainLine.Render(expected), str)
}

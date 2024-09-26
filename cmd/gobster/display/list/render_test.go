package list

import (
	"fmt"
	"strings"
	"testing"
	"time"

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

func TestRenderCategories(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	categories := []string{"programming", "go"}
	fmtCat := []string{}
	for _, cat := range categories {
		fmtCat = append(fmtCat, delegate.renderCategory(styles.category, cat))
	}

	str := delegate.renderCategories(styles.categories, styles.category, []string{"programming", "go"})

	assert.Equal(t, styles.categories.Render(strings.Join(fmtCat, " ")), str)
}

func TestRenderCategory(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	cat := "programming"
	catSpecificForeground, exists := catBackgrounds[cat]
	if !exists {
		catSpecificForeground = tagDefaultBackground
	}

	str := delegate.renderCategory(styles.category, cat)
	expectedCatStyle := styles.category.Foreground(catSpecificForeground)

	assert.Equal(t, expectedCatStyle.Render(fmt.Sprintf("<%s>", cat)), str)
}

func TestRenderItem(t *testing.T) {
	styles := newStyleProvider(false)
	delegate := itemDelegate{}
	item := Item{
		title: "Title",
		categories: []string{
			"programming",
			"go",
		},
		date: &testDate,
	}
	indexStr := delegate.renderIndex(styles.index, 1, false)
	titleStr := delegate.renderTitle(styles.title, item.title)
	categoriesStr := delegate.renderCategories(styles.categories, styles.category, item.categories)
	dateStr := delegate.renderDate(styles.date, item.date)
	expected := fmt.Sprintf("%s %s %s\n%[4]*s%s", indexStr, titleStr, categoriesStr, itemPrefixLength-styles.mainLine.GetPaddingLeft(), "", dateStr)

	str := delegate.renderItem(styles, item, 1, false)

	assert.Equal(t, styles.mainLine.Render(expected), str)
}

func TestRenderItemSelected(t *testing.T) {
	styles := newStyleProvider(true)
	delegate := itemDelegate{}
	item := Item{
		title: "Title",
		categories: []string{
			"programming",
			"go",
		},
		date: &testDate,
	}
	indexStr := delegate.renderIndex(styles.index, 1, true)
	titleStr := delegate.renderTitle(styles.title, item.title)
	categoriesStr := delegate.renderCategories(styles.categories, styles.category, item.categories)
	dateStr := delegate.renderDate(styles.date, item.date)
	expected := fmt.Sprintf("%s %s %s\n%[4]*s%s", indexStr, titleStr, categoriesStr, itemPrefixLength-styles.mainLine.GetPaddingLeft(), "", dateStr)

	str := delegate.renderItem(styles, item, 1, true)

	assert.Equal(t, styles.mainLine.Render(expected), str)
}

package list

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	testDate = time.Date(2024, 1, 1, 14, 30, 0, 0, time.UTC)
)

func TestRenderIndex(t *testing.T) {
	delegate := itemDelegate{}

	str := delegate.renderIndex(3, false)

	assert.Equal(t, " 4.", str)
}

func TestRenderTitle(t *testing.T) {
	delegate := itemDelegate{}

	str := delegate.renderTitle("Title", true)

	assert.Equal(t, "Title", str)
}

func TestRenderCategory(t *testing.T) {
	delegate := itemDelegate{}

	str := delegate.renderCategory("programming")

	assert.Equal(t, "<programming>", str)
}

func TestRenderCategories(t *testing.T) {
	delegate := itemDelegate{}

	str := delegate.renderCategories([]string{"programming", "go"})

	assert.Equal(t, "<programming> <go>", str)
}

func TestRenderMainLine(t *testing.T) {
	delegate := itemDelegate{}
	item := Item{
		title: "Title",
		categories: []string{
			"programming",
			"go",
		},
		date: &testDate,
	}

	str := delegate.renderMainLine(item, 1, false)
	strSelected := delegate.renderMainLine(item, 1, true)

	assert.Equal(t, "     2. Title <programming> <go>\n        Mon  1 Jan 2024 - 14:30 ", str)
	assert.Equal(t, "  ➜  2. Title <programming> <go>\n        Mon  1 Jan 2024 - 14:30 ", strSelected)
}

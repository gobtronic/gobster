package display

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func TestNewModel(t *testing.T) {
	feed := gofeed.Feed{
		Items: []*gofeed.Item{
			{Title: "Discussion 1"},
			{Title: "Discussion 2"},
			{Title: "Discussion 3"},
		},
	}

	model := NewModel(&feed)

	assert.Len(t, model.list.Items(), 3)
}

func TestUpdate(t *testing.T) {
	feed := gofeed.Feed{
		Items: []*gofeed.Item{
			{Title: "Discussion 1"},
			{Title: "Discussion 2"},
			{Title: "Discussion 3"},
		},
	}
	m := NewModel(&feed)

	updatedM, _ := m.Update(tea.WindowSizeMsg{Width: 10, Height: 10})

	assert.Equal(t, updatedM.(model).list.Height(), 10)
	assert.Equal(t, updatedM.(model).list.Width(), 10)
}

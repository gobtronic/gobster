package display

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/assert"
)

func feedMock() gofeed.Feed {
	return gofeed.Feed{
		Items: []*gofeed.Item{
			{Title: "Discussion 1"},
			{Title: "Discussion 2"},
			{Title: "Discussion 3"},
		},
	}
}

func TestNewModel(t *testing.T) {
	feed := feedMock()
	termSize := [2]int{60, 30}

	model := NewModel(&feed, termSize)

	assert.Len(t, model.list.Items(), 3)
	assert.Equal(t, model.list.Width(), termSize[0])
	assert.Equal(t, model.list.Height(), termSize[1])
}

func TestUpdateWindowSizeMsg(t *testing.T) {
	feed := feedMock()
	m := NewModel(&feed, [2]int{30, 30})

	updatedM, cmd := m.Update(tea.WindowSizeMsg{Width: 10, Height: 10})

	assert.Equal(t, updatedM.(model).list.Height(), 10)
	assert.Equal(t, updatedM.(model).list.Width(), 10)
	assert.Nil(t, cmd)
}

func TestUpdateQuitKeyMsg(t *testing.T) {
	feed := feedMock()
	m := NewModel(&feed, [2]int{30, 30})

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}})

	assert.IsType(t, tea.QuitMsg{}, cmd().(tea.QuitMsg))
}

func TestUpdateUnknownKeyMsg(t *testing.T) {
	feed := feedMock()
	m := NewModel(&feed, [2]int{30, 30})

	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'a'}})

	assert.Nil(t, cmd)
}

package display

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gobtronic/gobster/cmd/gobster/feed"
	"github.com/stretchr/testify/assert"
)

func feedMock() feed.LobsterFeed {
	return feed.LobsterFeed{
		feed.Item{
			Title:     "Title",
			CreatedAt: feed.ItemTime{Time: time.Now()},
			Tags:      []string{"programming", "go"},
		},
		feed.Item{
			Title:     "Title",
			CreatedAt: feed.ItemTime{Time: time.Now()},
			Tags:      []string{"programming", "go"},
		},
		feed.Item{
			Title:     "Title",
			CreatedAt: feed.ItemTime{Time: time.Now()},
			Tags:      []string{"programming", "go"},
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

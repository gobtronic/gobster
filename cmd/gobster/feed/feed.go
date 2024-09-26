package feed

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type feedType int

const (
	Active feedType = iota
)

type LobsterFeed []Item
type Item struct {
	ShortId       string   `json:"short_id"`
	ShortIdUrl    string   `json:"short_id_url"`
	CreatedAt     ItemTime `json:"created_at"`
	Title         string
	Url           string
	Score         int
	CommentCount  int    `json:"comment_count"`
	SubmitterUser string `json:"submitted_user"`
	UserIsAuthor  bool   `json:"user_is_author"`
	Tags          []string
}
type ItemTime struct {
	time.Time
}

func (d *ItemTime) UnmarshalJSON(data []byte) error {
	t, err := time.Parse(`"2006-01-02T15:04:05.000-07:00"`, string(data))
	if err != nil {
		return err
	}
	d.Time = t
	return nil
}

func (i Item) FilterValue() string { return i.Title + strings.Join(i.Tags, " ") }

var ErrUnknownFeedType = errors.New("Unknown feed type")

type HTTPErr struct {
	statusCode int
}

func (err HTTPErr) Error() string {
	return fmt.Sprintf("error statusCode: %d", err.statusCode)
}

func FetchFeed(ft feedType) (LobsterFeed, error) {
	feed := LobsterFeed{}
	var jsonName string
	switch ft {
	case Active:
		jsonName = "active"
	default:
		return feed, ErrUnknownFeedType
	}

	client := http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("https://lobste.rs/%s.json", jsonName), nil)
	if err != nil {
		return feed, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return feed, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, HTTPErr{resp.StatusCode}
	}

	err = json.NewDecoder(resp.Body).Decode(&feed)
	if err != nil {
		return feed, err
	}

	return feed, nil
}

package feed

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestFeedUnmarshal(t *testing.T) {
	itemJSON := `[{"short_id":"pca7h0","short_id_url":"https://lobste.rs/s/pca7h0","created_at":"2024-09-25T15:27:23.000-05:00","title":"Eliminating memory safety vulnerabilities at the source","url":"https://security.googleblog.com/2024/09/eliminating-memory-safety-vulnerabilities-Android.html","score":43,"flags":0,"comment_count":29,"description":"","description_plain":"","comments_url":"https://lobste.rs/s/pca7h0/eliminating_memory_safety","submitter_user":"heinrich5991","user_is_author":false,"tags":["android","security"]},{"short_id":"29a1eo","short_id_url":"https://lobste.rs/s/29a1eo","created_at":"2024-09-26T01:22:35.000-05:00","title":"Rewriting Rust","url":"https://josephg.com/blog/rewriting-rust/","score":72,"flags":0,"comment_count":43,"description":"","description_plain":"","comments_url":"https://lobste.rs/s/29a1eo/rewriting_rust","submitter_user":"mpweiher","user_is_author":false,"tags":["rust"]}]`
	feed := LobsterFeed{}
	expItemCreatedAtTime := time.Date(2024, time.September, 25, 15, 27, 23, 0, time.FixedZone("", -60*60*5))

	err := json.Unmarshal([]byte(itemJSON), &feed)
	if len(feed) == 0 {
		assert.FailNow(t, "feed should not be empty")
	}
	firstItem := feed[0]

	assert.Nil(t, err)
	assert.Equal(t, 2, len(feed))
	assert.Equal(t, "Eliminating memory safety vulnerabilities at the source", firstItem.Title)
	assert.Equal(t, 29, firstItem.CommentCount)
	assert.Equal(t, "https://security.googleblog.com/2024/09/eliminating-memory-safety-vulnerabilities-Android.html", firstItem.Url)
	assert.Equal(t, []string{"android", "security"}, firstItem.Tags)
	assert.Equal(t, 43, firstItem.Score)
	assert.Equal(t, expItemCreatedAtTime, firstItem.CreatedAt.Time)
}

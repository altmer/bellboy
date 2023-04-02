package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddText(t *testing.T) {
	teardown := setup()
	defer teardown()

	testCases := []struct {
		text Text
		err  error
	}{
		{Text{Title: "my new text", Body: "story text", PostID: 1}, nil},
		{Text{Title: "another text", Body: "yet another awesome story", PostID: 2}, nil},
	}

	for _, testCase := range testCases {
		err := repo.AddText(&testCase.text)
		checkErrors(t, testCase.err, err)
		if err == nil {
			var dbText Text
			DB.Get(
				&dbText,
				"SELECT id, created_at, post_id, title, body FROM texts WHERE post_id = ? LIMIT 1",
				testCase.text.PostID,
			)
			assert.NotEqual(t, uint(0), testCase.text.ID, "ID is not set for saved text")
			assert.NotEqual(t, uint(0), dbText.ID, "ID is not fetched from database")
			assert.Equal(t, testCase.text.PostID, dbText.PostID)
			assert.Equal(t, testCase.text.Title, dbText.Title)
			assert.Equal(t, testCase.text.Body, dbText.Body)
		}
	}
}

package media

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddLink(t *testing.T) {
	teardown := setup()
	defer teardown()

	testCases := []struct {
		link Link
		err  error
	}{
		{Link{URL: "http://google.com", PostID: 21434332}, nil},
		{Link{URL: "http://medium.com", PostID: 23432}, nil},
	}

	for _, testCase := range testCases {
		err := repo.AddLink(&testCase.link)
		checkErrors(t, testCase.err, err)
		if err == nil {
			var dbLink Link
			DB.Get(
				&dbLink,
				"SELECT id, created_at, post_id, url FROM links WHERE post_id = ? LIMIT 1",
				testCase.link.PostID,
			)
			assert.NotEqual(t, uint(0), testCase.link.ID, "ID is not set for saved link")
			assert.NotEqual(t, uint(0), dbLink.ID, "ID is not fetched from database")
			assert.Equal(t, testCase.link.PostID, dbLink.PostID)
			assert.Equal(t, testCase.link.URL, dbLink.URL)
		}
	}
}

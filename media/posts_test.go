package media

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPost(t *testing.T) {
	teardown := setup()
	defer teardown()

	testCases := []struct {
		post Post
		err  error
	}{
		{Post{
			ExternalID: "1", Status: "approved", SFW: true,
			Source: "blogName", Type: "text", ExternalURL: "http://blog.com",
		}, nil},
		{Post{ExternalID: "1"}, errors.New("UNIQUE constraint failed: posts.external_id")},
		{Post{
			ExternalID: "2", Status: "queued", SFW: false,
			Source: "another", Type: "video", ExternalURL: "http://blogee.com",
		}, nil},
	}

	for _, testCase := range testCases {
		err := repo.AddPost(&testCase.post)
		checkErrors(t, testCase.err, err)
		if err == nil {
			var dbPost Post
			DB.Get(
				&dbPost,
				"SELECT id, external_id, created_at, status, sfw, type, category, source, external_url FROM posts WHERE external_id = ? LIMIT 1",
				testCase.post.ExternalID,
			)

			assert.NotEqual(t, uint(0), testCase.post.ID, "ID is not set for saved object")
			assert.NotEqual(t, uint(0), dbPost.ID, "ID is not fetched from DB")
			assert.Equal(t, testCase.post.ExternalID, dbPost.ExternalID)
			assert.Equal(t, testCase.post.Status, dbPost.Status)
			assert.Equal(t, testCase.post.SFW, dbPost.SFW)
			assert.Equal(t, testCase.post.Source, dbPost.Source)
			assert.Equal(t, testCase.post.Type, dbPost.Type)
			assert.Equal(t, testCase.post.ExternalURL, dbPost.ExternalURL)
		}
	}
}

func TestPostExistsWithExternalID(t *testing.T) {
	teardown := setup()
	defer teardown()

	post := Post{ExternalID: "5"}
	err := repo.AddPost(&post)
	if err != nil {
		panic(err)
	}

	testCases := []struct {
		externalID string
		expected   bool
	}{
		{"5", true},
		{"32", false},
	}

	for _, testCase := range testCases {
		actual := repo.PostExistsWithExternalID(testCase.externalID)
		if actual != testCase.expected {
			t.Errorf("Expected result for [%#v] was [%#v], got [%#v]",
				testCase.externalID, testCase.expected, actual)
		}
	}
}

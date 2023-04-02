package media

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func checkPostHasTag(t *testing.T, post *Post, tagName string) {
	tag := Tag{}
	DB.Get(
		&tag,
		"SELECT tags.name, tags.id FROM tags INNER JOIN posts_tags ON tags.id = posts_tags.tag_id WHERE posts_tags.post_id = ? AND tags.name = ?",
		post.ID, tagName,
	)
	assert.Equal(t, tagName, tag.Name, "post does not have expected tag")
	assert.NotEqual(t, uint(0), tag.ID, "tag ID can not be 0")
}

func TestAddTagToPost(t *testing.T) {
	teardown := setup()
	defer teardown()

	post := &Post{
		Status:     "added",
		SFW:        false,
		Source:     "tumblr",
		Category:   "blog",
		ExternalID: "432432432432",
	}
	repo.AddPost(post)

	repo.AddTag(&Tag{Name: "existing-tag"})

	// check tags count
	var tagsCount int
	DB.Get(&tagsCount, "SELECT count(*) FROM tags")
	if tagsCount != 1 {
		t.Errorf("Expected tags count to be [1], got [%d]", tagsCount)
	}

	testCases := []struct {
		externalTag string
		err         error
	}{
		{"new-tag", nil},
		{"existing-tag", nil},
		{"another-tag", nil},
		{"new-tag", errors.New("UNIQUE constraint failed: posts_tags.post_id, posts_tags.tag_id")},
	}

	for _, testCase := range testCases {
		err := repo.AddTagToPost(post, testCase.externalTag)
		checkErrors(t, testCase.err, err)
		if err == nil {
			checkPostHasTag(t, post, testCase.externalTag)
		}
	}
	// check tags count
	DB.Get(&tagsCount, "SELECT count(*) FROM tags")
	if tagsCount != 3 {
		t.Errorf("Expected tags count to be [2], got [%d]", tagsCount)
	}
}

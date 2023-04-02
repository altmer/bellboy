package tumblr

import (
	"github.com/altmer/bellboy/media"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubsDown(t *testing.T) {
	mock := mockClient{
		DeletedPosts: []int{},
		UnlikedPosts: []int{},
	}

	teardown := setup()
	defer teardown()

	repo.AddSubscription(&media.Subscription{BlogName: "a"})
	repo.AddSubscription(&media.Subscription{BlogName: "b"})
	repo.AddSubscription(&media.Subscription{BlogName: "c"})
	repo.AddSubscription(&media.Subscription{BlogName: "d"})
	repo.AddSubscription(&media.Subscription{BlogName: "e"})
	var subsCount int
	DB.Get(&subsCount, "SELECT count(*) FROM subscriptions")

	assert.Equal(t, 5, subsCount)

	Syncer{
		BlogName: "blog_with_posts",
		Client:   &mock,
		Repo:     repo,
	}.SubsDown()

	expectedSubsCount := len(userFollowing.Blogs)

	DB.Get(&subsCount, "SELECT count(*) FROM subscriptions")

	if subsCount != expectedSubsCount {
		t.Errorf("Expected subscriptions count to be [%d], got [%d]", expectedSubsCount, subsCount)
	}

	var subscriptions []media.Subscription
	DB.Select(&subscriptions, "SELECT id, created_at, updated_at, url, blog_name, source, title, description FROM subscriptions")

	sub := subscriptions[0]
	assert.NotEqual(t, uint(0), sub.ID)
	assert.Equal(t, "http://tumblr.com/blog1", sub.URL)
	assert.Equal(t, "blog1", sub.BlogName)
	assert.Equal(t, "title1", sub.Title)
	assert.Equal(t, "description1", sub.Description)

	sub = subscriptions[1]
	assert.NotEqual(t, uint(0), sub.ID)
	assert.Equal(t, "http://tumblr.com/blog2", sub.URL)
	assert.Equal(t, "blog2", sub.BlogName)
	assert.Equal(t, "title2", sub.Title)
	assert.Equal(t, "description2", sub.Description)
}

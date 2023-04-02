package tumblr

import (
	"github.com/altmer/bellboy/media"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSubsUp(t *testing.T) {
	mock := mockClient{
		DeletedPosts:  []int{},
		UnlikedPosts:  []int{},
		FollowedBlogs: []string{},
	}

	teardown := setup()
	defer teardown()

	repo.AddSubscription(&media.Subscription{BlogName: "a", URL: "http://tumblr.com/a"})
	repo.AddSubscription(&media.Subscription{BlogName: "b", URL: "http://tumblr.com/b"})
	repo.AddSubscription(&media.Subscription{BlogName: "c", URL: "http://tumblr.com/c"})
	repo.AddSubscription(&media.Subscription{BlogName: "d", URL: "http://tumblr.com/d"})
	repo.AddSubscription(&media.Subscription{BlogName: "e", URL: "http://tumblr.com/e"})

	Syncer{
		BlogName: "blog_with_posts",
		Client:   &mock,
		Repo:     repo,
	}.SubsUp()

	assert.Equal(t, 5, len(mock.FollowedBlogs))
	assert.Equal(t, "http://tumblr.com/a", mock.FollowedBlogs[0])
	assert.Equal(t, "http://tumblr.com/b", mock.FollowedBlogs[1])
	assert.Equal(t, "http://tumblr.com/c", mock.FollowedBlogs[2])
	assert.Equal(t, "http://tumblr.com/d", mock.FollowedBlogs[3])
	assert.Equal(t, "http://tumblr.com/e", mock.FollowedBlogs[4])
}

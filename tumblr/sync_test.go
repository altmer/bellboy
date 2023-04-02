package tumblr

import (
	"io/ioutil"
	"testing"

	"github.com/altmer/bellboy/media"

	"github.com/stretchr/testify/assert"
)

func TestSync(t *testing.T) {

	mock := mockClient{
		DeletedPosts: []int{},
		UnlikedPosts: []int{},
	}

	teardown := setup()

	Syncer{
		BlogName: "blog_with_posts",
		Client:   &mock,
		Repo:     repo,
	}.Sync()

	expectedPostsCount := len(blogPosts.Posts) + len(likes.LikedPost)

	var postsCount int
	DB.Get(&postsCount, "SELECT count(*) FROM posts")

	if postsCount != expectedPostsCount {
		t.Errorf("Expected posts count to be [%d], got [%d]", expectedPostsCount, postsCount)
	}

	var posts []media.Post
	DB.Select(&posts, "SELECT id, type, external_id, category, external_id, external_url, source_category, likes, type, summary FROM posts")

	post := posts[0]
	assert.NotEqual(t, uint(0), post.ID)
	assert.Equal(t, "link", post.Type)
	assert.Equal(t, "10", post.ExternalID)
	assert.Equal(t, "https://tumblr.com/posts/1", post.ExternalURL)
	assert.Equal(t, "anotherlinksblog", post.SourceCategory)
	assert.Equal(t, 23, post.Likes)
	assert.Equal(t, "this is link", post.Summary)

	var link media.Link
	DB.Get(&link, "SELECT post_id, url FROM links WHERE post_id = ?", post.ID)

	assert.Equal(t, post.ID, link.PostID, "link.PostId is wrong")
	assert.Equal(t, "http://example.com/content", link.URL, "link url is wrong")

	post = posts[1]
	assert.NotEqual(t, uint(0), post.ID)
	assert.Equal(t, "text", post.Type)
	assert.Equal(t, "12", post.ExternalID)
	assert.Equal(t, "https://tumblr.com/posts/2", post.ExternalURL)
	assert.Equal(t, "anothertextsblog", post.SourceCategory)
	assert.Equal(t, 4357, post.Likes)
	assert.Equal(t, "great novel", post.Summary)

	var text media.Text
	DB.Get(&text, "SELECT post_id, title, body FROM texts WHERE post_id = ?", post.ID)

	assert.Equal(t, post.ID, text.PostID, "text.PostId is wrong")
	assert.Equal(t, "Novel title", text.Title, "text title is wrong")
	assert.Equal(t, "Novel body", text.Body, "text body is wrong")

	post = posts[2]
	assert.NotEqual(t, uint(0), post.ID)
	assert.Equal(t, "photo", post.Type)
	assert.Equal(t, "14", post.ExternalID)
	assert.Equal(t, "https://tumblr.com/posts/3", post.ExternalURL)
	assert.Equal(t, "photocoolblog", post.SourceCategory)
	assert.Equal(t, 22, post.Likes)
	assert.Equal(t, "photo exhibition", post.Summary)

	var photo media.Photo
	DB.Get(&photo, "SELECT post_id, caption, external_url, sfw FROM photos WHERE post_id = ?", post.ID)

	assert.Equal(t, post.ID, photo.PostID, "photo.PostId is wrong")
	assert.Equal(t, false, photo.SFW, "photo.SFW is wrong")
	assert.Equal(t, "photo inner caption", photo.Caption, "photo.Caption is wrong")
	assert.Equal(t, "http://photo.tumblr/photo.png", photo.ExternalURL, "photo.ExternalURL is wrong")
	var photoContents []byte
	photoContents, _ = ioutil.ReadFile("./photo_1.png")
	if string(photoContents) != "png file contents" {
		t.Errorf("Wrong photo file content: [%s]", photoContents)
	}

	post = posts[3]
	assert.NotEqual(t, uint(0), post.ID)
	assert.Equal(t, "video", post.Type)
	assert.Equal(t, "18", post.ExternalID)
	assert.Equal(t, "https://tumblr.com/posts/4", post.ExternalURL)
	assert.Equal(t, 0, post.Likes)
	assert.Equal(t, "movie", post.Summary)

	var video media.Video
	DB.Get(&video, "SELECT post_id, external_url, thumbnail_url FROM videos WHERE post_id = ?", post.ID)

	assert.Equal(t, post.ID, video.PostID, "video.PostId is wrong")
	assert.Equal(t, "http://photo.tumblr/video_thumb.png", video.ThumbnailURL, "video.Caption is wrong")
	assert.Equal(t, "http://photo.tumblr/video.mp4", video.ExternalURL, "video.ExternalURL is wrong")
	var videoContents, videoThumbnailContents []byte
	videoContents, _ = ioutil.ReadFile("./video_1.mp4")
	videoThumbnailContents, _ = ioutil.ReadFile("./video_1_thumbnail.png")
	if string(videoContents) != "mp4 file contents" {
		t.Errorf("Wrong video file content: [%s]", videoContents)
	}
	if string(videoThumbnailContents) != "thumbnail file contents" {
		t.Errorf("Wrong video thumbnail file content: [%s]", videoThumbnailContents)
	}

	teardown()
}

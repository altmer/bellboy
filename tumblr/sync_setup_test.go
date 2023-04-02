package tumblr

import (
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"os"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	"github.com/altmer/bellboy/context"
	"github.com/altmer/bellboy/media"
)

var blogPosts = BlogPosts{
	TotalPosts: 2,
	Posts: []Post{
		Post{
			ID:          10,
			Type:        "link",
			BlogName:    "linksblog",
			PostURL:     "https://tumblr.com/posts/1",
			SourceTitle: "anotherlinksblog",
			NoteCount:   23,
			Date:        "2017-01-02 12:33:44 CET",
			Summary:     "this is link",
			URL:         "http://example.com/content",
		},
		Post{
			ID:          12,
			Type:        "text",
			BlogName:    "textsblog",
			PostURL:     "https://tumblr.com/posts/2",
			SourceTitle: "anothertextsblog",
			NoteCount:   4357,
			Date:        "2017-06-01 09:33:44 CET",
			Summary:     "great novel",
			Title:       "Novel title",
			Body:        "Novel body",
		},
	},
}

var likes = Likes{
	LikedCount: 2,
	LikedPost: []Post{
		Post{
			ID:          14,
			Type:        "photo",
			BlogName:    "photoblog",
			PostURL:     "https://tumblr.com/posts/3",
			SourceTitle: "photocoolblog",
			NoteCount:   22,
			Date:        "2017-06-01 09:33:44 CET",
			Summary:     "photo exhibition",
			Caption:     "photo caption",
			Photos: []struct {
				Caption      string `json:"caption,omitempty"` // user supplied caption for the individual photo
				OriginalSize struct {
					Height int    `json:"height,omitempty"` // height of the image
					Width  int    `json:"width,omitempty"`  // width of the image
					URL    string `json:"url,omitempty"`    // location of the photo file
				} `json:"original_size,omitempty"`
				AlternateSizes []struct {
					Height int    `json:"height,omitempty"` // height of the photo
					Width  int    `json:"width,omitempty"`  // width of the photo
					URL    string `json:"url,omitempty"`    // Location of the photo file
				} `json:"alt_sizes,omitempty"` // alternate photo sizes
			}{
				{
					Caption: "photo inner caption",
					OriginalSize: struct {
						Height int    `json:"height,omitempty"` // height of the image
						Width  int    `json:"width,omitempty"`  // width of the image
						URL    string `json:"url,omitempty"`    // location of the photo file
					}{
						Height: 100,
						Width:  100,
						URL:    "http://photo.tumblr/photo.png",
					},
				},
			},
		},
		Post{
			ID:           18,
			Type:         "video",
			BlogName:     "videoblog",
			PostURL:      "https://tumblr.com/posts/4",
			NoteCount:    0,
			Date:         "2017-06-01 09:33:44 CET",
			Summary:      "movie",
			VideoURL:     "http://photo.tumblr/video.mp4",
			ThumbnailURL: "http://photo.tumblr/video_thumb.png",
		},
	},
}

var userFollowing = UserFollowing{
	TotalBlogs: 2,
	Blogs: []struct {
		Name        string `json:"name"`        // the user name attached to the blog that's being followed
		URL         string `json:"url"`         // the URL of the blog that's being followed
		Updated     int    `json:"updated"`     // the time of the most recent post, in seconds since the epoch
		Title       string `json:"title"`       // the title of the blog
		Description string `json:"description"` // the description of the blog
	}{
		{
			Name:        "blog1",
			URL:         "http://tumblr.com/blog1",
			Description: "description1",
			Title:       "title1",
		},
		{
			Name:        "blog2",
			URL:         "http://tumblr.com/blog2",
			Description: "description2",
			Title:       "title2",
		},
	},
}

type mockClient struct {
	UnlikedPosts  []int
	DeletedPosts  []int
	FollowedBlogs []string
}

func (client mockClient) BlogPosts(blogName string, params map[string]string) BlogPosts {
	return blogPosts
}

func (client *mockClient) PostDelete(blogName string, postID int) Meta {
	client.DeletedPosts = append(client.DeletedPosts, postID)
	return Meta{}
}

func (client mockClient) UserLikes(params map[string]string) Likes {
	return likes
}

func (client *mockClient) UserUnlike(postID int, reblogKey string) Meta {
	client.UnlikedPosts = append(client.UnlikedPosts, postID)
	return Meta{}
}

func (client mockClient) UserFollowing(params map[string]string) UserFollowing {
	return userFollowing
}

func (client *mockClient) UserFollow(followURL string) Meta {
	client.FollowedBlogs = append(client.FollowedBlogs, followURL)
	return Meta{}
}

var DB *sqlx.DB
var repo media.Repository

func setup() func() {
	httpmock.Activate()

	httpmock.RegisterResponder("GET", "http://photo.tumblr/photo.png",
		httpmock.NewStringResponder(200, "png file contents"))
	httpmock.RegisterResponder("GET", "http://photo.tumblr/video.mp4",
		httpmock.NewStringResponder(200, "mp4 file contents"))
	httpmock.RegisterResponder("GET", "http://photo.tumblr/video_thumb.png",
		httpmock.NewStringResponder(200, "thumbnail file contents"))

	testDBPath := "./test.db"
	testMediaPath := "./"

	viper.SetDefault("media_folder", testMediaPath)
	DB = context.NewDBConnection(testDBPath)

	repo = media.NewRepository(DB)

	return func() {
		DB.Close()
		os.Remove(testDBPath)
		os.Remove("./photo_1.png")
		os.Remove("./video_1.mp4")
		os.Remove("./video_1_thumbnail.png")
		httpmock.DeactivateAndReset()
	}
}

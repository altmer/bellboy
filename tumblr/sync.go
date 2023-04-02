package tumblr

import (
	"fmt"
	"strconv"
	"time"

	"github.com/altmer/bellboy/media"
)

// Syncer represents the main syncing entity
type Syncer struct {
	BlogName string
	Client   API
	Repo     media.Repository
}

// Sync syncs tumblr blog with given config
func (s Syncer) Sync() {
	fmt.Printf("Getting info from Tumblr blog [%s]\n", s.BlogName)
	totalPosts := s.Client.BlogPosts(s.BlogName, map[string]string{}).TotalPosts
	fmt.Printf("%d posts found\n", totalPosts)
	s.syncBlogPosts(totalPosts)
	totalLikes := s.Client.UserLikes(map[string]string{}).LikedCount
	fmt.Printf("%d user likes found\n", totalLikes)
	s.syncLikes(totalLikes)
	fmt.Println("Tumblr synced!")
}

// SubsDown imports subscriptions from tumblr
func (s Syncer) SubsDown() {
	totalSubscriptions := s.Client.UserFollowing(map[string]string{}).TotalBlogs
	fmt.Printf("%d user subscriptions found\n", totalSubscriptions)

	s.Repo.RemoveAllSubscriptions()

	limit := 20
	for offset := 0; offset < totalSubscriptions; offset += limit {
		fmt.Printf("Fetching subscriptions from [%d] to [%d]...\n", offset, offset+limit)
		subscriptions := s.Client.UserFollowing(map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  strconv.Itoa(limit),
		})
		for _, blog := range subscriptions.Blogs {
			s.Repo.AddSubscription(&media.Subscription{
				BlogName:    blog.Name,
				URL:         blog.URL,
				Description: blog.Description,
				Title:       blog.Title,
				Source:      "tumblr",
			})
		}
	}
}

// SubsUp exports subscriptions to tumblr blog (follows all of them)
func (s Syncer) SubsUp() {
	fmt.Println("Exporting subscriptions...")
	subs, err := s.Repo.ListSubscriptions()
	if err != nil {
		panic(err)
	}
	for _, sub := range subs {
		s.Client.UserFollow(sub.URL)
	}
}

func (s Syncer) syncBlogPosts(totalPosts int) {
	limit := 20
	for offset := 0; offset < totalPosts; offset += limit {
		fmt.Printf("Fetching posts from [%d] to [%d]...\n", offset, offset+limit)
		posts := s.Client.BlogPosts(s.BlogName, map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  strconv.Itoa(limit),
		})
		for _, post := range posts.Posts {
			result := s.syncBlogPost(&post, "added")
			if result {
				s.Client.PostDelete(s.BlogName, post.ID)
			}
		}
	}
}

func (s Syncer) syncLikes(totalLikes int) {
	limit := 20
	for offset := 0; offset < totalLikes; offset += limit {
		fmt.Printf("Fetching likes from [%d] to [%d]...\n", offset, offset+limit)
		likes := s.Client.UserLikes(map[string]string{
			"offset": strconv.Itoa(offset),
			"limit":  strconv.Itoa(limit),
		})
		for _, post := range likes.LikedPost {
			result := s.syncBlogPost(&post, "queued")
			if result {
				s.Client.UserUnlike(post.ID, post.ReblogKey)
			}
		}
	}
}

func (s Syncer) syncBlogPost(externalPost *Post, state string) bool {
	if s.Repo.PostExistsWithExternalID(strconv.Itoa(externalPost.ID)) {
		fmt.Printf("Post with id [%d] already exists\n", externalPost.ID)
		return false
	}

	post, err := createPost(externalPost)
	if err != nil {
		fmt.Printf("WARN: Post loading failed with error [%s] for post [%#v]\n", err, post)
		return false
	}
	post.Status = state
	err = s.Repo.AddPost(post)
	if err != nil {
		fmt.Printf("WARN: Post creation failed with error [%s] for post [%#v]\n", err, post)
		return false
	}

	switch post.Type {
	case "link":
		link := createLink(post, externalPost)
		err = s.Repo.AddLink(link)
		if err != nil {
			fmt.Printf("WARN: Link creation failed with error [%s] for link [%#v]", err, link)
			return false
		}
	case "text":
		text := createText(post, externalPost)
		err = s.Repo.AddText(text)
		if err != nil {
			fmt.Printf("WARN: Text creation failed with error [%s] for text [%#v]", err, text)
			return false
		}
	case "photo":
		for _, externalPhoto := range externalPost.Photos {
			photo := createPhoto(post, externalPhoto.OriginalSize.URL, externalPhoto.Caption)
			err = s.Repo.AddPhoto(photo)
			if err != nil {
				fmt.Printf("WARN: Photo creation failed with error [%s] for photo [%#v]", err, photo)
				return false
			}
		}
	case "video":
		video := createVideo(post, externalPost)
		err = s.Repo.AddVideo(video)
		if err != nil {
			fmt.Printf("WARN: Video creation failed with error [%s] for video [%#v]", err, video)
			return false
		}
	default:
		fmt.Printf("WARN: Unexpected post type: [%s]", post.Type)
		return false
	}

	// we load tags only for finished posts
	if post.Status != "added" {
		return true
	}

	for _, externalTag := range externalPost.Tags {
		s.Repo.AddTagToPost(post, externalTag)
	}
	return true
}

func createPost(externalPost *Post) (*media.Post, error) {
	releasedAt, err := time.Parse("2006-01-02 15:04:05 MST", externalPost.Date)
	if err != nil {
		return nil, err
	}
	return &media.Post{
		SFW:            false,
		Source:         "tumblr",
		Type:           externalPost.Type,
		Category:       externalPost.BlogName,
		ExternalID:     strconv.Itoa(externalPost.ID),
		ExternalURL:    externalPost.PostURL,
		SourceURL:      externalPost.SourceURL,
		SourceCategory: externalPost.SourceTitle,
		Likes:          externalPost.NoteCount,
		ReleasedAt:     releasedAt,
		Summary:        externalPost.Summary,
	}, nil
}

func createLink(post *media.Post, externalPost *Post) *media.Link {
	return &media.Link{PostID: post.ID, URL: externalPost.URL}
}

func createText(post *media.Post, externalPost *Post) *media.Text {
	return &media.Text{PostID: post.ID, Title: externalPost.Title, Body: externalPost.Body}
}

func createVideo(post *media.Post, externalPost *Post) *media.Video {
	return &media.Video{
		PostID:       post.ID,
		ExternalURL:  externalPost.VideoURL,
		ThumbnailURL: externalPost.ThumbnailURL,
	}
}

func createPhoto(post *media.Post, url, caption string) *media.Photo {
	return &media.Photo{
		PostID:      post.ID,
		Caption:     caption,
		ExternalURL: url,
		SFW:         false,
	}
}

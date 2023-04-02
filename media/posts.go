package media

import (
	"time"
)

func (r mediaRepo) PostExistsWithExternalID(externalID string) bool {
	var res int
	r.DB.Get(&res, "SELECT count(*) FROM posts WHERE external_id = ?", externalID)
	return res > 0
}

func (r mediaRepo) AddPost(post *Post) error {
	post.CreatedAt = time.Now()
	post.UpdatedAt = time.Now()
	res, err := r.DB.NamedExec(
		`INSERT INTO posts (
       created_at, updated_at, status, sfw, source, type, released_at, category, external_id,
       external_url, source_url, source_category, likes, summary)
     VALUES (
       :created_at, :updated_at, :status, :sfw, :source, :type, :released_at, :category,
       :external_id, :external_url, :source_url, :source_category, :likes, :summary
     )`,
		post,
	)
	if err != nil {
		return err
	}
	postID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	post.ID = uint(postID)
	return nil
}

// Post represents one post entity (tumblr post, fanfic, story, deviantart post).
type Post struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Status     string    // could be one of [approved queued]
	SFW        bool      // safe for work
	Source     string    // source of the post
	Type       string    // text, photo, video, link
	ReleasedAt time.Time `db:"released_at"` // when post was released initially

	Category       string // blog name or author
	ExternalID     string `db:"external_id"`     // ID in the external domain (f.ex. tumblr id)
	ExternalURL    string `db:"external_url"`    // URL to the post in the internet
	SourceURL      string `db:"source_url"`      // original post URL
	SourceCategory string `db:"source_category"` // original post author
	Likes          int    // number of likes
	Summary        string // caption to the post
}

// PostsSchema represents schema for "posts" table
var PostsSchema = `CREATE TABLE "posts" (
	"id" integer primary key autoincrement,
	"created_at" datetime,
	"updated_at" datetime,
	"status" varchar(255),
	"sfw" bool,
	"source" varchar(255),
	"type" varchar(255),
	"released_at" datetime,
	"category" varchar(255),
	"external_id" varchar(255) UNIQUE,
	"external_url" varchar(255),
	"source_url" varchar(255),
	"source_category" varchar(255),
	"likes" integer,
	"summary" varchar(255)
)`

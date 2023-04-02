package media

import (
	"time"
)

func (r mediaRepo) AddTag(tag *Tag) error {
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	res, err := r.DB.NamedExec(
		"INSERT INTO tags (created_at, updated_at, name) VALUES (:created_at, :updated_at, :name)",
		tag,
	)
	if err != nil {
		return err
	}
	tagID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	tag.ID = uint(tagID)
	return nil
}

func (r mediaRepo) AddTagToPost(post *Post, externalTag string) error {
	tag, err := r.findTag(externalTag)
	if err != nil {
		tag = Tag{Name: externalTag}
		err = r.AddTag(&tag)

		if err != nil {
			return err
		}

		tag, err = r.findTag(externalTag)

		if err != nil {
			return err
		}
	}
	_, err = r.DB.Exec(
		"INSERT INTO posts_tags (post_id, tag_id) VALUES (?, ?)",
		post.ID, tag.ID,
	)
	return err
}

// Tag is tag assigned to added post
type Tag struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	Name string
}

func (r mediaRepo) findTag(tagName string) (Tag, error) {
	tag := Tag{}
	err := r.DB.Get(&tag, "SELECT id, name from tags where name = ?", tagName)
	return tag, err
}

// TagsSchema represents schema for "tags" table
var TagsSchema = `CREATE TABLE "tags" (
	"id" integer PRIMARY KEY AUTOINCREMENT,
	"created_at" datetime,
	"updated_at" datetime,
	"name" varchar(255) UNIQUE
)`

// PostsTagsSchema represents schema for "posts_tags" table
// that connects posts with tags (has and belongs to many)
var PostsTagsSchema = `CREATE TABLE "posts_tags" (
	"post_id" integer,
	"tag_id" integer,
	PRIMARY KEY ("post_id","tag_id")
)`

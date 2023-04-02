package media

import (
	"time"
)

func (r mediaRepo) AddLink(link *Link) error {
	link.CreatedAt = time.Now()
	link.UpdatedAt = time.Now()
	res, err := r.DB.NamedExec(
		`INSERT INTO links (
       created_at, updated_at, post_id, url
		 )
     VALUES (
       :created_at, :updated_at, :post_id, :url
     )`,
		link,
	)
	if err != nil {
		return err
	}
	linkID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	link.ID = uint(linkID)
	return nil
}

// Link represents bookmark link
type Link struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	PostID uint   `db:"post_id"`
	URL    string `db:"url"`
}

// LinksSchema represents schema for "links" table
var LinksSchema = `CREATE TABLE "links" (
	"id" integer PRIMARY KEY AUTOINCREMENT,
	"created_at" datetime,
	"updated_at" datetime,
	"post_id" integer,
	"url" varchar(255)
)`

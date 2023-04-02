package media

import (
	"time"
)

func (r mediaRepo) AddText(text *Text) error {
	text.CreatedAt = time.Now()
	text.UpdatedAt = time.Now()
	res, err := r.DB.NamedExec(
		`INSERT INTO texts (
	    created_at, updated_at, post_id, title, body
		)
	  VALUES (
	    :created_at, :updated_at, :post_id, :title, :body
	  )`,
		text,
	)
	if err != nil {
		return err
	}
	textID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	text.ID = uint(textID)
	return nil
}

// Text represents one text story
type Text struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	PostID uint `db:"post_id"`
	Title  string
	Body   string
}

// TextsSchema represents schema for "texts" table
var TextsSchema = `CREATE TABLE "texts" (
	"id" integer primary key autoincrement,
	"created_at" datetime,
	"updated_at" datetime,
	"post_id" integer,
	"title" varchar(255),
	"body" varchar(255)
)`

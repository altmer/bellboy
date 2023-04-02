package media

import (
	"time"
)

func (r mediaRepo) AddSubscription(sub *Subscription) error {
	sub.CreatedAt = time.Now()
	sub.UpdatedAt = time.Now()
	res, err := r.DB.NamedExec(
		`INSERT INTO subscriptions (
       created_at, updated_at, blog_name, url, source, description, title
		 )
     VALUES (
       :created_at, :updated_at, :blog_name, :url, :source, :description, :title
     )`,
		sub,
	)
	if err != nil {
		return err
	}
	subID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	sub.ID = uint(subID)
	return nil
}

func (r mediaRepo) ListSubscriptions() ([]Subscription, error) {
	var subscriptions []Subscription
	err := r.DB.Select(
		&subscriptions,
		"SELECT id, created_at, updated_at, url, blog_name, source, title, description FROM subscriptions",
	)
	return subscriptions, err
}

func (r mediaRepo) RemoveAllSubscriptions() error {
	_, err := r.DB.Exec("DELETE FROM subscriptions")
	return err
}

// Subscription represents particular blog that user is subscribed to
type Subscription struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	BlogName    string `db:"blog_name"`
	Source      string
	URL         string `db:"url"`
	Description string
	Title       string
}

// SubscriptionsSchema represents schema for "subscriptions" table
var SubscriptionsSchema = `CREATE TABLE "subscriptions" (
	"id" integer PRIMARY KEY AUTOINCREMENT,
	"created_at" datetime,
	"updated_at" datetime,
	"blog_name" varchar(255),
	"source" varchar(255),
	"url" varchar(255),
	"description" varchar(255),
	"title" varchar(255)
)`

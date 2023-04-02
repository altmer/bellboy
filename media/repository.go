package media

import (
	"github.com/jmoiron/sqlx"
)

var tables = []string{PostsSchema, PhotosSchema, VideosSchema, TextsSchema, LinksSchema, TagsSchema, PostsTagsSchema, SubscriptionsSchema}

type mediaRepo struct {
	DB *sqlx.DB
}

// Repository represents objects that handles media objetcs persistence
type Repository interface {
	AddPost(*Post) error
	AddLink(*Link) error
	AddText(*Text) error
	AddPhoto(*Photo) error
	AddVideo(*Video) error
	AddSubscription(*Subscription) error
	AddTag(*Tag) error
	AddTagToPost(*Post, string) error

	ListSubscriptions() ([]Subscription, error)
	GetPhotoPath(*Photo) string
	GetVideoPath(*Video) string
	GetVideoThumbnailPath(*Video) string
	PostExistsWithExternalID(string) bool

	RemoveAllSubscriptions() error
}

// NewRepository initializes media repository object
func NewRepository(DB *sqlx.DB) Repository {
	for _, table := range tables {
		DB.Exec(table)
	}
	return &mediaRepo{DB}
}

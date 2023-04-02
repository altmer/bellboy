package media

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
)

func (r mediaRepo) AddPhoto(photo *Photo) error {
	trx := mediaTransaction{
		insertCallback: func() error {
			photo.CreatedAt = time.Now()
			photo.UpdatedAt = time.Now()
			res, err := r.DB.NamedExec(
				`INSERT INTO photos (
					created_at, updated_at, post_id, caption, external_url, sfw
				)
				VALUES (
					:created_at, :updated_at, :post_id, :caption, :external_url, :sfw
				)`,
				photo,
			)
			if err != nil {
				return err
			}
			photoID, err := res.LastInsertId()
			if err != nil {
				return err
			}
			photo.ID = uint(photoID)
			return nil
		},
	}
	trx.validateUrls([]string{photo.ExternalURL})
	trx.save()
	trx.downloadAll([]downloadTask{
		downloadTask{url: photo.ExternalURL, localPath: r.GetPhotoPath(photo)},
	})
	return trx.err
}

func (r mediaRepo) GetPhotoPath(photo *Photo) string {
	return filepath.Join(viper.GetString("media_folder"), photo.FileName())
}

// FileName returns local file name where current photo (should be) stored
func (photo Photo) FileName() string {
	extension := extension(photo.ExternalURL)
	return fmt.Sprintf("photo_%d%s", photo.ID, extension)
}

// Photo represents one photo media object
type Photo struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	PostID      uint `db:"post_id"`
	Caption     string
	ExternalURL string `db:"external_url"`
	SFW         bool
}

// PhotosSchema represents schema for "photos" table
var PhotosSchema = `CREATE TABLE "photos" (
	"id" integer primary key autoincrement,
	"created_at" datetime,
	"updated_at" datetime,
	"post_id" integer,
	"caption" varchar(255),
	"external_url" varchar(255),
	"sfw" bool
)`

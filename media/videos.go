package media

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"time"
)

func (r mediaRepo) AddVideo(video *Video) error {
	trx := mediaTransaction{
		insertCallback: func() error {
			video.CreatedAt = time.Now()
			video.UpdatedAt = time.Now()
			res, err := r.DB.NamedExec(
				`INSERT INTO videos (
					created_at, updated_at, post_id, external_url, thumbnail_url
				)
				VALUES (
					:created_at, :updated_at, :post_id, :external_url, :thumbnail_url
				)`,
				video,
			)
			if err != nil {
				return err
			}
			videoID, err := res.LastInsertId()
			if err != nil {
				return err
			}
			video.ID = uint(videoID)
			return nil
		},
	}
	trx.validateUrls([]string{video.ExternalURL, video.ThumbnailURL})
	trx.save()
	trx.downloadAll([]downloadTask{
		downloadTask{url: video.ExternalURL, localPath: r.GetVideoPath(video)},
		downloadTask{url: video.ThumbnailURL, localPath: r.GetVideoThumbnailPath(video)},
	})
	return trx.err
}

func (r mediaRepo) GetVideoPath(video *Video) string {
	return filepath.Join(viper.GetString("media_folder"), video.FileName())
}

func (r mediaRepo) GetVideoThumbnailPath(video *Video) string {
	return filepath.Join(viper.GetString("media_folder"), video.ThumbnailFileName())
}

// FileName returns local file name where current video (should be) stored
func (video Video) FileName() string {
	extension := extension(video.ExternalURL)
	return fmt.Sprintf("video_%d%s", video.ID, extension)
}

// ThumbnailFileName returns local file name where current video thumbnail (should be) stored
func (video Video) ThumbnailFileName() string {
	extension := extension(video.ThumbnailURL)
	return fmt.Sprintf("video_%d_thumbnail%s", video.ID, extension)
}

// Video represents one video media object
type Video struct {
	ID        uint
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`

	PostID       uint   `db:"post_id"`
	ExternalURL  string `db:"external_url"`
	ThumbnailURL string `db:"thumbnail_url"`
}

// VideosSchema represents schema for "videos" table
var VideosSchema = `CREATE TABLE "videos" (
	"id" integer primary key autoincrement,
	"created_at" datetime,
	"updated_at" datetime,
	"post_id" integer,
	"external_url" varchar(255),
	"thumbnail_url" varchar(255)
)`

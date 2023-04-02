package media

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	httpmock "gopkg.in/jarcoal/httpmock.v1"
)

func TestVideoFileName(t *testing.T) {
	testCases := []struct {
		video    Video
		expected string
	}{
		{Video{ID: 1, ExternalURL: "http://example.com/some/path/filename.avi"}, "video_1.avi"},
		{Video{ID: 1234, ExternalURL: "http://example.com/filename.mpeg?q=t#fragment"}, "video_1234.mpeg"},
		{Video{ID: 54543, ExternalURL: "http://example.com"}, "video_54543"},
		{Video{ID: 2, ExternalURL: "not a url"}, "video_2"},
		{Video{ID: 10, ExternalURL: "http://vtt.tumblr.com/tumblr_nr3hcb0o041th7ykt_720.mp4"}, "video_10.mp4"},
	}

	for _, testCase := range testCases {
		actual := testCase.video.FileName()
		if actual != testCase.expected {
			t.Errorf("Expected vido with url [%#v] filename to be [%#v], got [%#v]", testCase.video.ExternalURL, testCase.expected, actual)
		}
	}
}
func TestVideoThumbnailFileName(t *testing.T) {
	testCases := []struct {
		video    Video
		expected string
	}{
		{Video{ID: 4, ThumbnailURL: "http://example.com/some/path/filename.jpg"}, "video_4_thumbnail.jpg"},
		{Video{ID: 65464, ThumbnailURL: "http://example.com/filename.jpeg?q=t#fragment"}, "video_65464_thumbnail.jpeg"},
		{Video{ID: 34543, ThumbnailURL: "http://example.com"}, "video_34543_thumbnail"},
		{Video{ID: 9, ThumbnailURL: "not a url"}, "video_9_thumbnail"},
		{Video{ID: 52, ThumbnailURL: "https://31.media.tumblr.com/tumblr_loivwok9aP2w1jfb9_frame1.jpg"}, "video_52_thumbnail.jpg"},
	}

	for _, testCase := range testCases {
		actual := testCase.video.ThumbnailFileName()
		if actual != testCase.expected {
			t.Errorf("Expected vido with url [%#v] filename to be [%#v], got [%#v]", testCase.video.ExternalURL, testCase.expected, actual)
		}
	}
}

func TestAddVideo(t *testing.T) {
	teardown := setup()
	defer teardown()

	testCases := []struct {
		video Video
		err   error
	}{
		{Video{ExternalURL: "http://example.com/path/filename.avi", ThumbnailURL: "http://example.com/path/filename.jpg", PostID: 1}, nil},
		{Video{ThumbnailURL: "http://example.com/path/filename.jpg"}, errors.New("parse : empty url")},
		{Video{ExternalURL: "path/filename.avi", ThumbnailURL: "http://example.com/path/filename.jpg"}, errors.New("parse path/filename.avi: invalid URI for request")},
	}

	for _, testCase := range testCases {
		httpmock.Activate()

		httpmock.RegisterResponder("GET", testCase.video.ExternalURL,
			httpmock.NewStringResponder(200, "video file contents"))

		httpmock.RegisterResponder("GET", testCase.video.ThumbnailURL,
			httpmock.NewStringResponder(200, "thumbnail file contents"))

		err := repo.AddVideo(&testCase.video)
		checkErrors(t, testCase.err, err)
		if err == nil {
			var dbVideo Video
			DB.Get(
				&dbVideo,
				"SELECT id, created_at, post_id, external_url, thumbnail_url FROM videos WHERE external_url = ? LIMIT 1",
				testCase.video.ExternalURL,
			)

			assert.NotEqual(t, uint(0), testCase.video.ID, "ID is not set for saved video")
			assert.NotEqual(t, uint(0), dbVideo.ID, "ID is not fetched from database")
			assert.Equal(t, testCase.video.PostID, dbVideo.PostID)
			assert.Equal(t, testCase.video.ThumbnailURL, dbVideo.ThumbnailURL)
			assert.Equal(t, testCase.video.ExternalURL, dbVideo.ExternalURL)
		}

		expectedCalls := []string{
			fmt.Sprintf("GET %s", testCase.video.ExternalURL),
			fmt.Sprintf("GET %s", testCase.video.ThumbnailURL),
		}
		expectedCallCount := 1
		if testCase.err != nil {
			expectedCallCount = 0
		}

		info := httpmock.GetCallCountInfo()
		for _, expectedCall := range expectedCalls {
			if info[expectedCall] != expectedCallCount {
				t.Errorf("Expected to receive %d http calls to [%s], got %d", expectedCallCount, expectedCall, info[expectedCall])
			}
		}

		// check file
		if testCase.err == nil {
			contents, _ := ioutil.ReadFile(repo.GetVideoPath(&testCase.video))
			if string(contents) != "video file contents" {
				t.Errorf("Wrong video file content: [%s]", contents)
			}
			contents, _ = ioutil.ReadFile(repo.GetVideoThumbnailPath(&testCase.video))
			if string(contents) != "thumbnail file contents" {
				t.Errorf("Wrong thumbnail file content: [%s]", contents)
			}
			// clean up
			os.Remove(repo.GetVideoPath(&testCase.video))
			os.Remove(repo.GetVideoThumbnailPath(&testCase.video))
		}

		// clean up
		httpmock.DeactivateAndReset()
	}
}

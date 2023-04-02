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

func TestPhotoFileName(t *testing.T) {
	testCases := []struct {
		photo    Photo
		expected string
	}{
		{Photo{ID: 1, ExternalURL: "http://example.com/some/path/filename.jpg"}, "photo_1.jpg"},
		{Photo{ID: 1234, ExternalURL: "http://example.com/filename.jpeg?q=t#fragment"}, "photo_1234.jpeg"},
		{Photo{ID: 54543, ExternalURL: "http://example.com"}, "photo_54543"},
		{Photo{ID: 2, ExternalURL: "not a url"}, "photo_2"},
		{Photo{ID: 10, ExternalURL: "https://77.media.tumblr.com/1595a97bc7e471e0c7bbbe371ce7429c/tumblr_mjeiuxCZTA1rktr9ro1_1280.jpg"}, "photo_10.jpg"},
	}

	for _, testCase := range testCases {
		actual := testCase.photo.FileName()
		if actual != testCase.expected {
			t.Errorf("Expected photo with url [%#v] filename to be [%#v], got [%#v]", testCase.photo.ExternalURL, testCase.expected, actual)
		}
	}
}

func TestAddPhoto(t *testing.T) {
	teardown := setup()
	defer teardown()

	testCases := []struct {
		photo Photo
		err   error
	}{
		{Photo{ExternalURL: "http://example.com/path/filename.jpg", Caption: "some photo", PostID: 1, SFW: true}, nil},
		{Photo{ExternalURL: "http://example.com/filename.png?q=43", Caption: "some other photo", PostID: 2, SFW: false}, nil},
		{Photo{Caption: "wrong photo"}, errors.New("parse : empty url")},
		{Photo{ExternalURL: "http//fdfds", Caption: "wrong photo url 2"}, errors.New("parse http//fdfds: invalid URI for request")},
	}

	for _, testCase := range testCases {
		httpmock.Activate()

		httpmock.RegisterResponder("GET", testCase.photo.ExternalURL,
			httpmock.NewStringResponder(200, "jpg file contents"))

		err := repo.AddPhoto(&testCase.photo)
		checkErrors(t, testCase.err, err)

		if err == nil {
			var dbPhoto Photo
			DB.Get(
				&dbPhoto,
				"SELECT id, created_at, post_id, caption, external_url, sfw FROM photos WHERE external_url = ? LIMIT 1",
				testCase.photo.ExternalURL,
			)
			assert.NotEqual(t, uint(0), testCase.photo.ID, "ID is not set for saved photo")
			assert.NotEqual(t, uint(0), dbPhoto.ID, "ID is not fetched from database")
			assert.Equal(t, testCase.photo.PostID, dbPhoto.PostID)
			assert.Equal(t, testCase.photo.Caption, dbPhoto.Caption)
			assert.Equal(t, testCase.photo.ExternalURL, dbPhoto.ExternalURL)
		}

		expectedCall := fmt.Sprintf("GET %s", testCase.photo.ExternalURL)
		expectedCallCount := 1
		if testCase.err != nil {
			expectedCallCount = 0
		}

		info := httpmock.GetCallCountInfo()
		if info[expectedCall] != expectedCallCount {
			t.Errorf("Expected to receive %d http calls to [%s], got %d", expectedCallCount, expectedCall, info[expectedCall])
		}

		// check file
		if testCase.err == nil {
			contents, _ := ioutil.ReadFile(repo.GetPhotoPath(&testCase.photo))
			if string(contents) != "jpg file contents" {
				t.Errorf("Wrong photo file content: [%s]", contents)
			}
			// clean up
			os.Remove(repo.GetPhotoPath(&testCase.photo))
		}

		// clean up
		httpmock.DeactivateAndReset()
	}

}

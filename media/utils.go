package media

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
)

func download(url, filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

func extension(url string) string {
	extensionRegexp := regexp.MustCompile("(?:(?:http://)|(?:https://)).+/.+\\.([a-zA-Z0-9]{2,6})(?:\\?.*)?$")
	matches := extensionRegexp.FindStringSubmatch(url)
	if len(matches) == 2 {
		return fmt.Sprintf(".%s", matches[1])
	}
	return ""
}

type mediaTransaction struct {
	err            error
	insertCallback func() error
}

func (trx *mediaTransaction) validateUrls(urls []string) error {
	if trx.err != nil {
		return trx.err
	}
	for _, ur := range urls {
		_, err := url.ParseRequestURI(ur)
		if err != nil {
			trx.err = err
			return err
		}
	}
	return nil
}

func (trx *mediaTransaction) downloadAll(objects []downloadTask) error {
	if trx.err != nil {
		return trx.err
	}
	for _, object := range objects {
		err := download(object.url, object.localPath)
		if err != nil {
			trx.err = err
			return err
		}
	}
	return nil
}

func (trx *mediaTransaction) save() error {
	if trx.err != nil {
		return trx.err
	}
	err := trx.insertCallback()
	if err != nil {
		trx.err = err
		return err
	}
	return nil
}

type downloadTask struct {
	url       string
	localPath string
}

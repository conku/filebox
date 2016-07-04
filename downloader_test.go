package downloader_test

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/qor/downloader"
	"github.com/qor/roles"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var Downloader *downloader.Downloader

func init() {
	root, _ := os.Getwd()
	os.Remove(root + "/test/downloads/a.csv.meta")
	Downloader = downloader.New(root + "/test/downloads")
}

func TestDownloader(t *testing.T) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	Downloader.MountTo(mux)
	req, err := http.Get(server.URL + "/downloads/a.csv")
	if err != nil || req.StatusCode != 200 {
		t.Errorf(color.RedString(fmt.Sprintf("Downloader error: can't get file")))
	}
	body, _ := ioutil.ReadAll(req.Body)
	if string(body) != "Column1,Column2\n" {
		t.Errorf(color.RedString(fmt.Sprintf("Downloader error: file'content is incorrect")))
	}

	permission := roles.Allow(roles.Read, "admin")
	if err := Downloader.Get("a.csv").SetPermission(permission); err != nil {
		t.Errorf(color.RedString(fmt.Sprintf("Downloader error: create meta file failure (%v)", err)))
	}

	req, err = http.Get(server.URL + "/downloads/a.csv")
	if err != nil || req.StatusCode != 404 {
		t.Errorf(color.RedString(fmt.Sprintf("Downloader error: should can't download file")))
	}
}

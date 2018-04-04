package radibrary_downloader_go

import (
	"net/http"
	"mime"
	"os"
	"log"
	"fmt"
	"io"
	"github.com/PuerkitoBio/goquery"
)

func GetDocFromUrl(url string) (*goquery.Document, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:58.0) Gecko/20100101 Firefox/58.0")

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	return goquery.NewDocumentFromResponse(res)
}


func DownloadBinaryFile(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:58.0) Gecko/20100101 Firefox/58.0")

	headRes, headErr := client.Do(req)
	if headErr != nil {
		return headErr
	}

	defer headRes.Body.Close()

	contentDisposition := headRes.Header.Get("Content-Disposition")

	_, params, err := mime.ParseMediaType(contentDisposition)
	filename := params["filename"]
	// Create the file
	out, err := os.Create(filename)
	if err != nil  {
		return err
	}
	defer out.Close()

	// Get the data
	downloadClient := &http.Client{}
	getReq, getErr := http.NewRequest("GET", url, nil)
	if getErr != nil {
		log.Fatal(fmt.Sprintf("Fail to create request object: %s", url))
		return getErr
	}
	getReq.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.13; rv:58.0) Gecko/20100101 Firefox/58.0")

	res, err := downloadClient.Do(getReq)
	if err != nil {
		log.Fatal(fmt.Sprintf("Fail to download: %s", url))
		return err
	}
	defer res.Body.Close()

	// Check server response
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", res.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil  {
		log.Fatal(fmt.Sprintf("Fail to copy: %s", url))
		return err
	}

	return nil
}

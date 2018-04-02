package main
import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"mime"
	"os"
	"io"
)

const hostURL = "http://radibrary.tistory.com/"

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

func GetSearchResults(url string) ([]string, error) {
	log.Printf("GetItems(url: %s)", url)

	doc, err := GetDocFromUrl(url)
	if err != nil {
		return nil, err
	}

	log.Print("Start to find CSS selector")
	// Find the review items
	items := doc.Find("div.list_content > .link_post").Map(func(i int, s *goquery.Selection) string {
		itemPath, err := s.Attr("href")
		if err != true {
			log.Fatal(err)
			return ""
		}

		title := s.Find(".tit_post").Text()
		log.Printf("Item %s -> %s", title, itemPath)

		return itemPath
	})

	return items, nil
}

func SearchPage(query string) []string {
	var allUrls []string

	for i := 1; ; i++ {
		encodedPageUrl := fmt.Sprintf("%s/search/%s?page=%d", hostURL, query, i)
		log.Println(encodedPageUrl)

		items, err := GetSearchResults(encodedPageUrl)
		if err != nil {
			break
		}

		if len(items) == 0 {
			break
		}

		for _, item := range items {
			allUrls = append(allUrls, item)
		}
	}

	return allUrls
}

func ExtractRadiofileFromPage(pageUrl string, result chan<- radiofile) error {
	doc, err := GetDocFromUrl(fmt.Sprintf("%s%s", hostURL, pageUrl))
	if err != nil {
		return err
	}

	title := doc.Find(".tit_post").Text()

	fileUrls := doc.Find(".moreless_content a").Map(func(i int, s *goquery.Selection) string {
		url, err := s.Attr("href")
		if err != true {
			return ""
		}

		return url
	})

	var radioFiles []radiofile
	for _, fileUrl := range fileUrls {
		radioFiles = append(radioFiles, radiofile{title: title, url: fileUrl})
	}

	for _, file := range radioFiles {
		result <- file
	}

	return nil
}


func radiofileDownloadWorker(id int, results <-chan radiofile) {
	for result := range results {
		log.Print(result)
		err := downloadBinaryFile(result.url)
		if err != nil {
			log.Print(fmt.Sprintf("Download Fail: %s", err))
		}
	}
}

func downloadBinaryFile(url string) error {
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

type radiofile struct {
	title string
	url  string
}

func main() {
	results := make(chan radiofile, 8)

	for w := 1; w <= 8; w++ {
		go radiofileDownloadWorker(w, results)
	}

	pageUrls := SearchPage("MELODY FLAG")
	log.Print("Pages", pageUrls)

	for i, pageUrl := range pageUrls {
		log.Print(fmt.Sprintf("Send %d: %s", i, pageUrl))
		ExtractRadiofileFromPage(pageUrl, results)
	}

	log.Print("Finish")
}

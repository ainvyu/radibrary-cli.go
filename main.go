package main
import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"./lib"
	"strings"
	"math"
	"flag"
)

const hostURL = "http://radibrary.tistory.com/"

func GetItemsFromSearchPage(url string) ([]string, error) {
	log.Printf("GetItems(url: %s)", url)

	doc, err := lib.GetDocFromUrl(url)
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

		items, err := GetItemsFromSearchPage(encodedPageUrl)
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
	doc, err := lib.GetDocFromUrl(fmt.Sprintf("%s%s", hostURL, pageUrl))
	if err != nil {
		return err
	}

	title := strings.TrimSpace(doc.Find(".area_title .tit_post").Text())

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

	log.Printf("Page %s: End to extract", pageUrl)

	return nil
}

func RadiofileDownloadWorker(id int, results <-chan radiofile, done chan<- bool) {
	for result := range results {
		log.Print(result)
		err := lib.DownloadBinaryFile(result.url)
		if err != nil {
			log.Printf("Download Fail %s - %s: %s", result.title, result.url, err)
		}
	}

	log.Printf("%d: Start return done", id)
	done <- true
	log.Printf("%d: End return done", id)
}

type radiofile struct {
	title string
	url  string
}

func main() {
	query := flag.String("query", "", "Query Sentence")
	log.Printf("Query: %s", *query)

	results := make(chan radiofile, int(math.Pow(2, 16)))
	done := make(chan bool)

	pageUrls := SearchPage(*query)
	log.Print("Pages", pageUrls)

	for i, pageUrl := range pageUrls {
		log.Printf("Send page URL %d: %s", i, pageUrl)
		//go ExtractRadiofileFromPage(pageUrl, results)
		ExtractRadiofileFromPage(pageUrl, results)
	}

	close(results)

	workerCount := 8
	log.Printf("Create download worker")
	for w := 0; w < workerCount; w++ {
		go RadiofileDownloadWorker(w, results, done)
	}

	for i := 0; i < workerCount; i++ {
		log.Printf("Wait i: %d", i)
		<- done
		log.Printf("Received %d", i)
	}

	log.Printf("Finish")
}

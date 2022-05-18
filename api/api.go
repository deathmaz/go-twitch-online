package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

type FetchedData struct {
	URL         string
	IsLive      bool
	Description string
}

func FetchDataForStream(link string, c chan FetchedData) {
	res, err := http.Get(link)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
        c <- FetchedData{
            URL: link,
            Description: fmt.Sprintf("There was an error while fetching stream: %d %s", res.StatusCode, res.Status),
            IsLive: false,
        };
        return
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	script := doc.Find("script[type='application/ld+json']")
	var isLive bool
	description := "No description"
	if script.Text() != "" {
		type node map[string]interface{}
		var resJson []node
		jsonError := json.Unmarshal([]byte(script.Text()), &resJson)
		if jsonError != nil {
			fmt.Println("Error while parsing json", jsonError)
		}
		if resJson[0]["publication"] != nil {
			publication := resJson[0]["publication"].(map[string]interface{})
			if publication == nil {
				isLive = false
			} else {
				isLive = true
			}
		}
		description = "No description"
		if resJson[0]["description"] != nil {
			description = resJson[0]["description"].(string)
		}
	}

	payload := FetchedData{URL: link, Description: description, IsLive: isLive}
	c <- payload
}

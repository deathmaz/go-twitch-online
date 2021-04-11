package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fatih/color"
)

type FetchedData struct {
	url         string
	isLive      bool
	description string
}

type Stream struct {
	id          string
	url         string
	isLive      bool
	description string
	index       int
}

func (s Stream) show() {
	statusStyles := color.New(color.FgRed, color.Bold)
	statusText := statusStyles.Sprintf("%s", "offline")
	if s.isLive {
		statusStyles = color.New(color.FgGreen, color.Bold)
		statusText = statusStyles.Sprintf("%s", "live")
	}
	linkText := color.New(color.FgBlue, color.Underline)
	indexText := color.New(color.FgGreen, color.Bold)
	fmt.Println(fmt.Sprintf("%s %s is %s", indexText.Sprintf("%d.", s.index), linkText.Sprintf("%s", s.url), statusText))

	descriptionStyles := color.New(color.FgMagenta)
	fmt.Println("Streaming:", descriptionStyles.Sprintf("%s", s.description))
	fmt.Println("")
}

type StreamList struct {
	inner   []Stream
	channel chan FetchedData
}

func (sl *StreamList) createFromIds(ids []string) {
	for i, id := range ids {
		sl.inner = append(sl.inner, Stream{
			id:    id,
			url:   "https://m.twitch.tv/" + id,
			index: i,
		})
	}
}

func (sl *StreamList) fetchAll() {
	fmt.Println("")
	fmt.Println("Fetching all steams")
	for _, stream := range sl.inner {
		go fetchDataForStream(stream.url, sl.channel)
	}

	for i := 0; i < len(sl.inner); i++ {
		data := <-sl.channel
		for i := range sl.inner {
			t := &sl.inner[i]
			if t.url == data.url {
				t.description = data.description
				t.isLive = data.isLive
			}
		}
	}
}

func (sl StreamList) show() {
	fmt.Println("")
	fmt.Println("Displaying all data")
	sort.Slice(sl.inner, func(i, j int) bool { return sl.inner[i].isLive })
	for _, stream := range sl.inner {
		stream.show()
	}
}

func (sl StreamList) fetchAndShow() {
	sl.fetchAll()
	sl.show()
}

func (sl StreamList) showOnlyLive() {
	fmt.Println("")
	fmt.Println("Displaying all data")
	sort.Slice(sl.inner, func(i, j int) bool { return sl.inner[i].isLive })
	for _, stream := range sl.inner {
		if stream.isLive {
			stream.show()
		}
	}
}

func main() {
	userList := readConfig()

	c := make(chan FetchedData)
	streamList := StreamList{
		channel: c,
	}
	streamList.createFromIds(userList)
	streamList.fetchAll()
	streamList.show()

	mainMenu(streamList)
}

func mainMenu(sl StreamList) {
	for {
		fmt.Println("")
		fmt.Println("What would you like to do?")
		fmt.Println("1. Play stream")
		fmt.Println("2. Refetch data")
		fmt.Println("")
		fmt.Println("Enter selection:")

		input := getInput()
		switch input {
		case "1":
			playStream(sl)
		case "2":
			sl.fetchAndShow()
		}
	}
}

func playStream(sl StreamList) {
	sl.showOnlyLive()
	fmt.Println("Enter video index:")
	input := getInput()
	for _, stream := range sl.inner {
		if input == fmt.Sprintf("%d", stream.index) {
			if playVideo(stream.url) != nil {
				fmt.Println("cant't play stream")
			}
		}
	}
}

func playVideo(url string) error {
	cmd := exec.Command("bash", "-c", "streamlink --player=mpv "+url)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func fetchDataForStream(link string, c chan FetchedData) {
	res, err := http.Get(link)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	script := doc.Find("script[type='application/ld+json']")
	type node map[string]interface{}
	var resJson []node
	json.Unmarshal([]byte(script.Text()), &resJson)
	var isLive bool
	if resJson[0]["publication"] != nil {
		publication := resJson[0]["publication"].(map[string]interface{})
		if publication == nil {
			isLive = false
		} else {
			isLive = true
		}
	}
	description := "No description"
	if resJson[0]["description"] != nil {
		description = resJson[0]["description"].(string)
	}

	payload := FetchedData{url: link, description: description, isLive: isLive}
	c <- payload
}

func getInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("An error occured while reading input. Please try again", err)
	}

	input = strings.TrimSuffix(input, "\n")
	return input
}

func readConfig() []string {
	usr, _ := user.Current()
	file, err := os.Open(usr.HomeDir + "/.config/go-twitch-online/users")
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var userList []string

	for scanner.Scan() {
		userList = append(userList, scanner.Text())
	}

	file.Close()

	return userList
}

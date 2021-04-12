package stream

import (
	"fmt"
	"sort"

	"github.com/deathmaz/go-twitch-online/api"
	"github.com/fatih/color"
)

type Stream struct {
	id          string
	URL         string
	isLive      bool
	description string
	Index       int
}

type StreamList struct {
	Inner   []Stream
	Channel chan api.FetchedData
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
	fmt.Println(fmt.Sprintf("%s %s is %s", indexText.Sprintf("%d.", s.Index), linkText.Sprintf("%s", s.URL), statusText))

	descriptionStyles := color.New(color.FgMagenta)
	fmt.Println("Streaming:", descriptionStyles.Sprintf("%s", s.description))
	fmt.Println("")
}

func (sl *StreamList) CreateFromIds(ids []string) {
	for i, id := range ids {
		sl.Inner = append(sl.Inner, Stream{
			id:    id,
			URL:   "https://m.twitch.tv/" + id,
			Index: i,
		})
	}
}

func (sl *StreamList) FetchAll() {
	fmt.Println("")
	fmt.Println("Fetching all steams")
	for _, stream := range sl.Inner {
		go api.FetchDataForStream(stream.URL, sl.Channel)
	}

	for i := 0; i < len(sl.Inner); i++ {
		data := <-sl.Channel
		for i := range sl.Inner {
			t := &sl.Inner[i]
			if t.URL == data.URL {
				t.description = data.Description
				t.isLive = data.IsLive
			}
		}
	}
}

func (sl StreamList) Show() {
	fmt.Println("")
	fmt.Println("Displaying all data")
	sort.Slice(sl.Inner, func(i, j int) bool { return sl.Inner[i].isLive })
	for _, stream := range sl.Inner {
		stream.show()
	}
}

func (sl StreamList) FetchAndShow() {
	sl.FetchAll()
	sl.Show()
}

func (sl StreamList) ShowOnlyLive() {
	fmt.Println("")
	fmt.Println("Displaying all data")
	sort.Slice(sl.Inner, func(i, j int) bool { return sl.Inner[i].isLive })
	for _, stream := range sl.Inner {
		if stream.isLive {
			stream.show()
		}
	}
}
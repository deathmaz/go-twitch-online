package menu

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/deathmaz/go-twitch-online/stream"
)

func MainMenu(sl stream.List) {
	for {
		fmt.Println("")
		fmt.Println("What would you like to do?")
		fmt.Println("1. Play stream")
		fmt.Println("2. Refetch data")
		fmt.Println("3. Show only live streams")
		fmt.Println("")
		fmt.Println("Enter selection:")

		input := getInput()
		switch input {
		case "1":
			playStream(sl)
		case "2":
			sl.FetchAndShow()
		case "3":
			sl.ShowOnlyLive()
		}
	}
}

func playStream(sl stream.List) {
	sl.ShowOnlyLive()
	fmt.Println("Enter video index:")
	input := getInput()
	for _, stream := range sl.Inner {
		if input == fmt.Sprintf("%d", stream.Index) {
			if playVideo(stream.URL) != nil {
				fmt.Println("cant't play stream")
			}
		}
	}
}

func playVideo(url string) error {
	cmd := exec.Command("bash", "-c", "streamlink --player=mpv "+url)
	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	return nil
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

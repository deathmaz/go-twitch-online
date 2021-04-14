package main

import (
	"github.com/deathmaz/go-twitch-online/api"
	"github.com/deathmaz/go-twitch-online/config"
	"github.com/deathmaz/go-twitch-online/menu"
	"github.com/deathmaz/go-twitch-online/stream"
)

func main() {
	userList := config.ReadConfig()

	c := make(chan api.FetchedData)
	streamList := stream.List{
		Channel: c,
	}
	streamList.CreateFromIds(userList)
	streamList.FetchAll()
	streamList.Show()

	menu.MainMenu(streamList)
}

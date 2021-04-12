package config

import (
	"bufio"
	"log"
	"os"
	"os/user"
)

func ReadConfig() []string {
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

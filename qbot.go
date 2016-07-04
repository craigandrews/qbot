package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/slack"
	"github.com/doozr/qbot/usercache"
)

type Notification struct {
	Channel string
	Message string
}

func listen(name string, connection *slack.Slack,
	messageChan chan slack.RtmMessage, userUpdateChan chan slack.UserInfo) {

	for {
		// read each incoming message
		m, err := connection.GetEvent()
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if m.Type == "message" {
			msg := slack.ConvertEventToMessage(m)
			if strings.HasPrefix(msg.Text, name) || strings.HasPrefix(msg.Text, "<@"+ connection.Id+">") {
				messageChan <- msg
			}
		}

		// see if it's a user update
		if m.Type == "user_change" {
			uc := slack.ConvertEventToUserChange(m)
			userUpdateChan <- uc.User
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: qbot <token> <data file>")
		os.Exit(1)
	}

	// Get command line parameters
	token := os.Args[1]
	filename := os.Args[2]

	// Instantiate state
	connection := connectToSlack(token)
	userCache := getUserList(connection)
	name := getBotName(userCache, connection.Id)
	q := loadQueue(filename)

	// Set up command and response processors
	notifications := notification.New(userCache)
	commands := command.New(notifications, userCache)

	// Create channels
	messageChan := make(chan slack.RtmMessage, 100)
	saveChan := make(chan queue.Queue, 5)
	notifyChan := make(chan Notification, 5)
	userUpdateChan := make(chan slack.UserInfo, 5)

	// Start goroutines
	go MessageDispatch(name, q, commands, messageChan, saveChan, notifyChan)
	go Save(filename, saveChan)
	go Notify(connection, notifyChan)
	go UpdateUser(userCache, userUpdateChan)

	// Dispatch incoming events
	log.Println("Ready to receive events")
	listen(name, connection, messageChan, userUpdateChan)
}

func connectToSlack(token string) (connection *slack.Slack) {
	log.Print("Connecting to Slack")
	connection, err := slack.New(token)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func getUserList(connection *slack.Slack) (userCache *usercache.UserCache) {
	log.Println("Getting user list")
	users, err := connection.GetUserList()
	if err != nil {
		log.Fatal(err)
	}
	userCache = usercache.New(users)
	return
}

func getBotName(userCache *usercache.UserCache, id string) (name string) {
	name = userCache.GetUserName(id)
	if name == "" {
		log.Fatal("Could not get username of bot")
	}
	log.Printf("Responding to requests directed to @%s", name)
	return
}

func loadQueue(filename string) (q queue.Queue) {
	log.Printf("Attempting to load queue from %s", filename)
	q, err := queue.Load(filename)
	if err != nil {
		log.Fatal("Error loading queue: %s", err)
	}
	return
}


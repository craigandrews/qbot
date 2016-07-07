package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/dispatch"
	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/slack"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/util"
)

func listen(name string, connection *slack.Slack, messageChan dispatch.MessageChan, userChan dispatch.UserChan) {

	for {
		// read each incoming message
		e, err := connection.GetEvent()
		if err != nil {
			log.Fatal(err)
		}

		// see if we're mentioned
		if e.Type == "message" {
			m := slack.ConvertEventToMessage(e)
			if strings.HasPrefix(m.Text, name) || strings.HasPrefix(m.Text, "<@"+connection.Id+">") {
				_, m.Text = util.StringPop(m.Text)
				messageChan <- m
			}
		}

		// see if it's a user update
		if e.Type == "user_change" {
			uc := slack.ConvertEventToUserChange(e)
			userChan <- uc.User
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
	messageChan := make(dispatch.MessageChan, 100)
	saveChan := make(dispatch.SaveChan, 5)
	notifyChan := make(dispatch.NotifyChan, 5)
	userChan := make(dispatch.UserChan, 5)

	// Start goroutines
	go dispatch.Message(name, q, commands, messageChan, saveChan, notifyChan)
	go dispatch.Save(filename, saveChan)
	go dispatch.Notify(connection, notifyChan)
	go dispatch.User(userCache, userChan)

	// Dispatch incoming events
	log.Println("Ready to receive events")
	listen(name, connection, messageChan, userChan)
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

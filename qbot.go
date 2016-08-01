package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/doozr/goslack"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/dispatch"
	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
	"github.com/doozr/qbot/util"
)

// Version is the current release version
const Version = "1.2"

func listen(name string, connection *goslack.Connection, messageChan dispatch.MessageChan, userChan dispatch.UserChan) {

	for {
		// read each incoming message
		e := <-connection.RealTime

		// see if we're mentioned
		if e.Type == "message" {
			m, err := e.RtmMessage()
			if err != nil {
				log.Println(err)
				continue
			}

			directedAtUs := strings.HasPrefix(m.Text, name) || strings.HasPrefix(m.Text, "<@"+connection.ID+">")
			if directedAtUs {
				_, m.Text = util.StringPop(m.Text)
				messageChan <- m
			} else if util.IsPrivateChannel(m.Channel) {
				messageChan <- m
			}
		}

		// see if it's a user update
		if e.Type == "user_change" {
			uc, err := e.RtmUserChange()
			if err != nil {
				log.Println(err)
				continue
			}
			userChan <- uc.User
		}
	}
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: qbot <token> <data file>")
		os.Exit(1)
	}

	log.Printf("Qbot version %s", Version)

	// Get command line parameters
	token := os.Args[1]
	filename := os.Args[2]

	// Instantiate state
	connection := connectToSlack(token)
	userCache := getUserList(connection)
	name := getBotName(userCache, connection.ID)
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

func connectToSlack(token string) (connection *goslack.Connection) {
	log.Print("Connecting to Slack")
	connection, err := goslack.New(token)
	if err != nil {
		log.Fatal(err)
	}
	return
}

func getUserList(connection *goslack.Connection) (userCache *usercache.UserCache) {
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
		log.Fatalf("Error loading queue: %s", err)
	}
	return
}

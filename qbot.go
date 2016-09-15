package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/dispatch"
	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

// Version is the current release version
var Version string

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: qbot <token> <data file>")
		os.Exit(1)
	}

	if Version != "" {
		log.Printf("Qbot version %s", Version)
	} else {
		log.Printf("Qbot <unversioned build>")
	}

	// Get command line parameters
	token := os.Args[1]
	filename := os.Args[2]

	// Turn on jot if required
	if os.Getenv("QBOT_DEBUG") == "true" {
		jot.Enable()
	}

	// Synchronisation primitives
	waitGroup := sync.WaitGroup{}
	done := make(chan struct{})
	defer func() {
		jot.Print("Closing done channel")
		close(done)

		jot.Print("Awaiting all goroutines")
		waitGroup.Wait()

		jot.Print("Shutdown complete")
	}()

	// Connect to Slack
	client, err := guac.New(token).PersistentRealTime()
	if err != nil {
		log.Fatal("Error connecting to Slack ", err)
	}

	// Instantiate state
	userCache := getUserList(client.WebClient)
	name := client.Name()
	q := loadQueue(filename)

	// Set up command and response processors
	notifications := notification.New(userCache)
	commands := command.New(notifications, userCache)

	// Create channels
	messageChan := make(dispatch.MessageChan, 100)
	saveChan := make(dispatch.SaveChan, 5)
	notifyChan := make(dispatch.NotifyChan, 5)
	userChan := make(dispatch.UserChan, 5)
	defer func() {
		jot.Println("Closing channels")
		close(messageChan)
		close(saveChan)
		close(notifyChan)
		close(userChan)
	}()

	// Start goroutines
	waitGroup.Add(4)
	go dispatch.Message(name, q, commands, messageChan, saveChan, notifyChan, &waitGroup)
	go dispatch.Save(filename, saveChan, &waitGroup)
	go dispatch.Notify(client, notifyChan, &waitGroup)
	go dispatch.User(userCache, userChan, &waitGroup)

	// Dispatch incoming events
	jot.Println("Ready to receive events")
	waitGroup.Add(1)
	go listen(name, client, messageChan, userChan, done, &waitGroup)

	// Wait for signals to stop
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGKILL)

	// Wait for a signal
	s := <-sig
	log.Printf("Received %s signal - shutting down", s)
	client.Close()
}

func getUserList(client guac.WebClient) (userCache *usercache.UserCache) {
	log.Println("Getting user list")
	users, err := client.UsersList()
	if err != nil {
		log.Fatal(err)
	}
	userCache = usercache.New(users)
	return
}

func loadQueue(filename string) (q queue.Queue) {
	q, err := queue.Load(filename)
	if err != nil {
		log.Fatalf("Error loading queue: %s", err)
	}
	log.Printf("Loaded queue from %s", filename)
	return
}

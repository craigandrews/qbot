package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/doozr/guac"
	"github.com/doozr/jot"
	"github.com/doozr/qbot/command"
	"github.com/doozr/qbot/notification"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

// Version is the current release version
var Version string

// DoneChan is a channel used for informing go routines to shut down
type DoneChan chan struct{}

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
	done := make(DoneChan)

	// Connect to Slack
	client, err := guac.New(token).RealTime()
	if err != nil {
		log.Fatal("Error connecting to Slack ", err)
	}
	log.Print("Connected to slack as ", client.Name())

	// Instantiate state
	userCache := getUserList(client)
	name := client.Name()
	jot.Print("qbot: name is ", name)
	q := loadQueue(filename)

	// Set up command and response processors
	notifications := notification.New(userCache)
	commands := command.New(notifications, userCache)

	// Create dispatchers
	notify := createNotifier(client)
	persist := createPersister(filename)
	messageHandler := createMessageHandler(client.ID(), client.Name(), q, commands, notify, persist)
	userChangeHandler := createUserChangeHandler(userCache)

	// keepalive
	waitGroup.Add(1)
	go keepalive(client, done, &waitGroup)

	// Receive incoming events
	receiver := createReceiver(client)
	events := receive(receiver, done, &waitGroup)

	// Dispatch incoming events
	jot.Println("qbot: ready to receive events")
	dispatcher := createDispatcher(client, 1*time.Minute, messageHandler, userChangeHandler)
	abort := dispatch(dispatcher, events, done, &waitGroup)

	// Wait for signals to stop
	sig := addSignalHandler()

	// Wait for a signal
	select {
	case err := <-abort:
		if err != nil {
			log.Print("Error: ", err)
		}
		log.Print("Execution terminated - shutting down")
	case s := <-sig:
		log.Printf("Received %s signal - shutting down", s)
	}

	close(done)
	client.Close()
	waitGroup.Wait()

	jot.Print("qbot: shutdown complete")
}

func receive(receiver Receiver, done DoneChan, waitGroup *sync.WaitGroup) (events guac.EventChan) {
	events = make(guac.EventChan)

	waitGroup.Add(1)
	jot.Print("receive starting up")
	go func() {
		err := receiver(events, done)
		if err != nil {
			log.Print("Error receiving events: ", err)
		}

		close(events)
		jot.Print("receive done")
		waitGroup.Done()
	}()
	return
}

func dispatch(dispatcher Dispatcher, events guac.EventChan, done DoneChan, waitGroup *sync.WaitGroup) (abort chan error) {
	abort = make(chan error)

	waitGroup.Add(1)
	jot.Print("dispatch starting up")
	go func() {
		err := dispatcher(events, done)
		if err != nil {
			abort <- err
		}

		close(abort)
		jot.Print("dispatch done")
		waitGroup.Done()
	}()
	return
}

func addSignalHandler() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGKILL)
	return sig
}

func getUserList(client guac.WebClient) (userCache *usercache.UserCache) {
	log.Println("Getting user list")
	users, err := client.UsersList()
	if err != nil {
		log.Fatal(err)
	}
	userCache = usercache.New(users)
	jot.Print("loaded user list: ", userCache)
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

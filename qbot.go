package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
var Version = "<unversioned build>"

// DoneChan is a channel used for informing go routines to shut down
type DoneChan chan struct{}

func main() {
	log.Printf("Qbot version %s", Version)

	// Turn on jot if required
	if os.Getenv("QBOT_DEBUG") == "true" {
		jot.Enable()
	}

	token, filename := parseArgs()

	// Synchronisation primitives
	waitGroup := sync.WaitGroup{}
	done := make(DoneChan)

	// Connect to Slack
	client := connectToSlack(token)
	log.Print("Connected to slack as ", client.Name())

	// Instantiate state
	userCache := getUserList(client)
	q := loadQueue(filename)

	// Set up command and response processors
	notifications := notification.New(userCache)
	commands := command.New(notifications, userCache)

	// Create dispatchers
	notify := createNotifier(client)
	persist := createPersister(filename)
	messageHandler := createMessageHandler(client.ID(), client.Name(), q, commands, notify, persist)
	userChangeHandler := createUserChangeHandler(userCache)

	// start keepalive
	startKeepAlive(client, done, &waitGroup)

	// Receive incoming events
	receiver := createReceiver(client)
	events := receive(receiver, done, &waitGroup)

	// Dispatch incoming events
	jot.Println("qbot: ready to receive events")
	dispatcher := createDispatcher(client, 1*time.Minute, messageHandler, userChangeHandler)
	abort := dispatch(dispatcher, events, done, &waitGroup)

	// Wait for signals to stop
	sig := addSignalHandler()

	// Wait for a signal or an error to kill the process
	wait(sig, abort)

	// Shut it down
	close(done)
	client.Close()
	waitGroup.Wait()

	jot.Print("qbot: shutdown complete")
}

func parseArgs() (token, filename string) {
	if len(os.Args) < 3 {
		fmt.Println("Usage: qbot <token> <data file>")
		os.Exit(1)
	}
	// Get command line parameters
	token = os.Args[1]
	filename = os.Args[2]
	return
}

func connectToSlack(token string) guac.RealTimeClient {
	client, err := guac.New(token).RealTime()
	if err != nil {
		log.Fatal(err)
	}
	return client
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
	q = queue.Queue{}
	if _, err := os.Stat(filename); err == nil {
		dat, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Fatalf("Error loading queue: %s", err)
		}
		json.Unmarshal(dat, &q)
		jot.Printf("loadQueue: read queue from %s: %v", filename, q)
		log.Printf("Loaded queue from %s", filename)
	}
	return q
}

func addSignalHandler() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGKILL)
	return sig
}

func wait(sig chan os.Signal, abort chan error) {
	select {
	case err := <-abort:
		if err != nil {
			log.Print("Error: ", err)
		}
		log.Print("Execution terminated - shutting down")
	case s := <-sig:
		log.Printf("Received %s signal - shutting down", s)
	}
}

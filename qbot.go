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
	"github.com/doozr/qbot/dispatch"
	"github.com/doozr/qbot/queue"
	"github.com/doozr/qbot/usercache"
)

// Version is the current release version
var Version string

// World is all the state that could be global, but isn't
//
// Useful for making sure we can tear down and rebuild on connection failure
type World struct {
	Client    guac.WebClient
	Q         queue.Queue
	UserCache *usercache.UserCache
	SaveChan  dispatch.SaveChan
	UserChan  dispatch.UserChan
	Done      chan struct{}
	WaitGroup *sync.WaitGroup
}

func run(world World) {
	for {
		select {
		case <-world.Done:
			jot.Print("WaitGroup.Done: run")
			world.WaitGroup.Done()
			return
		default:
			r, ok := mustConnect(world.Client, world.Done)
			if ok {
				log.Printf("Responding as %s", r.Name())
				listen(r, world)
			}
		}
	}
}

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

	if os.Getenv("JOTTER_ENABLE") == "true" {
		jot.Enable()
	}

	// Get command line parameters
	token := os.Args[1]
	filename := os.Args[2]

	client := guac.New(token)
	var wg sync.WaitGroup

	// Load user list
	userCache, err := getUserList(client)
	if err != nil {
		log.Print("Error loading user list", err)
		os.Exit(1)
	}

	// Load queue
	queue, err := loadQueue(filename)
	if err != nil {
		log.Print("Error loading queue", err)
		os.Exit(1)
	}

	// Instantiate state
	world := World{
		Client:    client,
		Q:         queue,
		UserCache: userCache,
		SaveChan:  make(dispatch.SaveChan, 5),
		UserChan:  make(dispatch.UserChan, 5),
		Done:      make(chan struct{}),
		WaitGroup: &wg,
	}

	// Listen for notifications to save and update global state
	wg.Add(2)
	go dispatch.Save(filename, world.SaveChan, &wg)
	go dispatch.User(world.UserCache, world.UserChan, &wg)

	// Run the receiver
	wg.Add(1)
	go run(world)

	// Wait for signals to stop
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	signal.Notify(sig, syscall.SIGTERM)
	signal.Notify(sig, syscall.SIGKILL)

	// Wait for a signal
	s := <-sig
	log.Printf("Received %s signal - shutting down", s)

	// Shut down all the qeueus
	close(world.Done)
	close(world.SaveChan)
	close(world.UserChan)

	// Wait for all the goroutines to stop
	wg.Wait()
	log.Println("Shutdown complete")
}

func getUserList(client guac.WebClient) (userCache *usercache.UserCache, err error) {
	log.Println("Getting user list")
	users, err := client.UsersList()
	if err != nil {
		return
	}
	userCache = usercache.New(users)
	return
}

func loadQueue(filename string) (q queue.Queue, err error) {
	log.Printf("Attempting to load queue from %s", filename)
	return queue.Load(filename)
}

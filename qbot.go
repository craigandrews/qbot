package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/doozr/guac"
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

	// Get command line parameters
	token := os.Args[1]
	filename := os.Args[2]

	client := guac.New(token)
	var wg sync.WaitGroup

	// Instantiate state
	world := World{
		Client:    client,
		Q:         loadQueue(filename),
		UserCache: getUserList(client),
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

	// Run forever or until the done channel is closed
	select {
	case s := <-sig:
		log.Printf("Received %s signal - shutting down", s)
		close(world.Done)
		close(world.SaveChan)
		close(world.UserChan)
	}

	// Wait for everything to shut down
	wg.Wait()
	log.Println("Shutdown complete")
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
	log.Printf("Attempting to load queue from %s", filename)
	q, err := queue.Load(filename)
	if err != nil {
		log.Fatalf("Error loading queue: %s", err)
	}
	return
}

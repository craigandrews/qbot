package qbot

// Version is the current release version
var Version = "<unversioned build>"

// DoneChan is a channel used for informing go routines to shut down
type DoneChan chan struct{}

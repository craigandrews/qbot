package qbot

// version is the current release version.
var version = "<unversioned build>"

// DoneChan is a channel used for informing go routines to shut down.
type DoneChan chan struct{}

// Version of the running Qbot or `<unversioned build>` if built locally.
func Version() string {
	return version
}

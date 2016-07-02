package notification

import (
	"github.com/doozr/qbot/queue"
	"fmt"
)

func Join(i queue.Item) string {
	return fmt.Sprintf("@%s has joined the queue (%s)", i.Name, i.Reason)
}

func Active(i queue.Item) string {
	return fmt.Sprintf("@%s now has the token (%s)", i.Name, i.Reason)
}
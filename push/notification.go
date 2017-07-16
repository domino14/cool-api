// Package push implements our push notifications. For now, this will
// just be a mock.
package push

import (
	"fmt"
	"strings"
)

func notify(userid string, notification string) {
	delimiter := strings.Repeat("-", 30) + "\n"
	templateStr := delimiter + fmt.Sprintf(
		"[Push Notification for user %v]\n", userid) +
		"     " + notification + "\n" +
		delimiter + "\n"

	fmt.Println(templateStr)
}

func Notify(userid string, notification string) {
	// This goroutine, and the fact that the API in general uses goroutines
	// for HTTP requests, allows the API to scale more easily. We don't block
	// until all push notifications are delivered, instead we hand them off
	// in a goroutine and exit. This could be a separate microservice
	// or job queue later on.
	go func() {
		notify(userid, notification)
	}()
}

func NotifyMultiple(userids []string, notification string) {
	for _, userid := range userids {
		Notify(userid, notification)
	}
}

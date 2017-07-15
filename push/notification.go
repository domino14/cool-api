// Package push implements our push notifications. For now, this will
// just be a mock.
package push

import (
	"fmt"
	"strings"
)

func Notify(userid string, notification string) {
	delimiter := strings.Repeat("-", 30) + "\n"
	templateStr := delimiter + fmt.Sprintf(
		"[Push Notification for user %v]\n", userid) +
		"     " + notification + "\n" +
		delimiter + "\n"

	fmt.Println(templateStr)
}

func NotifyMultiple(userids []string, notification string) {
	// This goroutine, and the fact that the API in general uses goroutines
	// for HTTP requests, allows the API to scale more easily. We don't block
	// until all push notifications are delivered, instead we hand them off
	// in a goroutine and exit. This could be a separate microservice
	// or job queue later on.
	for _, userid := range userids {
		go func(uid string, notification string) {
			Notify(uid, notification)
		}(userid, notification)
	}
}

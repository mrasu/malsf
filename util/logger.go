package util

import (
	"log"
)

var debug bool = false

func SetDebug(b bool) {
	debug = b
}

func LogAction(name string, id int32, message string) {
	printf("%s (Action) Id: %d, Message: %s", name, id, message)
}

func LogActionReceived(name string, fromName string, fromId int32, message string) {
	printf("%s (ReceiveAction) From: %s(%d), Message: %s", name, fromName, fromId, message)
}

func LogReaction(name string, id int32, fromName string, fromId int32, message string) {
	printf("%s (Reaction) Id: %d, From: %s(%d), Message: %s", name, id, fromName, fromId, message)
}

func LogSwimMethod(isServer bool, phase string, message string) {
	if isServer {
		printf("SWIM: Send: %s (%s)", phase, message)
	} else {
		printf("SWIM: Catch: %s (%s)", phase, message)
	}
}

func Logf(format string, v ...interface{}) {
	printf(format, v...)
}

func printf(format string, v ...interface{}) {
	if debug {
		log.Printf(format, v...)
	}
}

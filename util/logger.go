package util

import (
	"log"
)

func LogAction(name string, id int32, message string) {
	log.Printf("%s (Action) Id: %d, Message: %s", name, id, message)
}

func LogActionReceived(name string, fromName string, fromId int32, message string) {
	log.Printf("%s (ReceiveAction) From: %s(%d), Message: %s", name, fromName, fromId, message)
}

func LogReaction(name string, id int32, fromName string, fromId int32, message string) {
	log.Printf("%s (Reaction) Id: %d, From: %s(%d), Message: %s", name, id, fromName, fromId, message)
}

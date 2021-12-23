package main

import (
	"flag"
	"log"
	"strings"

	"github.com/mvanwaaijen/wirepusher"
)

const DefaultType string = "Default"

// example database of users
var users map[string]*wirepusher.User

func init() {
	// populate example database
	users = make(map[string]*wirepusher.User)
	users["user001"] = wirepusher.NewUser("user001", "x8Xx7Xx6X", "myS3cret")     // send encrypted
	users["user002"] = wirepusher.NewUser("user002", "x5Yx9Yx0X", "myS3cretP@ss") // send encrypted
	users["user003"] = wirepusher.NewUser("user003", "z3Xx2Zx1X", "")             // send plaintext
}

func main() {
	var (
		user      string
		title     string
		body      string
		msgid     int
		actionURL string
		imgURL    string
		action    string
		msgtype   string
	)
	flag.StringVar(&user, "user", "", "name of user to send message to")
	flag.StringVar(&msgtype, "type", DefaultType, "defined type in the app on the mobile phone under 'Types'")
	flag.StringVar(&title, "title", "", "title of message that will be shown in the notification drawer")
	flag.StringVar(&body, "msg", "", "body or content of the notification")
	flag.IntVar(&msgid, "msgid", 0, "(optional) message-id for the message. you can optionally overwrite or clear a previous message")
	flag.StringVar(&actionURL, "url", "", "what intent should be triggered when clicking the notification, in the form of a uri")
	flag.StringVar(&imgURL, "image", "", "display an image in the notification, in the form of a uri")
	flag.StringVar(&action, "action", "send", "what kind of message to send (send = send message, clear = clear one or all messages)")
	flag.Parse()

	usr, found := users[user]
	if !found {
		log.Fatalf("user %q not found!", user)
	}

	wp := wirepusher.New()

	var msg *wirepusher.Message
	if msgid == 0 {
		msg = wirepusher.NewMsg(msgtype, title, body)
	} else {
		msg = wirepusher.MsgWithID(msgid, msgtype, title, body)
	}
	if len(actionURL) > 0 {
		msg.ActionURL = actionURL
	}
	if len(imgURL) > 0 {
		msg.ImageURL = imgURL
	}

	if strings.EqualFold(action, "clear") {
		if msgid != 0 {
			if err := wp.ClearMsg(msgid, usr); err != nil {
				log.Fatalf("error clearing message %d: %v", msgid, err)
			}
		} else {
			if err := wp.ClearAllMsg(usr); err != nil {
				log.Fatalf("error clearing all messages: %v", err)
			}
		}
		log.Print("message(s) cleared!")
	} else {
		if err := wp.Send(msg, usr); err != nil {
			log.Fatalf("error sending message: %v", err)
		}
		log.Print("message sent!")
	}
}

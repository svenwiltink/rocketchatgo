package main

import (
	"fmt"
	"github.com/svenwiltink/rocketchatgo"
	"os"
	"os/signal"
	"log"
)

func main() {
	client, err := rocketchatgo.NewClient("***")

	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Login("***", "", "***")
	if err != nil {
		fmt.Println(err)
	}

	client.AddHandler(OnMessageCreate)
	client.AddHandler(OnChannelJoin)
	client.AddHandler(OnChannelLeave)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	<-done
}

func OnMessageCreate(session *rocketchatgo.Session, event *rocketchatgo.MessageCreateEvent) {
	if session.GetUserID() == event.Message.Sender.ID{
		return
	}

	session.SendMessage(event.Message.ChannelID, "fuck it")
}

func OnChannelJoin(session *rocketchatgo.Session, event *rocketchatgo.ChannelJoinEvent) {
	log.Printf("WTF bitch, why join channel %s", event.Channel.Name)
}

func OnChannelLeave(session *rocketchatgo.Session, event *rocketchatgo.ChannelLeaveEvent) {
	log.Printf("WTF bitch, why leave channel %s", event.Channel.Name)
}
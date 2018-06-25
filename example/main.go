package main

import (
	"fmt"
	"github.com/svenwiltink/rocketchatgo"
	"os"
	"os/signal"
)

func main() {
	client, err := rocketchatgo.NewClient("xxx", 443, true)

	if err != nil {
		fmt.Println(err)
		return
	}

	err = client.Login("xxx", "", "xxx")
	if err != nil {
		fmt.Println(err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt)

	<-done
}

package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
    // --- APNA NUMBER (92...) ---
    myNumber := "923133164345" 
    // ---------------------------

	dbLog := waLog.Stdout("Database", "ERROR", true)
	container, err := sqlstore.New("sqlite3", "file:examplestore.db?_foreign_keys=on", dbLog)
	if err != nil { panic(err) }

	deviceStore, err := container.GetFirstDevice()
	client := whatsmeow.NewClient(deviceStore, waLog.Stdout("Client", "ERROR", true))

	// Message Aane Par Kya Kare?
	client.AddEventHandler(func(evt interface{}) {
		switch v := evt.(type) {
		case *events.Message:
			fmt.Println("Message Aya:", v.Message.GetConversation())
            // Yahan baad mein Jazz Drive ka logic lagega
		}
	})

	if client.Store.ID == nil {
        // Connect aur Pairing
		err = client.Connect()
		if err != nil { panic(err) }

        // Context Fix ke sath
		code, err := client.PairPhone(myNumber, true, whatsmeow.PairClientChrome, "Linux")
		if err != nil { panic(err) }

		fmt.Println("PAIRING CODE:", code)
	} else {
		err = client.Connect()
		if err != nil { panic(err) }
		fmt.Println("Bot Connected!")
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	client.Disconnect()
}

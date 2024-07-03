package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"

	"google.golang.org/api/option"
)

var app *firebase.App
func initFirebase() error {
	opt := option.WithCredentialsFile("languagebuddy-3a272-firebase-adminsdk-gtk54-4179a964ef.json")
	var err error
	app, err = firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return fmt.Errorf("error initializing app: %v", err)
	}
	return nil
}


func sendNotificationToToken(token,title,body string) error {
	// Obtain a messaging.Client from the App.
	
	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	// This registration token comes from the client FCM SDKs.

	// See documentation on defining a message payload.
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body: body,
		},
		Token: token,
	}

	// Send a message to the device corresponding to the provided
	// registration token.
	response, err := client.Send(ctx, message)
	if err != nil {
		log.Println(err)
	}else {
		// Response is a message ID string.
		fmt.Println("Successfully sent message:", response)
	}
	return err
}
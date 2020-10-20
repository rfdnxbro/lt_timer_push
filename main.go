package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {
	opt := option.WithCredentialsFile("key.json")
	config := &firebase.Config{ProjectID: "lt-timer-e8850"}
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		log.Fatalf("error getting Messaging client: %v\n", err)
	}

	fbc, err := app.Firestore(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer fbc.Close()

	topic := "brides20201030"
	jst, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now()
	formatTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), 0, 0, jst)
	oneMinuteAfterFormatTime := time.Date(now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute()+1, 0, 0, jst)

	startTimeIter := fbc.Collection("times").Doc(topic).Collection("times").Where("starts_at", "==", formatTime).Documents(context.Background())
	for {
		timeDoc, err := startTimeIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		timeData := timeDoc.Data()
		sendMessages(client, fbc, timeData["title"].(string), "スタートです", topic)
	}

	endTimeIter := fbc.Collection("times").Doc(topic).Collection("times").Where("ends_at", "==", formatTime).Documents(context.Background())
	for {
		timeDoc, err := endTimeIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		timeData := timeDoc.Data()
		sendMessages(client, fbc, timeData["title"].(string), "終了です！！！", topic)
	}

	oneStartTimeIter := fbc.Collection("times").Doc(topic).Collection("times").Where("starts_at", "==", oneMinuteAfterFormatTime).Documents(context.Background())
	for {
		timeDoc, err := oneStartTimeIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		timeData := timeDoc.Data()
		sendMessages(client, fbc, timeData["title"].(string), "あと1分でスタートです", topic)
	}

	oneEndTimeIter := fbc.Collection("times").Doc(topic).Collection("times").Where("ends_at", "==", oneMinuteAfterFormatTime).Documents(context.Background())
	for {
		timeDoc, err := oneEndTimeIter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			break
		}
		timeData := timeDoc.Data()
		sendMessages(client, fbc, timeData["title"].(string), "終了まで残り1分です", topic)
	}
}

func sendMessages(client *messaging.Client, fbc *firestore.Client, title string, body string, topic string) {
	message := &messaging.Message{
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Topic: topic,
	}

	_, err := client.Send(context.Background(), message)
	if err != nil {
		log.Fatalln(err)
	}
}

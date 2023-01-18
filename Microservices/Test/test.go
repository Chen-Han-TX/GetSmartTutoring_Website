package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

type UserScore struct {
	Score int `json:"score"`
}

func main() {
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "https://eti-assignment-2-default-rtdb.asia-southeast1.firebasedatabase.app",
	}

	// fetch service account key
	opt := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	// -------New/Update Data--------

	// create ref at path user_scores/:userId
	ref := client.NewRef("user_scores/" + fmt.Sprint(1))

	if err := ref.Set(context.TODO(), map[string]interface{}{"score": 40}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("score added/updated successfully!")

	// -------Retreiving------
	// get database reference to user score
	ref2 := client.NewRef("user_scores/1")

	// read from user_scores using ref
	var s UserScore
	if err := ref2.Get(context.TODO(), &s); err != nil {
		log.Fatalln("error in reading from firebase DB: ", err)
	}
	fmt.Println("retrieved user's score is: ", s.Score)

	// ----DELETE-----
	/*
		ref3 := client.NewRef("user_scores/1")

		if err := ref3.Delete(context.TODO()); err != nil {
			log.Fatalln("error in deleting ref: ", err)
		}
		fmt.Println("user's score deleted successfully:)")
	*/

}

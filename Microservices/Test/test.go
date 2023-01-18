package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/option"
)

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

	// create ref at path user_scores/:userId
	ref := client.NewRef("user_scores/" + fmt.Sprint(1))

	if err := ref.Set(context.TODO(), map[string]interface{}{"score": 40}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("score added/updated successfully!")

}

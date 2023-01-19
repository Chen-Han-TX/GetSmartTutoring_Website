package main

import (
	"context"
	"fmt"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type UserScore struct {
	Score int `json:"score"`
}

func RealTimeDatabaseExample() {
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

func FirestoreExample() {
	// https://firebase.google.com/docs/firestore/quickstart#go

	ctx := context.Background()

	// Use a service account
	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// ---- Add data ----

	// Create new collection and a document
	/*
		_, _, err := client.Collection("users").Add(ctx, map[string]interface{}{
			"first": "Ada",
			"last":  "Lovelace",
			"born":  1815,
		})
		if err != nil {
			log.Fatalf("Failed adding alovelace: %v", err)
		}
	*/

	// Add another document
	_, _, err = client.Collection("User").Add(ctx, map[string]interface{}{
		"Name":     "Yuxuan",
		"Email":    "xuange@gmail.com",
		"Password": "12345678",
		"UserID":   1912,
	})
	if err != nil {
		log.Fatalf("Failed adding aturing: %v", err)
	}

	// ----- Read Data -----
	iter := client.Collection("User").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		fmt.Println(doc.Data())
	}

}

func AuthenticationExample() {
	// https://firebase.google.com/docs/auth/admin/manage-users

	ctx := context.Background()

	uid := "OpoyFtGk74TZ0ZQgK2tkYTlBCJ33"

	// Initialize default app
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		log.Fatalf("error initializing app: %v\n", err)
	}

	// Access auth service from the default app
	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting Auth client: %v\n", err)
	}

	// ---- Get User by uid -----

	u, err := client.GetUser(ctx, uid)
	if err != nil {
		log.Fatalf("error getting user %s: %v\n", uid, err)
	}
	fmt.Println("Email:", u.Email)

	// ---- Create a new user ----
	params := (&auth.UserToCreate{}).
		Email("benlow2@gmail.com").
		Password("12345678")

	newUser, err := client.CreateUser(ctx, params)
	if err != nil {
		log.Fatalf("error creating user: %v\n", err)
	}
	fmt.Println("New user with id: ", newUser.UID)

}

func main() {

	//RealTimeDatabaseExample()
	//FirestoreExample()
	AuthenticationExample()

}

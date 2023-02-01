package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

// Struct for JWT Token stored in the Cookie
type Claims struct {
	EmailAddress string `json:"email_address"`
	UserType     string `json:"user_type"`
	UserID       string `json:"user_id"`
	jwt.RegisteredClaims
}

type User struct {
	UserID         string              `json:"user_id" firestore:"UserID"`
	UserType       string              `json:"user_type" firestore:"UserType"`
	Name           string              `json:"name" firestore:"Name"`
	Email          string              `json:"email" firestore:"Email"`
	Password       string              `json:"password" firestore:"Password"`
	School         string              `json:"school,omitempty" firestore:"School,omitempty"`
	AreaOfInterest map[string][]string `json:"area_of_interest" firestore:"AreaOfInterest"`
	CertOfEvidence []string            `json:"cert_of_evidence,omitempty" firestore:"CertOfEvidence,omitempty"`
}

type Message struct {
	SenderID    string    `json:"sender_id"`
	RecipientID string    `json:"recipient_id"`
	Message     string    `json:"message"`
	Timestamp   time.Time `json:"timestamp"`
}

var jwtKey = []byte("lhdrDMjhveyEVcvYFCgh1dBR2t7GM0YJ") // PLEASE DO NOT SHARE

func verifyJWT(w http.ResponseWriter, r *http.Request) (Claims, error) {
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return Claims{}, err
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return Claims{}, err
	}

	// Get the JWT string from the cookie
	tknStr := c.Value
	// Initialize a new instance of `Claims`
	claims := &Claims{}
	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return *claims, err
		}
		w.WriteHeader(http.StatusBadRequest)
		return *claims, err
	}

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return *claims, err
	}
	// Token is valid

	return *claims, nil
}

func main() {
	router := mux.NewRouter()

	//create a handler for the route for getting the chat list for the user
	router.HandleFunc("/api/user/messages", GetMessages).Methods("GET, OPTIONS")
	router.HandleFunc("/api/user/chatting", getChatList).Methods("GET","OPTIONS")
	//post a message to the chat
	router.HandleFunc("/api/user/chatting", postMessage).Methods("POST","OPTIONS")



	fmt.Println("Listening at port 5070")
	log.Fatal(http.ListenAndServe(":5070", router))

}

func GetMessages(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the cookie
	claims, err := verifyJWT(w, r)
	if err != nil {
		return
	}
	//get the messages for the user
	//get the user id from the claims
	userID := claims.UserID
	//get the recipient id from the request
	recipientID := r.URL.Query().Get("recipient_id")
	//get the messages from the database
	ctx := context.Background()
	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

	// conf := &firebase.Config{
	// 	ProjectID: "project-id",
	// }
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}


	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()
	//get the messages from the database
	messages := []Message{}
	iter := client.Collection("Messages").Where("SenderID", "==", userID).Where("RecipientID", "==", recipientID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var message Message
		doc.DataTo(&message)
		messages = append(messages, message)
	}
	iter = client.Collection("Messages").Where("SenderID", "==", recipientID).Where("RecipientID", "==", userID).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var message Message
		doc.DataTo(&message)
		messages = append(messages, message)
	}
	//sort the messages by timestamp
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Timestamp.Before(messages[j].Timestamp)
	})
	//send the messages to the user
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}

func postMessage(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the cookie
	claims, err := verifyJWT(w, r)
	if err != nil {
		return
	}
	//get the messages for the user
	//get the user id from the claims
	userID := claims.UserID
	//get the recipient id from the request
	recipientID := r.URL.Query().Get("recipient_id")
	//get the message from the request
	message := r.URL.Query().Get("message")
	//get the messages from the database
	ctx := context.Background()
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
	//get the messages from the database
	_, _, err = client.Collection("Messages").Add(ctx, map[string]interface{}{
		"SenderID":    userID,
		"RecipientID": recipientID,
		"Message":     message,
		"Timestamp":   time.Now(),
	})
	if err != nil {
		log.Fatalf("Failed adding alovelace: %v", err)
	}
}

//for getting the chat name for the user, the user must be logged in, return the chat list from firebase firestore
func getChatList(w http.ResponseWriter, r *http.Request) {
	//check if the user is logged in
	claims, err := verifyJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	//get the user id from the claims
	userID := claims.UserID

	//get the user type from the claims
	userType := claims.UserType

	//get the user from the firestore
	user, err := getUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//get the chat list from the user
	chatList := user.AreaOfInterest[userType]

	//get the chat list from the firestore
	chats, err := getChats(chatList)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//return the chat list
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(chats)
}



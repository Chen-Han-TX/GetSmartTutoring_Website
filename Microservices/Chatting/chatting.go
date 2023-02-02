package main

import (
	"context"
	"fmt"

	// "encoding/json"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/golang-jwt/jwt/v4"
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

type ChatList struct {
	StudentID string            `json:"student_id" firestore:"StudentID"`
	TutorID   string            `json:"tutor_id" firestore:"TutorID"`
	Messages  map[string]string `json:"messages" firestore:"Messages"`
}

type Application struct {
	StudentID        string `json:"student_id" firestore:"StudentID"`
	TutorID          string `json:"tutor_id" firestore:"TutorID"`
	Subject          string `json:"subject" firestore:"Subject"`
	ApplicatonStatus string `json:"application_status" firestore:"ApplicationStatus"`
	SessionLength    int    `json:"session_length" firestore:"SessionLength"`
	HourlyRate       int    `json:"hourly_rate" firestore:"HourlyRate"`
}

// var jwtKey = []byte("lhdrDMjhveyEVcvYFCgh1dBR2t7GM0YJ") // PLEASE DO NOT SHARE

// func verifyJWT(w http.ResponseWriter, r *http.Request) (Claims, error) {
// 	c, err := r.Cookie("token")
// 	if err != nil {
// 		if err == http.ErrNoCookie {
// 			// If the cookie is not set, return an unauthorized status
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return Claims{}, err
// 		}
// 		// For any other type of error, return a bad request status
// 		w.WriteHeader(http.StatusBadRequest)
// 		return Claims{}, err
// 	}

// 	// Get the JWT string from the cookie
// 	tknStr := c.Value
// 	// Initialize a new instance of `Claims`
// 	claims := &Claims{}
// 	// Parse the JWT string and store the result in `claims`.
// 	// Note that we are passing the key in this method as well. This method will return an error
// 	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
// 	// or if the signature does not match
// 	tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil {
// 		if err == jwt.ErrSignatureInvalid {
// 			w.WriteHeader(http.StatusUnauthorized)
// 			return *claims, err
// 		}
// 		w.WriteHeader(http.StatusBadRequest)
// 		return *claims, err
// 	}

// 	if !tkn.Valid {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return *claims, err
// 	}
// 	// Token is valid

// 	return *claims, nil
// }

func main() {
	router := mux.NewRouter()

	//create a handler for the route for getting the chat list for the user
	// router.HandleFunc("/api/user/messages", GetMessages).Methods("GET, OPTIONS")
	// router.HandleFunc("/api/user/chatting", getChatList).Methods("GET","OPTIONS")
	// router.HandleFunc("/api/user/chatting", postMessage).Methods("POST","OPTIONS")
	// router.HandleFunc("/api/user/messages", GetMessages).Methods("GET, OPTIONS")
	// router.HandleFunc("/api/user/chatting", getChatList).Methods("GET","OPTIONS")
	// router.HandleFunc("/api/user/chatting", postMessage).Methods("POST","OPTIONS")


	router.HandleFunc("/api/chatlist", createChatList).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/chatlist", createChatList).Methods("POST", "OPTIONS")
	fmt.Println("Listening at port 5070")
	log.Fatal(http.ListenAndServe(":5070", router))

}

// func GetMessages(w http.ResponseWriter, r *http.Request) {
// 	// Get the JWT token from the cookie
// 	claims, err := verifyJWT(w, r)
// 	if err != nil {
// 		return
// 	}
// 	//get the messages for the user
// 	//get the user id from the claims
// 	userID := claims.UserID
// 	//get the recipient id from the request
// 	recipientID := r.URL.Query().Get("recipient_id")
// 	//get the messages from the database
// 	ctx := context.Background()
// 	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

// 	app, err := firebase.NewApp(ctx, nil, sa)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	client, err := app.Firestore(ctx)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer client.Close()
// 	//get the messages from the database
// 	messages := []Message{}
// 	iter := client.Collection("Messages").Where("SenderID", "==", userID).Where("RecipientID", "==", recipientID).Documents(ctx)
// 	for {
// 		doc, err := iter.Next()
// 		if err != nil {
// 			break
// 		}
// 		var message Message
// 		doc.DataTo(&message)
// 		messages = append(messages, message)
// 	}
// 	iter = client.Collection("Messages").Where("SenderID", "==", recipientID).Where("RecipientID", "==", userID).Documents(ctx)
// 	for {
// 		doc, err := iter.Next()
// 		if err != nil {
// 			break
// 		}
// 		var message Message
// 		doc.DataTo(&message)
// 		messages = append(messages, message)
// 	}
// 	//sort the messages by timestamp
// 	sort.Slice(messages, func(i, j int) bool {
// 		return messages[i].Timestamp.Before(messages[j].Timestamp)
// 	})
// 	//send the messages to the user
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(messages)
// }

// func postMessage(w http.ResponseWriter, r *http.Request) {
// 	// Get the JWT token from the cookie
// 	claims, err := verifyJWT(w, r)
// 	if err != nil {
// 		return
// 	}
// 	//get the messages for the user
// 	//get the user id from the claims
// 	userID := claims.UserID
// 	//get the recipient id from the request
// 	recipientID := r.URL.Query().Get("recipient_id")
// 	//get the message from the request
// 	message := r.URL.Query().Get("message")
// 	//get the messages from the database
// 	ctx := context.Background()
// 	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

// 	app, err := firebase.NewApp(ctx, nil, sa)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}

// 	client, err := app.Firestore(ctx)
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	defer client.Close()
// 	//get the messages from the database
// 	_, _, err = client.Collection("Messages").Add(ctx, map[string]interface{}{
// 		"SenderID":    userID,
// 		"RecipientID": recipientID,
// 		"Message":     message,
// 		"Timestamp":   time.Now(),
// 	})
// 	if err != nil {
// 		log.Fatalf("Failed adding alovelace: %v", err)
// 	}
// }

func createChatList(w http.ResponseWriter, r *http.Request) {
	// Get the JWT token from the cookie
	// claims, err := verifyJWT(w, r)
	// if err != nil {
	// 	return
	// }
	//get the user id from the claims
	// userID := claims.UserID
	//get the messages from the database
	ctx := context.Background()
	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err.Error())
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer client.Close()
	//get the messages from the database
	chatList := []ChatList{}
	iter := client.Collection("Application").Where("ApplicationStatus", "==", "success").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		var application Application
		doc.DataTo(&application)
		chatList = append(chatList, ChatList{StudentID: application.StudentID, TutorID: application.TutorID})
	}

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200
		return
	} else if r.Method == "POST" {
		//check if the chatlist already exists in firestore, if not create a new chatlist
		for _, chat := range chatList {
			iter := client.Collection("ChatList").Where("StudentID", "==", chat.StudentID).Where("TutorID", "==", chat.TutorID).Documents(ctx)
			if iter == nil {
				_, _, err = client.Collection("ChatList").Add(ctx, map[string]interface{}{
					"StudentID": chat.StudentID,
					"TutorID":   chat.TutorID,
				})
				if err != nil {
					log.Fatalf("Failed adding alovelace: %v", err)
				}
			}
		}
	}

}

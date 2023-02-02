package main

import (
	"context"
	"encoding/json"
	"fmt"

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
	StudentID string `json:"student_id" firestore:"StudentID"`
	TutorID   string `json:"tutor_id" firestore:"TutorID"`
	//message array contains senderID content and timestamp, can be null
	Messages []map[string]Message `json:"messages" firestore:"Messages"`
}

type Message struct {
	SenderID string `json:"sender_id" firestore:"SenderID"`
	Content  string `json:"content" firestore:"Content"`
	Time     string `json:"time" firestore:"Time"`
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

	router.HandleFunc("/api/chatlist", createChatList).Methods("POST", "OPTIONS")

	fmt.Println("Listening at port 5070")
	log.Fatal(http.ListenAndServe(":5070", router))

}

// create a chat for a student and a tutor once the applicaton became success
func createChatList(w http.ResponseWriter, r *http.Request) {
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

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200
		return
	} else if r.Method == "POST" {
		//get the messages from the database
		chatList := []ChatList{}

		var chatListData ChatList
		//chatlist contains tutor id, student id and message array, message array contains senderID content and timestamp
		iter := client.Collection("Applications").Where("ApplicationStatus", "==", "Success").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err != nil {
				break
			}
			var application Application
			doc.DataTo(&application)
			chatListData.TutorID = application.TutorID
			chatListData.StudentID = application.StudentID
			chatListData.Messages = []map[string]Message{}
			chatList = append(chatList, chatListData)
		}
		//send the messages to the user
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(chatListData)

		for _, chat := range chatList {
			//check the chatlist in firebase if the tutorID and studentid is already in the chatlist, if not add it to the chatlist
			iter := client.Collection("ChatList").Where("TutorID", "==", chat.TutorID).Where("StudentID", "==", chat.StudentID).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatList ChatList
				doc.DataTo(&chatList)
				if chatList.TutorID == chat.TutorID && chatList.StudentID == chat.StudentID {
					break
				}
			}
			//add the chatlist to the firebase
			_, _, err := client.Collection("ChatList").Add(ctx, chat)
			if err != nil {
				log.Fatalf("Failed adding chatlist: %v", err)
			}

		}
	}
}

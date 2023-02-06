package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/api/option"
)

var cred_file = "/eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json"

// var url = "https://react-app-4dcnj7fm6a-uc.a.run.app
var url = "http://104.154.110.27"

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
	ChatId      string `json:"chat_id" firestore:"ChatID"`
	StudentID   string `json:"student_id" firestore:"StudentID"`
	TutorID     string `json:"tutor_id" firestore:"TutorID"`
	StudentName string `json:"student_name" firestore:"StudentName"`
	TutorName   string `json:"tutor_name" firestore:"TutorName"`
	//message array contains senderID content and timestamp, can be null
	Messages []Message `json:"messages" firestore:"Messages"`
}

type Message struct {
	SenderID string `json:"sender_id" firestore:"SenderID"`
	Content  string `json:"content" firestore:"Content"`
	//unix timestamp
	Timestamp int64 `json:"timestamp" firestore:"Timestamp"`
}

type Application struct {
	SessionID        string `json:"session_id" firestore:"SessionID"`
	StudentID        string `json:"student_id" firestore:"StudentID"`
	StudentName      string `json:"student_name" firestore:"StudentName"`
	TutorID          string `json:"tutor_id" firestore:"TutorID"`
	TutorName        string `json:"tutor_name" firestore:"TutorName"`
	Subject          string `json:"subject" firestore:"Subject"`
	ApplicatonStatus string `json:"application_status" firestore:"ApplicationStatus"`
	SessionLength    int    `json:"session_length" firestore:"SessionLength"`
	HourlyRate       int    `json:"hourly_rate" firestore:"HourlyRate"`
}

func main() {
	router := mux.NewRouter()

	// Create a chat list if the tutor accepts the application (when the application status is changed to "accepted")
	router.HandleFunc("/api/createchatlist", createChatList).Methods("POST", "OPTIONS")
	// For the user, retrieve all the chatlist that he is involved in
	router.HandleFunc("/api/getlist/{user_id}/{user_type}", getChatList).Methods("GET", "OPTIONS")
	// For the user, retrieve all the messages for a specific chat in a chatlist (Chen Han: we didnt really use this in front-end so..)
	// router.HandleFunc("/api/getmessages/{userid_opp}", getMessages).Methods("GET", "OPTIONS")
	// For the user, send a message to a chat in a chatlist
	router.HandleFunc("/api/sendmessages/{user_id}/{userid_opp}/{user_type}", sendMessage).Methods("POST", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{url},
		AllowCredentials: true,
	})

	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)

	fmt.Println("Listening at port 5053")
	log.Fatal(http.ListenAndServe(":5053", handler))

}

// create a chat for a student and a tutor once the applicaton became success
func createChatList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", url)

	ctx := context.Background()
	sa := option.WithCredentialsFile(cred_file)

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
		iter := client.Collection("Applications").Where("ApplicationStatus", "==", "Accepted").Documents(ctx)
		for {
			doc, err := iter.Next()
			if err != nil {
				break
			}
			var application Application
			doc.DataTo(&application)
			chatListData.TutorID = application.TutorID
			chatListData.StudentID = application.StudentID
			chatListData.StudentName = application.StudentName
			chatListData.TutorName = application.TutorName
			chatListData.Messages = []Message{}
			chatList = append(chatList, chatListData)
		}
		/*
			//send the messages to the user
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(chatListData)
		*/

		for _, chat := range chatList {
			// Check if the chatList entry already exists in the database
			found := false
			iter := client.Collection("ChatList").Where("TutorID", "==", chat.TutorID).Where("StudentID", "==", chat.StudentID).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatList ChatList
				doc.DataTo(&chatList)
				if chatList.TutorID == chat.TutorID && chatList.StudentID == chat.StudentID {
					found = true
					break
				}
			}
			// Add the chatList to the database if it does not already exist
			if !found {

				ref := client.Collection("ChatList").NewDoc()
				chat.ChatId = ref.ID

				//_, _, err2 := client2.Collection("Applications").Add(ctx, xApplication)
				_, err := ref.Set(ctx, chat)
				if err != nil {
					log.Fatalf("Failed adding chatlist: %v", err)
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(chat)
			}
		}
	}
}

func getChatList(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", url)
	w.Header().Set("Allow-Control-Allow-Origin", url)

	params := mux.Vars(r)
	user_id := params["user_id"]
	user_type := params["user_type"]

	ctx := context.Background()
	sa := option.WithCredentialsFile(cred_file)

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
	} else if r.Method == "GET" {
		chatList := []ChatList{}
		//get the messages from the database, user data from jwt
		if user_type == "Tutor" {
			iter := client.Collection("ChatList").Where("TutorID", "==", user_id).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatListData ChatList
				doc.DataTo(&chatListData)
				chatList = append(chatList, chatListData)
			}
		}
		if user_type == "Student" {
			iter := client.Collection("ChatList").Where("StudentID", "==", user_id).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatListData ChatList
				doc.DataTo(&chatListData)
				chatList = append(chatList, chatListData)
			}
		}
		if len(chatList) == 0 {
			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "No chatlist found")
			return
		} else {
			// w.Header().Set("Content-Type", "application/json")
			// w.WriteHeader(http.StatusAccepted)
			json.NewEncoder(w).Encode(chatList)
		}
	}
}

/*
// write the getmessage function, get all the messages from the firebase firestore in Message struct. Sort the messages by timestamp
func getMessages(w http.ResponseWriter, r *http.Request) {
	claims, _ := verifyJWT(w, r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotAcceptable)
	// 	fmt.Println(err.Error())
	// }
	userid := claims.UserID

	vars := mux.Vars(r)
	anotherUserId := vars["userid_opp"]

	ctx := context.Background()
	sa := option.WithCredentialsFile(cred_file)

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
	} else if r.Method == "GET" {
		messages := []Message{}
		//get the messages from the database, user data from jwt
		if claims.UserType == "Tutor" {
			iter := client.Collection("ChatList").Where("TutorID", "==", userid).Where("StudentID", "==", anotherUserId).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatList ChatList
				doc.DataTo(&chatList)
				//messages is an array under ChatList
				messages = chatList.Messages
			}

			if len(messages) == 0 {
				w.Header().Set("Content-type", "text/plain")
				// w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "No messages found")
				return
			} else {
				// w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusAccepted)
				//display the messages in the array
				json.NewEncoder(w).Encode(messages)
			}
		} else if claims.UserType == "Student" {
			iter := client.Collection("ChatList").Where("StudentID", "==", userid).Where("TutorID", "==", anotherUserId).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatList ChatList
				doc.DataTo(&chatList)
				//messages is an array under ChatList
				messages = chatList.Messages
				//sort the messages by timestamp, latest first
				sort.Slice(messages, func(i, j int) bool {
					return messages[i].Timestamp > messages[j].Timestamp
				})

			}

			if len(messages) == 0 {
				w.Header().Set("Content-type", "text/plain")
				// w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "No messages found")
				return
			} else {
				// w.Header().Set("Content-Type", "application/json")
				// w.WriteHeader(http.StatusAccepted)
				//display the messages in the array
				json.NewEncoder(w).Encode(messages)
			}
		}
	}
}
*/

func sendMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", url)

	vars := mux.Vars(r)
	user_id := vars["user_id"]
	anotherUserId := vars["userid_opp"]
	user_type := vars["user_type"]

	ctx := context.Background()
	sa := option.WithCredentialsFile(cred_file)

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
		var message Message
		_ = json.NewDecoder(r.Body).Decode(&message)
		// message.Content = strings.TrimSpace(message.Content)
		message.SenderID = user_id
		message.Timestamp = time.Now().Unix()
		//if the user is a student
		if user_type == "Student" {
			iter := client.Collection("ChatList").Where("StudentID", "==", user_id).Where("TutorID", "==", anotherUserId).Documents(ctx)

			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatList ChatList
				doc.DataTo(&chatList)
				chatList.Messages = append(chatList.Messages, message)
				_, err = client.Collection("ChatList").Doc(doc.Ref.ID).Set(ctx, chatList)
				if err != nil {
					log.Fatalf("Failed adding chatlist: %v", err)
				}
			}
		}
		//if the user is a tutor
		if user_type == "Tutor" {
			iter := client.Collection("ChatList").Where("TutorID", "==", user_id).Where("StudentID", "==", anotherUserId).Documents(ctx)
			for {
				doc, err := iter.Next()
				if err != nil {
					break
				}
				var chatList ChatList
				doc.DataTo(&chatList)
				chatList.Messages = append(chatList.Messages, message)
				_, err = client.Collection("ChatList").Doc(doc.Ref.ID).Set(ctx, chatList)
				if err != nil {
					log.Fatalf("Failed adding chatlist: %v", err)
				}
			}
		}
	}
}

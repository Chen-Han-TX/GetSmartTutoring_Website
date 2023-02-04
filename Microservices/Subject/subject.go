package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	firebase "firebase.google.com/go/v4"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// http://localhost:5053/api/getsubjects/all
// {"A-level": ["a","b"], "O-level": [...]..}
// http://localhost:5053/api/getsubjects/psle
// ["a","b",..]
// http://localhost:5053/api/getsubjects/olevel
// http://localhost:5053/api/getsubjects/alevel

func Subject(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	req_type := params["type"]

	ctx := context.Background()

	// Use a service account
	sa := option.WithCredentialsFile("/eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	defer client.Close()

	// get all the data
	iter := client.Collection("Global Data").Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}

		data := doc.Data()
		if data != nil {
			if req_type == "all" {
				json.NewEncoder(w).Encode(data)
				return
			} else if req_type == "psle" {
				for k := range data {
					if k == "PSLE" {
						json.NewEncoder(w).Encode(data[k])
					}
				}
				return
			} else if req_type == "olevel" {
				for k := range data {
					if k == "O-Level" {
						json.NewEncoder(w).Encode(data[k])
					}
				}
				return
			} else if req_type == "alevel" {
				for k := range data {
					if k == "A-Level" {
						json.NewEncoder(w).Encode(data[k])
					}
				}
				return
			} else {
				w.WriteHeader(http.StatusNotAcceptable) // 406
				return
			}
		}
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/getsubjects/{type}", Subject).Methods("GET", "OPTIONS")

	fmt.Println("Listening at port 5051")
	log.Fatal(http.ListenAndServe(":5051", router))
}

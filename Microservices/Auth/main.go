package main

import (
	//"encoding/json"

	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/api/option"
)

var cred_file = "/eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json"

// var url = "https://react-app-4dcnj7fm6a-uc.a.run.app
// var url = "http://localhost:3000"
// var url = "http://192.168.1.4:80"
var url = "http://104.154.110.27"

// ==== STRUCTs ========
// struct for user login credentials
type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Struct for User object
type User struct {
	UserID         string                       `json:"user_id" firestore:"UserID"`
	UserType       string                       `json:"user_type" firestore:"UserType"`
	Name           string                       `json:"name" firestore:"Name"`
	Email          string                       `json:"email" firestore:"Email"`
	Password       string                       `json:"password" firestore:"Password"`
	AreaOfInterest map[string][]string          `json:"area_of_interest" firestore:"AreaOfInterest"`
	School         string                       `json:"school,omitempty" firestore:"School,omitempty"`
	HourlyRate     int                          `json:"hourly_rate,omitempty" firestore:"HourlyRate,omitempty"`
	Availability   map[string]map[string]string `json:"availability,omitempty" firestore:"Availability,omitempty"`
	CertOfEvidence []string                     `json:"cert_of_evidence,omitempty" firestore:"CertOfEvidence,omitempty"`
}

// ========= HANDLER FUNCTIONS ==========

// RETURN 200 -> Registered
// RETURN 406 -> Duplicated account (email)
func SignUp(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", url)

	// POST http://localhost:5050/api/auth/signup/student
	// {"name": "xyz", "email": "..", "password", "area of interest": {"olevel":"..."...}, "certificate":[]}
	// POST http://localhost:5050/api/auth/signup/tutor
	//{"name": "xyz", "email": "..", "password", "area of interest": {"olevel":"..."...}, "certificate":"..."}

	params := mux.Vars(r)
	user_type := params["user_type"]

	ctx := context.Background()
	sa := option.WithCredentialsFile(cred_file)

	// ---Authentication--
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
	}

	// Access auth service from the default app
	client, err := app.Auth(ctx)
	if err != nil {
		fmt.Printf("error getting Auth client: %v\n", err)
	}

	// ----Firestore----
	app2, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Println(err.Error())
	}

	client2, err := app2.Firestore(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer client2.Close()

	// new user
	var user User

	// Check req methods
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "POST" {
		fmt.Println(r.Body)
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			panic(err.Error())
		}

		// Step 1: Check if student or tutor
		if user_type == "student" {
			user.UserType = "Student"
		} else if user_type == "tutor" {
			user.UserType = "Tutor"
		} else {
			w.WriteHeader(http.StatusNotFound) // 404
			return
		}

		// ---- Create a new Auth user ----
		params := (&auth.UserToCreate{}).
			Email(user.Email).
			Password(user.Password)

		newUser, err := client.CreateUser(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusNotAcceptable)
			fmt.Println(err.Error())
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		fmt.Println("New user with id: ", newUser.UID)
		user.UserID = newUser.UID

		// ---- Create a new Firebase record
		// Set UserID as the document name
		success, err := client2.Collection("User").Doc(user.UserID).Set(ctx, user)
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			fmt.Printf("An error has occurred: %s", err)
		}
		fmt.Println(success)
		json.NewEncoder(w).Encode(user)

	} else { // Other req methods
		w.WriteHeader(http.StatusNotFound) //404
		return
	}
}

func Login(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", url)

	var creds Credentials
	var UserID string
	var user User

	ctx := context.Background()
	sa := option.WithCredentialsFile(cred_file)

	// ---Authentication--
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
	}

	// Access auth service from the default app
	client, err := app.Auth(ctx)
	if err != nil {
		fmt.Printf("error getting Auth client: %v\n", err)
	}

	app2, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Println(err.Error())
	}

	client2, err := app2.Firestore(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer client2.Close()

	// Check req methods
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200
		return
	} else if r.Method == "POST" {
		// Receive user login information in JSON
		// and decode into User
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		// verify user email and password

		u, err := client.GetUserByEmail(ctx, creds.Email)
		if err != nil {
			fmt.Printf("Error getting user: %v\n", err)
			return
		}

		UserID = u.UID

		// Get user password and verify (For now, do this way)
		dsnap, err := client2.Collection("User").Doc(UserID).Get(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dsnap.DataTo(&user)

		if creds.Password != user.Password {
			w.WriteHeader(http.StatusNotAcceptable) //406
			json.NewEncoder(w).Encode("Password not matched!")
			return
		}

		// Return user object
		json.NewEncoder(w).Encode(user)
		return

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

}

// TEST - Check Cookie JWT token and return something
func GetUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Access-Control-Allow-Origin", url)

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200
		return
	} else if r.Method == "GET" {

		params := mux.Vars(r)
		user_id := params["user_id"]

		ctx := context.Background()
		sa := option.WithCredentialsFile(cred_file)

		// ----Firestore----
		app2, err := firebase.NewApp(ctx, nil, sa)
		if err != nil {
			fmt.Println(err.Error())
		}

		client2, err := app2.Firestore(ctx)
		if err != nil {
			fmt.Println(err.Error())
		}

		defer client2.Close()

		// new user
		var user User

		dsnap, err := client2.Collection("User").Doc(user_id).Get(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dsnap.DataTo(&user)
		user.Password = ""

		json.NewEncoder(w).Encode(user)

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}

}

/*
func Logout(w http.ResponseWriter, r *http.Request) {
	// immediately clear the token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Expires:  time.Now(),
		Path:     "/",
		HttpOnly: true,
	})
	fmt.Println(r.Header)
}
*/

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/auth/signup/{user_type}", SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/login", Login).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/get/{user_id}", GetUser).Methods("GET", "OPTIONS")
	//router.HandleFunc("/api/auth/refresh", Refresh)
	//router.HandleFunc("/api/auth/logout", Logout).Methods("GET", "OPTIONS")
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{url},
		AllowCredentials: true,
	})

	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)

	fmt.Println("Listening at port 5050")
	log.Fatal(http.ListenAndServe(":5050", handler))
}

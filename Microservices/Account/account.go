package main

import (
	//"encoding/json"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

// ==== STRUCTs ========
// Struct for User object
type User struct {
	UserID         string              `json:"user_id" firestore:"UserID"`
	UserType       string              `json:"user_type" firestore:"UserType"`
	Name           string              `json:"name" firestore:"Name"`
	Email          string              `json:"email" firestore:"Email"`
	Password       string              `json:"password" firestore:"Password"`
	AreaOfInterest map[string][]string `json:"area_of_interest" firestore:"AreaOfInterest"`
	CertOfEvidence []string            `json:"cert_of_evidence,omitempty" firestore:"CertOfEvidence,omitempty"`
}

// Struct for JWT Token stored in the Cookie
type Claims struct {
	EmailAddress string `json:"email_address"`
	UserType     string `json:"user_type"`
	UserID       string `json:"user_id"`
	jwt.RegisteredClaims
}

// ====== GLOBAL VARIABLES ========
var jwtKey = []byte("lhdrDMjhveyEVcvYFCgh1dBR2t7GM0YK") // A secure JWT Token for decoding, DO NOT SHARE

// ====== FUNCTONS =========
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

// ======= HANDLER FUNCTIONS ========
// GET user info
// UPDATE user in the db
func GetUser(w http.ResponseWriter, r *http.Request) {
	// Verify the JWT Token and get the Claim info
	claims, err := verifyJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		panic(err.Error())
	}

	// Variables
	user_id := claims.UserID
	var user User

	// Init connection to firestore
	ctx := context.Background()
	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		fmt.Println(err.Error())
	}

	client, err := app.Firestore(ctx)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer client.Close()

	// Check req method
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) //200
		return
	} else if r.Method == "GET" {
		dsnap, err := client.Collection("User").Doc(user_id).Get(ctx)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		dsnap.DataTo(&user)
		fmt.Printf("Document data: %#v\n", user)

		// Hide the user password
		user.Password = ""

		// Return response
		json.NewEncoder(w).Encode(user)
		return

	} else if r.Method == "POST" {
		return

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

// This handler update the user's password regardless of user type
func UpdatePassword(w http.ResponseWriter, r *http.Request) {
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/user", GetUser).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/api/user/password", UpdatePassword).Methods("POST", "OPTIONS")
	//router.HandleFunc("/api/user/getuser", GetUser).Methods("GET", "PUT", "OPTIONS")
	//router.HandleFunc("/api/user/password", UpdatePassword).Methods("PUT", "OPTIONS")

	fmt.Println("Listening at port 5051")
	log.Fatal(http.ListenAndServe(":5051", router))

}

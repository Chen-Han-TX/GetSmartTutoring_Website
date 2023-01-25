package main

import (
	//"encoding/json"
	//"encoding/json"
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
)

type SubjectOfInterest struct {
	SubjectID    int    `json:"subject_id"`
	SubjectName  string `json:"subject_name"`
	SubjectLevel string `json:"subject_level"`
}

// Struct for JWT Token stored in the Cookie
type Claims struct {
	EmailAddress string `json:"email_address"`
	UserType     string `json:"user_type"`
	UserID       string `json:"user_id"`
	jwt.RegisteredClaims
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

	router.HandleFunc("/api/user/matchtutors", matchTutors).Methods("GET", "POST", "OPTIONS")
	router.HandleFunc("/api/user/handletutorrequest", handleTutorRequest).Methods("POST", "OPTIONS")
	//router.HandleFunc("/api/user/getuser", GetUser).Methods("GET", "PUT", "OPTIONS")
	//router.HandleFunc("/api/user/password", UpdatePassword).Methods("PUT", "OPTIONS")

	fmt.Println("Listening at port 5052")
	log.Fatal(http.ListenAndServe(":5052", router))

}

func matchTutors(w http.ResponseWriter, r *http.Request) {

}

func handleTutorRequest(w http.ResponseWriter, r *http.Request) {

}

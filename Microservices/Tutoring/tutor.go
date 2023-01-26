package main

import (
	//"encoding/json"
	//"encoding/json"
	"context"
	"encoding/json"
	"fmt"

	//"io/ioutil"
	"log"
	"net/http"

	firebase "firebase.google.com/go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"google.golang.org/api/iterator"
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

	fmt.Println("Listening at port 5054")
	log.Fatal(http.ListenAndServe(":5054", router))

}

// For method matchTutors():
// Example of JSON sent with no matches (so far)
// {"O-Level":["Literature", "Chinese"],
//
//	"A-Level":["Chemistry", "Malay"]}

// Example of JSON returned with 2 matches (so far)
// [
//
//	{
//	    "Email": "benlow@gmail.com",
//	    "Name": "Ben Low",
//	    "UserID": "9QvRMWvL9gUTcVllfnStT7IrHi62"
//	},
//	{
//	    "Email": "thomerson@gmail.com",
//	    "Name": "Thomerson Yang",
//	    "UserID": " Ed8DKn9WKCqHOcbJpED0"
//	}
//
// ]

func matchTutors(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("../eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

	// // ---Authentication--
	// app, err := firebase.NewApp(ctx, nil, sa)
	// if err != nil {
	// 	fmt.Printf("error initializing app: %v\n", err)
	// }

	// // Access auth service from the default app
	// client, err := app.Auth(ctx)
	// if err != nil {
	// 	fmt.Printf("error getting Auth client: %v\n", err)
	// }

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

	var areaOfInterests map[string][]string
	var allTutors []User
	var matchedTutors []map[string]string
	tutors := client2.Collection("User")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200
		return
	} else if r.Method == "POST" {

		err := json.NewDecoder(r.Body).Decode(&areaOfInterests)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		keys := make([]string, 0, len(areaOfInterests))
		for k := range areaOfInterests {
			keys = append(keys, k)
		}

		fmt.Println(keys)

		query := tutors.Where("UserType", "==", "Tutor").Documents(ctx)
		for {
			var xTutor User
			doc, err := query.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatalf("Failed to iterate: %v", err)
			}
			doc.DataTo(&xTutor)
			allTutors = append(allTutors, xTutor)
		}
	tutorloop:
		for _, v := range allTutors {
			fmt.Println(v.Name)
			tutorAOIList := make([]string, 0, len(v.AreaOfInterest))
			for k := range v.AreaOfInterest {
				tutorAOIList = append(tutorAOIList, k)
			}

			studentAOIList := make([]string, 0, len(areaOfInterests))
			for k := range areaOfInterests {
				studentAOIList = append(studentAOIList, k)
			}

			for _, y := range studentAOIList {
				if contains(tutorAOIList, y) {
					tutorSubjList := v.AreaOfInterest[y]
					studentSubjList := areaOfInterests[y]

					if checkSimilar(tutorSubjList, studentSubjList) {
						tutorObj := map[string]string{
							"Name":   v.Name,
							"UserID": v.UserID,
							"Email":  v.Email,
						}
						matchedTutors = append(matchedTutors, tutorObj)
						continue tutorloop
					}
				}
			}
		}

		if len(matchedTutors) != 0 {
			w.Header().Set("Content-type", "application/json")
			res, _ := json.MarshalIndent(matchedTutors, "", "\t")
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, string(res))
		} else {
			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "Error - no tutors matched!")
		}

	}
}

func handleTutorRequest(w http.ResponseWriter, r *http.Request) {

}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func checkSimilar(listOne, listTwo []string) bool {
	for _, i := range listOne {
		for _, j := range listTwo {
			if i == j {
				return true
			}
		}
	}
	return false
}

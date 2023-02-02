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

	"cloud.google.com/go/firestore"
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

type Application struct {
	StudentID        string `json:"student_id" firestore:"StudentID"`
	TutorID          string `json:"tutor_id" firestore:"TutorID"`
	Subject          string `json:"subject" firestore:"Subject"`
	ApplicatonStatus string `json:"application_status" firestore:"ApplicationStatus"`
	SessionLength    int    `json:"session_length" firestore:"SessionLength"`
	HourlyRate       int    `json:"hourly_rate" firestore:"HourlyRate"`
}

var jwtKey = []byte("lhdrDMjhveyEVcvYFCgh1dBR2t7GM0YK") // PLEASE DO NOT SHARE

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

	router.HandleFunc("/api/tutoring/matchtutors", matchTutors).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tutoring/apply", applyForTutor).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tutoring/getapplications", getApplications).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tutoring/handleapplications", handleApplications).Methods("PUT", "OPTIONS")
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
	// Verify JWT token to continue using
	// _, err := verifyJWT(w, r)
	// if err != nil {
	// 	w.WriteHeader(http.StatusNotAcceptable) //406
	// 	json.NewEncoder(w).Encode(err.Error())
	// 	return
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
	matchedTutors := []map[string]interface{}{}

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

		//tutorloop:
		for _, v := range allTutors {
			var fullsubjlist []string
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
				var tutorSubjList []string
				var studentSubjList []string
				if contains(tutorAOIList, y) {
					fmt.Println("match aoi - " + y)
					tutorSubjList = append(tutorSubjList, v.AreaOfInterest[y]...)
					for i, v := range tutorSubjList {
						tutorSubjList[i] = y + " - " + v
					}
					studentSubjList = append(studentSubjList, areaOfInterests[y]...)
					fmt.Print("ASDJASDJ>")
					fmt.Println(studentSubjList)
					for i, v := range studentSubjList {
						studentSubjList[i] = y + " - " + v
					}

					fmt.Println(tutorSubjList)
					fmt.Println(studentSubjList)

					similar, matchedList := checkSimilar(tutorSubjList, studentSubjList)
					fmt.Print("MATCHEDLIST> ")
					fmt.Println(similar)
					fmt.Println(matchedList)

					if similar {
						fullsubjlist = append(fullsubjlist, matchedList...)
					}
				}
			}
			if len(fullsubjlist) != 0 {
				tutorObj := map[string]interface{}{
					"Name":               v.Name,
					"UserID":             v.UserID,
					"Email":              v.Email,
					"MatchedSubjectList": fullsubjlist,
					"HourlyRate":         v.HourlyRate,
					"Availability":       v.Availability,
				}
				matchedTutors = append(matchedTutors, tutorObj)
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

func applyForTutor(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("Microservices/eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

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

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK) // 200
		return
	} else if r.Method == "POST" {
		var xApplication Application
		err := json.NewDecoder(r.Body).Decode(&xApplication)
		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		_, _, err2 := client2.Collection("Applications").Add(ctx, xApplication)
		if err2 != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusConflict)
			fmt.Fprintf(w, "Error - unable to post application!")
		} else {
			w.Header().Set("Content-type", "text/plain")
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "Successfully posted application!")
		}
	}
}

func getApplications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("Microservices/eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

	claims, err := verifyJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable) //406
		json.NewEncoder(w).Encode(err.Error())
		return
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

	if claims.UserType != "Tutor" {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error - This method is only for tutors to use!")
	} else {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK) // 200
			return
		} else if r.Method == "GET" {
			var applicationList []Application
			allApplications := client2.Collection("Applications")
			query := allApplications.Where("TutorID", "==", claims.UserID).Documents(ctx)
			for {
				var xApplication Application
				doc, err := query.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalf("Failed to iterate: %v", err)
				}
				doc.DataTo(&xApplication)
				applicationList = append(applicationList, xApplication)
			}

			if len(applicationList) != 0 {
				w.Header().Set("Content-type", "application/json")
				res, _ := json.MarshalIndent(applicationList, "", "\t")
				w.WriteHeader(http.StatusAccepted)
				fmt.Fprintf(w, string(res))
			} else {
				w.Header().Set("Content-type", "text/plain")
				w.WriteHeader(http.StatusNotFound)
				fmt.Fprintf(w, "Error - no applications found!")
			}
		}
	}
}

func handleApplications(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	sa := option.WithCredentialsFile("Microservices/eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json")

	claims, err := verifyJWT(w, r)
	if err != nil {
		w.WriteHeader(http.StatusNotAcceptable) //406
		json.NewEncoder(w).Encode(err.Error())
		return
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

	if claims.UserType != "Tutor" {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error - This method is only for tutors to use!")
	} else {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK) // 200
			return
		} else if r.Method == "PUT" {
			var app Application
			err := json.NewDecoder(r.Body).Decode(&app)
			if err != nil {
				// If the structure of the body is wrong, return an HTTP error
				w.WriteHeader(http.StatusBadRequest) //400
				return
			}
			query := client2.Collection("Applications").Where("TutorID", "==", app.TutorID).Where("StudentID", "==", app.StudentID).Where("Subject", "==", app.Subject).Documents(ctx)
			for {

				doc, err := query.Next()
				if err == iterator.Done {
					break
				}
				if err != nil {
					log.Fatalf("Failed to iterate: %v", err)
				}

				_, err = doc.Ref.Update(ctx, []firestore.Update{
					{
						Path:  "ApplicationStatus",
						Value: app.ApplicatonStatus,
					},
				})
				if err != nil {
					// Handle any errors in an appropriate way, such as returning them.
					log.Printf("An error has occurred: %s", err)
					w.Header().Set("Content-type", "text/plain")
					w.WriteHeader(http.StatusConflict)
					fmt.Fprintf(w, "Error - Request went wrong...")
				} else {
					w.Header().Set("Content-type", "text/plain")
					w.WriteHeader(http.StatusAccepted)
					fmt.Fprintf(w, "Application handled successfully!")
				}
			}
		}
	}
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func checkSimilar(listOne, listTwo []string) (bool, []string) {
	var subjlist []string
	for _, i := range listOne {
		for _, j := range listTwo {
			if i == j {
				subjlist = append(subjlist, i)
			}
		}
	}
	if len(subjlist) != 0 {
		return true, subjlist
	} else {
		return false, subjlist
	}
}

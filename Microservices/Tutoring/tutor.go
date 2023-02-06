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
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

var cred_file = "eti-assignment-2-firebase-adminsdk-6r9lk-85fb98eda4.json"

// var url = "https://react-app-4dcnj7fm6a-uc.a.run.app

var url = "http://104.154.110.27"

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

	router.HandleFunc("/api/tutoring/matchtutors", matchTutors).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tutoring/apply", applyForTutor).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/tutoring/getapplications/{user_id}/{user_type}", getApplications).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/tutoring/handleapplications/{user_type}", handleApplications).Methods("POST", "OPTIONS")
	//router.HandleFunc("/api/user/getuser", GetUser).Methods("GET", "PUT", "OPTIONS")
	//router.HandleFunc("/api/user/password", UpdatePassword).Methods("PUT", "OPTIONS")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{url},
		AllowCredentials: true,
	})

	handler := cors.Default().Handler(router)
	handler = c.Handler(handler)

	fmt.Println("Listening at port 5052")
	log.Fatal(http.ListenAndServe(":5052", handler))

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
	w.Header().Set("Access-Control-Allow-Origin", url)

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
					tutorSubjList = append(tutorSubjList, v.AreaOfInterest[y]...)
					for i, v := range tutorSubjList {
						tutorSubjList[i] = y + " - " + v
					}
					studentSubjList = append(studentSubjList, areaOfInterests[y]...)
					for i, v := range studentSubjList {
						studentSubjList[i] = y + " - " + v
					}

					similar, matchedList := checkSimilar(tutorSubjList, studentSubjList)

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
	w.Header().Set("Access-Control-Allow-Origin", url)

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
		ref := client2.Collection("Applications").NewDoc()
		xApplication.SessionID = ref.ID

		//_, _, err2 := client2.Collection("Applications").Add(ctx, xApplication)
		_, err2 := ref.Set(ctx, xApplication)
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
	w.Header().Set("Access-Control-Allow-Origin", url)

	params := mux.Vars(r)
	user_id := params["user_id"]
	user_type := params["user_type"]

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

	if user_type == "Student" {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK) // 200
			return
		} else if r.Method == "GET" {
			var applicationList []Application
			allApplications := client2.Collection("Applications")
			query := allApplications.Where("StudentID", "==", user_id).Documents(ctx)
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
	} else {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK) // 200
			return
		} else if r.Method == "GET" {
			var applicationList []Application
			allApplications := client2.Collection("Applications")
			query := allApplications.Where("TutorID", "==", user_id).Documents(ctx)
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
	w.Header().Set("Access-Control-Allow-Origin", url)

	params := mux.Vars(r)
	user_type := params["user_type"]

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

	if user_type != "Tutor" {
		w.Header().Set("Content-type", "text/plain")
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Error - This method is only for tutors to use!")
	} else {
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK) // 200
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT")
			return
		} else if r.Method == "POST" {
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

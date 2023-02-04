package main

import (
	//"encoding/json"
	//"encoding/json"
	//"database/sql"

	"encoding/json"
	"fmt"

	"log"
	//"net/http"
	"context"
	"net/http"

	firestore "cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"google.golang.org/api/option"
)

type Payment struct {
	Amount    int    `json:"amount" firestore:"Amount"`
	TutorID   string `json:"tutor_id" firestore:"TutorID"`
	StudentID string `json:"student_id" firestore:"StudentID"`
	SessionID string `json:"session_id" firestore:"SessionID"`
}

/*
var (
	c           *paypal.Client
	err         error
	accessToken *paypal.TokenResponse
)
*/

func main() {

	/*
		// Create a client instance
		clientID := "ATNfsZGLPPVL7arlBKnfzN-EHuez8RCa32pQt2Oeg_F8-mg707LkxZht-zpiLNKxkV1VuIriz1IpYD7c"
		clientSecret := "EChsTM9D--MWl2wZrG_Ko9pewK8QrUBoPZc8PHemaOJrVleqQSPftbnbBxxXxvsQQNFJbM-EnD5Hg8yD"
		c, err = paypal.NewClient(clientID, clientSecret, paypal.APIBaseSandBox)
		if err != nil {
			fmt.Println(err.Error())
		}

		accessToken, err = c.GetAccessToken(context.Background())
		if err != nil {
			fmt.Println(err.Error())
		}
	*/
	router := mux.NewRouter()
	//router.HandleFunc("/api/payment/", SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/payment", GetPayment).Methods("POST", "OPTIONS")

	fmt.Println("Listening at port 5053")
	log.Fatal(http.ListenAndServe(":5053", router))
}

// ======= HANDLER FUNCTIONS ========
// GET user info
// UPDATE user in the db
func GetPayment(w http.ResponseWriter, r *http.Request) {

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

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	} else if r.Method == "POST" {

		var payment Payment

		err := json.NewDecoder(r.Body).Decode(&payment)

		fmt.Println(payment.Amount, payment.SessionID, payment.StudentID, payment.TutorID)

		if err != nil {
			// If the structure of the body is wrong, return an HTTP error
			w.WriteHeader(http.StatusBadRequest) //400
			return
		}

		_, err = client.Collection("Payment").Doc(payment.SessionID).Set(ctx, payment)

		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}

		// update the application status to Payment Made
		_, err = client.Collection("Applications").Doc(payment.SessionID).Update(ctx, []firestore.Update{
			{
				Path:  "ApplicationStatus",
				Value: "Payment Made",
			},
		})
		if err != nil {
			// Handle any errors in an appropriate way, such as returning them.
			log.Printf("An error has occurred: %s", err)
		}

		json.NewEncoder(w).Encode(payment)
		return

	} else {
		w.WriteHeader(http.StatusNotFound)
		return
	}
}

/*
// Add new trip record to database
func GetPayment(w http.ResponseWriter, r *http.Request) {

	// Create a new http client
	client := &http.Client{}

	if r.Method == "POST" {
		// Prepare the payment data
		data := map[string]interface{}{
			"intent": "sale",
			"payer": map[string]string{
				"payment_method": "paypal",
			},
			"transactions": []map[string]interface{}{
				{
					"amount": map[string]string{
						"total":    "10",
						"currency": "USD",
					},
					"description": "This is the payment transaction description.",
				},
			},
			"redirect_urls": map[string]string{
				"return_url": "http://localhost:3000/success",
				"cancel_url": "http://localhost:3000/cancel",
			},
		}

		// Convert the data to json
		jsonData, err := json.Marshal(data)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Prepare the request
		req, err := http.NewRequest("POST", "https://api.sandbox.paypal.com/v1/payments/payment", bytes.NewBuffer(jsonData))
		if err != nil {
			fmt.Println(err)
			return
		}

		// Set the headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+accessToken.Token)

		// Send the request
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// Read the response
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Print the response
		fmt.Println(string(body))
	} else {
		req, _ := http.NewRequest("GET", "https://api.sandbox.paypal.com/v1/payments/payment", nil)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+accessToken.Token)
		res, _ := client.Do(req)
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		fmt.Println(string(body))
	}
}
*/

/*
func receivePayment(w http.ResponseWriter, r *http.Request) {

}

func makePayment(w http.ResponseWriter, r *http.Request) {

}

func authorizePayment(w http.ResponseWriter, r *http.Request) {

}

*/

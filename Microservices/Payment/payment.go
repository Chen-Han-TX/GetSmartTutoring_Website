package main

import (
	//"encoding/json"
	//"encoding/json"
	//"database/sql"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"log"
	//"net/http"
	"context"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/plutov/paypal/v4"
)

type Payment struct {
	PaymentID int `json:"payment_id"`
	Amount    int `json:"amount"`
	TutorID   int `json:"tutor_id"`
	StudentID int `json:"student_id"`
	SessionID int `json:"session_id"`
}

var (
	c           *paypal.Client
	err         error
	accessToken *paypal.TokenResponse
)

func main() {

	// Create a client instance
	clientID := "AYfInMQ1PMfzy31ogmES6Li1EZ7QBXDkWLeiVLzy4qMDNsWb2pFnsG1kRto_KkAyrg9a9gZQ40tQLe88"
	clientSecret := "EMTu8swH13MPTnYS-HUArKcWtHZ266nwJ98GgHFqp5pVClQVbTKkeb4hoJki7k3NcJkyzKCPfc-i8c9L"
	c, err = paypal.NewClient(clientID, clientSecret, paypal.APIBaseSandBox)
	if err != nil {
		fmt.Println(err.Error())
	}

	accessToken, err = c.GetAccessToken(context.Background())
	if err != nil {
		fmt.Println(err.Error())
	}
	router := mux.NewRouter()
	//router.HandleFunc("/api/payment/", SignUp).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/payment/makepayment/", createPayment).Methods("POST")

	fmt.Println("Listening at port 5053")
	log.Fatal(http.ListenAndServe(":5053", router))
}

// Add new trip record to database
func createPayment(w http.ResponseWriter, r *http.Request) {

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

	// Create a new http client
	client := &http.Client{}

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
}

func receivePayment(w http.ResponseWriter, r *http.Request) {

}

func makePayment(w http.ResponseWriter, r *http.Request) {

}

func authorizePayment(w http.ResponseWriter, r *http.Request) {

}

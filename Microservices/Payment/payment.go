package main

import (
	//"encoding/json"
	//"encoding/json"
	//"database/sql"
	"fmt"
	//"io/ioutil"
	//"log"
	//"net/http"
	"context"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
	//router := mux.NewRouter()
}

// Add new trip record to database

func createPayment(p Payment) {
	//paymentResponse, err := c.CreateOrder()

func createPayment(w http.ResponseWriter, r *http.Request) {
	//paymentResponse, err := c.CreateOrder()
	amount := paypal.Amount{
		Total:    "",
		Currency: "SGD",
	}
	//URL to redirect to after PayPal has complete the online payment
	redirectURI := "http://example.com/redirect-uri"
	//URL to redirect to if user clicks cancel
	cancelURI := "http://example.com/cancel-uri"
	description := "Description for this payment"
	paymentResult, err := c.CreateDirectPaypalPayment(amount, redirectURI, cancelURI, description)


}

func makePayment(w http.ResponseWriter, r *http.Request) {

}

func authorizePayment(w http.ResponseWriter, r *http.Request) {

}

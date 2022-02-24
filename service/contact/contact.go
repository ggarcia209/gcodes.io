package main

// contact relays a contact request to the Admin email address provided in the SAM template.

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/ggarcia209/portfolio/service/util/awsops"
	"github.com/ggarcia209/portfolio/service/util/httpops"
	"github.com/ggarcia209/portfolio/service/util/sesops"
)

const route = "/contact/submit" // POST

const failMsg = "On no! Something went wrong! Please contact me on LinkedIn for now."
const successMsg = "Your message has been sent. I will follow up at the email provided as soon as possible."

type request struct {
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// RootHandler handles HTTP requests
func RootHandler(w http.ResponseWriter, r *http.Request) {
	// verify content-type
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		httpops.ErrResponse(w, "Content-Type is not application/json", failMsg, http.StatusUnsupportedMediaType)
		return
	}

	// decode JSON object from http request
	data := request{}
	var unmarshalErr *json.UnmarshalTypeError

	// get contact request info
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	err := decoder.Decode(&data)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			httpops.ErrResponse(w, "Bad Request: Wrong type provided for field "+unmarshalErr.Field, failMsg, http.StatusBadRequest)
		} else {
			httpops.ErrResponse(w, "Bad Request: "+err.Error(), failMsg, http.StatusBadRequest)
		}
		return
	}

	// send email
	svc := awsops.NewDefaultAWSSession()
	ses := sesops.InitClient(svc)
	target := sesops.AdminEmailAddress()
	src := sesops.ServiceEmailAddress()

	err = sesops.SendContactRequest(ses, src, target, data.Email, data.Subject, data.Message)
	if err != nil {
		log.Printf("contactSubmit failed: %v", err)
		resp := response{Success: false, Message: failMsg}
		httpops.ErrResponse(w, "Request failed", resp, http.StatusInternalServerError)
		return
	}

	resp := response{Success: true, Message: successMsg}
	httpops.ErrResponse(w, "Success", resp, http.StatusOK)
	return
}

func main() {
	httpops.RegisterRoutes(route, RootHandler)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

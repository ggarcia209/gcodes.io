package httpops

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ggarcia209/portfolio/service/util/htmlops"
)

// HttpResponse contains a status code, message, and body to return to the client
// (Content-Type: application/json)
type HttpResponse struct {
	StatusCode      int         `json:"status_code"`
	Message         string      `json:"message"`
	Body            interface{} `json:"body"`
	IsBase64Encoded bool        `json:"isbase64encoded"`
}

// ErrResponse writes an http response with the given message, body, and HTTP status code.
func ErrResponse(w http.ResponseWriter, message string, body interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := HttpResponse{
		StatusCode: httpStatusCode,
		Message:    message,
		Body:       body,
	}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

// ErrResponse writes an http response with the given message, body, and HTTP status code.
func TextHtmlResponse(w http.ResponseWriter, body string, httpStatusCode int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(httpStatusCode)
	w.Write([]byte(body))
}

// ErrResponse writes an http response with the given message, body, and HTTP status code.
func Base64Response(w http.ResponseWriter, message string, body interface{}, httpStatusCode int) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(httpStatusCode)
	resp := HttpResponse{
		StatusCode:      httpStatusCode,
		Message:         message,
		Body:            body,
		IsBase64Encoded: true,
	}
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

// RegisterRoutes registers the HTTP handler func of each route
func RegisterRoutes(route string, rootHandler func(http.ResponseWriter, *http.Request)) {
	http.Handle(route, h(rootHandler))
}

// RegisterRoutes registers the HTTP handler func of each route for serving HTML.
func RegisterRoutesHtml(route string, rootHandler func(http.ResponseWriter, *http.Request)) {
	http.Handle(route, hhtml(rootHandler))
}

// h wraps a http.HandlerFunc and adds common headers.
// Access-Control-Allow headers are required to return responses from API Gateway Lambda proxy functions.
func h(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

// hhtml wraps a http.HandlerFunc and adds a text/html Content-Type header.
// Access-Control-Allow headers are required to return responses from API Gateway Lambda proxy functions.
func hhtml(next http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Amz-Date,Authorization,X-Api-Key,X-Amz-Security-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(w, r)
	})
}

// GetQueryStringParams returns an HTTP request's query string parameters.
func GetQueryStringParams(r *http.Request) map[string]string {
	data := make(map[string]string) // key value pairs for params
	params := r.URL.Query()
	for k, v := range params {
		data[k] = v[0]
	}
	return data
}

// HtmlHandler handles HTTP GET requests to serve HTML pages from html files stored in the lambda CodeUri Zip.
func HtmlHandler(w http.ResponseWriter, r *http.Request, filepath string) {
	failMsg := "<h1>Oh no! Something went wrong!</h1>"
	// get item
	data, err := htmlops.GetLocalHtml(filepath)
	if err != nil {
		log.Printf("HtmlHandler failed: %v", err)
		TextHtmlResponse(w, failMsg, http.StatusInternalServerError)
		return
	}

	html, err := htmlops.CreateHtmlTemplate(string(data), nil)
	if err != nil {
		log.Printf("HtmlHandler failed: %v", err)
		TextHtmlResponse(w, failMsg, http.StatusInternalServerError)
		return
	}

	// return order to admin
	TextHtmlResponse(w, html, http.StatusOK)
	return
}

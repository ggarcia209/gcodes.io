package main

/* about serves the about page. */

import (
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/ggarcia209/portfolio/service/util/httpops"
)

const route = "/about" // GET
const path = "./about.html"

// RootHandler handles HTTP request
func RootHandler(w http.ResponseWriter, r *http.Request) {
	httpops.HtmlHandler(w, r, path)
}

func main() {
	httpops.RegisterRoutesHtml(route, RootHandler)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

package main

/* resources serves the resources page. */

import (
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/ggarcia209/portfolio/service/util/httpops"
)

const route = "/resources" // GET
const path = "./resources.html"

// RootHandler handles HTTP request
func RootHandler(w http.ResponseWriter, r *http.Request) {
	httpops.HtmlHandler(w, r, path)
}

func main() {
	httpops.RegisterRoutesHtml(route, RootHandler)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

package main

/* index serves the home page. */

import (
	"log"
	"net/http"

	"github.com/apex/gateway"
	"github.com/ggarcia209/portfolio/service/util/httpops"
)

const route = "/" // GET
const path = "./index.html"

// RootHandler handles HTTP request to the root '/'
func RootHandler(w http.ResponseWriter, r *http.Request) {
	httpops.HtmlHandler(w, r, path)
}

func main() {
	httpops.RegisterRoutesHtml(route, RootHandler)
	log.Fatal(gateway.ListenAndServe(":3000", nil))
}

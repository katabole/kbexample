package actions

import (
	"fmt"
	"net/http"

	"github.com/golang/gddo/httputil"
)

type AcceptType int

const (
	AcceptJSON AcceptType = iota
	AcceptHTML
	AcceptCSV
)

// AcceptContentType figures out which of the above AcceptTypes are indicated by the "Accept" header, defaulting to
// "text/html". When you have an action that you want to return different content based on this header, so:
//
//	switch AcceptContentType(c) {
//	case AcceptCSV:
//		// Render the response as CSV data
//	default:
//		// Render JSON or HTML or whatever your usual format is
//	}
func AcceptContentType(r *http.Request) AcceptType {
	val := httputil.NegotiateContentType(r, []string{"application/json", "text/html", "text/csv"}, "text/html")
	switch val {
	case "application/json":
		return AcceptJSON
	case "text/html":
		return AcceptHTML
	case "text/csv":
		return AcceptCSV
	default:
		panic(fmt.Sprintf("got negotiated Accept header that makes no sense: %s", val))
	}
}

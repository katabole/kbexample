package actions

import (
	"net/http"
	"reflect"

	// If Go adds accept header negotiation to the standard library we may want to use it, see
	// https://github.com/golang/go/issues/19307
	"github.com/elnormous/contenttype"
)

type ContentType int

const (
	ContentTypeHTML = iota
	ContentTypeJSON
)

var (
	mediaTypeHTML contenttype.MediaType = contenttype.NewMediaType("text/html")
	mediaTypeJSON                       = contenttype.NewMediaType("application/json")

	// NOTE: the first in this list will be used as the default when negotiating Accept headers below.
	availableMediaTypes []contenttype.MediaType = []contenttype.MediaType{mediaTypeHTML, mediaTypeJSON}
)

// GetContentType figures out which of the above supported content/media types we should use with our request, by
// checking first the "Content-Type" header, or falling back to the "Accept" header. Use it like so:
//
//	switch GetContentType(r) {
//	case ContentTypeHTML:
//		// Render the response as HTML data
//	case ContentTypeJSON:
//		// Render JSON
//	}
func GetContentType(r *http.Request) ContentType {
	mediaType, err := contenttype.GetMediaType(r)
	if err == nil && !reflect.ValueOf(mediaType).IsZero() {
		for _, mt := range availableMediaTypes {
			if mt.Type == mediaType.Type && mt.Subtype == mediaType.Subtype {
				return mediaTypeToContentType(mt)
			}
		}
	}
	// We got an error or there is not content-type header. Try using Accept instead.

	accepted, _, err := contenttype.GetAcceptableMediaType(r, availableMediaTypes)
	if err != nil {
		// No accept header either, just go with default
		return ContentTypeHTML
	}
	return mediaTypeToContentType(accepted)
}

func mediaTypeToContentType(mt contenttype.MediaType) ContentType {
	if mt.Type == mediaTypeHTML.Type && mt.Subtype == mediaTypeHTML.Subtype {
		return ContentTypeHTML
	} else if mt.Type == mediaTypeJSON.Type && mt.Subtype == mediaTypeJSON.Subtype {
		return ContentTypeJSON
	}
	panic("unknown media type: " + mt.String())
}

package crawler

import (
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PageData struct {
	URL             string
	FinalURL        string
	StatusCode      int
	Headers         http.Header
	Document        *goquery.Document
	HTML            string
	BodyText        string
	Depth           int
	ResponseTime    time.Duration
	ContentLength   int64
	ContentType     string
	RedirectChain   []string
	InternalLinks   []LinkData
	ExternalLinks   []LinkData
	Images          []ImageData
	Scripts         []ResourceData
	Stylesheets     []ResourceData
}

type LinkData struct {
	URL        string
	AnchorText string
	Rel        string
	Target     string
	IsInternal bool
	StatusCode int
}

type ImageData struct {
	URL    string
	Alt    string
	Width  string
	Height string
	Loading string
	Size   int64
}

type ResourceData struct {
	URL      string
	Size     int64
	Type     string
	IsAsync  bool
	IsDefer  bool
	Location string
}

package request

import (
	"net/http"
)

// Params paramter lsit
type Params map[string]string

//UploadFile contains all info of upload 
type UploadFile struct {
	Path, Name, Mime string 
	Size int
	Buffer []byte
}

//Uploads is a list of UploadFile
type Uploads map[string]UploadFile

//Request extents the HTTP requst with RouteParams & Files
type Request struct {
	*http.Request
	RouteParams Params
	Files Uploads
}

// Redirect redirects the client
func (r *Request) Redirect(w http.ResponseWriter, path string) {
	req := r.Request
	http.Redirect(w, req, path, http.StatusMovedPermanently)
}

// NewRequest returns a new simple-router request
func NewRequest(r *http.Request) *Request {
	return &Request{r, nil, nil}
}


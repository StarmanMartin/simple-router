package router

import (
	"net/http"
)

//UploadFile contains all info of upload 
type UploadFile struct {
	Path, Name, Mime string 
	Size int
	Buffer []byte
}

type uploads map[string]UploadFile

//Request extents the HTTP requst with RouteParams & Files
type Request struct {
	*http.Request
	RouteParams params
	Files uploads
}

// Redirect redirects the client
func (r *Request) Redirect(w http.ResponseWriter, path string) {
	req := r.Request
	http.Redirect(w, req, path, http.StatusMovedPermanently)
}

func newRequest(r *http.Request) *Request {
	return &Request{r, nil, nil}
}


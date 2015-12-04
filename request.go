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

func newRequest(r *http.Request) *Request {
	return &Request{r, nil, nil}
}


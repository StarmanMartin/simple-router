package rtest

import (
    "os"
    "bytes"
    "io"
	"mime/multipart"    
    "net/http"
    "net/url"
    "path/filepath"
    "regexp"
    "fmt"
    "strings"
    "github.com/starmanmartin/simple-router/request"
)

// NewFileUploadRequest creates a new file upload http request with optional extra params
func NewFileUploadRequest(uri string, params map[string]string, paramName, path string) (*request.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	hRequest, _ := http.NewRequest("POST", uri, body)
	hRequest.Header.Add("Content-type", "multipart/form-data;boundary="+writer.Boundary()+";")
	return request.NewRequest(hRequest), nil
}

// NewPostReqeust generates a POST request. It is simpel tu use ist as a test request.
func NewPostReqeust(urlRequest string, vals url.Values) (*request.Request, error) {
    var valsBuffer *bytes.Buffer
    if vals != nil {
       vals = url.Values{}
    }
    
     valsBuffer = bytes.NewBufferString(vals.Encode())
    
    req, err := http.NewRequest("Post", urlRequest, valsBuffer)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    return request.NewRequest(req), nil
}

// NewXHRPostRequest generates a XHR POST request. It is simpel tu use ist as a test request.
func NewXHRPostRequest(urlRequest string, vals url.Values) (*request.Request, error){
    req, err := NewPostReqeust(urlRequest, vals)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    
    return req, nil
}


// NewGetRequest generates a GET request. It is simpel tu use ist as a test request.
func NewGetRequest(urlRequest string, vals map[string]string) (*request.Request, error) {
    hasQuestoinReg := regexp.MustCompile(`\?`) 
    if len(vals) > 0 && !hasQuestoinReg.MatchString(urlRequest) {
        valsAsArray := make([]string, 0, len(vals))
        for key, vals := range vals {
            valsAsArray = append(valsAsArray, fmt.Sprintf("%s=%s", key, vals))
        }
        
        urlRequest = fmt.Sprintf("%s?%s", urlRequest, strings.Join(valsAsArray, "&"))
    }   
    
    req, err := http.NewRequest("Get", urlRequest, nil)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    
    return request.NewRequest(req), nil
}

// NewXHRGetRequest generates a XHR GET request. It is simpel tu use ist as a test request.
func NewXHRGetRequest(urlRequest string, vals map[string]string) (*request.Request, error){
    req, err := NewGetRequest(urlRequest, vals)
    if err != nil {
        return nil, err
    }
    
    req.Header.Set("X-Requested-With", "XMLHttpRequest")
    
    return req, nil
}
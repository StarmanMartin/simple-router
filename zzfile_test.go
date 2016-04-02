package router

import (
	"strings"
	"testing"
    "net/http/httptest"
	"github.com/starmanmartin/simple-fs"
	"github.com/starmanmartin/simple-router/rtest"
)

func TestUploadAndPublic(t *testing.T) {
	router := NewRouter()

	handler := router.UploadPath("/test", true)
	params := map[string]string{
		"name":    "Maritn",
		"surname": "bar",
	}

	r, err := rtest.NewFileUploadRequest("http://localhost:80", params, "image", "test/image.png")

	if err != nil {
		t.Error("Error", err)
		t.SkipNow()
	}

	if isNext, err := handler(nil, r); !isNext {
		t.Error("No next on upload")
	} else if err != nil {
		t.Error("Error", err)
	}
    
    imageName :=  r.Files["image"].Name
    imagePath := r.Files["image"].Path

	if ex, err := fs.Exists(imagePath + "/" + imageName); !ex {
		t.Error("File not exists")
	} else if err != nil {
		t.Error("Error", err)
	}

	r, err = rtest.NewFileUploadRequest("http://localhost:80", params, "image", "test/image.png")

	if err != nil {
		t.Error("Error", err)
		t.SkipNow()
	}

	r.Header.Set("Content-type", "multipart/form-data;")

	if isNext, err := handler(nil, r); !isNext {
		t.Error("No next on upload")
	} else if err == nil {
		t.Error("No error")
	}

	r, err = rtest.NewPostReqeust("http://localhost:80", nil)

	if err != nil {
		t.Error("Error", err)
		t.SkipNow()
	}

	r.Header.Set("Content-type", "multipart/form-data;boundary=foo;")

	if isNext, err := handler(nil, r); !isNext {
		t.Error("No next on upload")
	} else if err == nil {
		t.Error("No error")
	}

	router.Public("test")
    

	if finalHandler, ok := findListOfHandler(routerList[mGet], strings.Split("test/" + imageName, "/"), false); ok {
		for _, finalHanderTemp := range finalHandler {
			for _, handerTemp := range finalHanderTemp.hanlder {
                w := httptest.NewRecorder()
                req, _ := rtest.NewGetRequest("http://localhost.com/test/" + imageName, nil)            
				handerTemp(w, req)
                if w.Header().Get("Content-Type") != "image/png" {
                    t.Error("Public download wrong Content-Type")
                }
                
                w = httptest.NewRecorder()
                req, _ = rtest.NewGetRequest("http://localhost.com/test/not.png", nil)              
                handerTemp(w, req)
                if w.Code != 404 {
                    t.Error("Not not found", w.Code)
                }
			}
		}

	}
}

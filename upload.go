package router

import (
	"github.com/google/uuid"
	"github.com/starmanmartin/simple-router/request"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var cleanUpSync sync.Once
var uploadPath *string
var sourceinfo os.FileInfo

func initUpload() {
	cleanUpSync = sync.Once{}
}

func newUploadPaser(path string, isBuffer bool) HTTPHandler {
	cwd, err := filepath.Abs(filepath.Dir(os.Args[0]))
	sourceinfo, err = os.Stat(cwd)
	if err != nil {
		log.Fatal(err)
	}

	path = filepath.Join(cwd, path)
	uploadPath = &path

	cleanUpSync.Do(func() {
		cleanUpUpload()
	})

	return func(w http.ResponseWriter, r *request.Request) (isNext bool, err error) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		isNext = true
		if err != nil {
			return
		}

		r.PostForm = make(map[string][]string)
		r.Files = make(request.Uploads)
		if strings.HasPrefix(mediaType, "multipart/") {
			mr, err := r.MultipartReader()
			if err == io.EOF {
				return isNext, nil
			} else if err != nil {
				return isNext, err
			}
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					return isNext, nil
				} else if err != nil {
					return isNext, err
				}

				slurp, err := ioutil.ReadAll(p)

				if err != nil {
					return isNext, err
				}

				if len(p.Header["Content-Type"]) > 0 {
					uuidS, _ := uuid.NewUUID()
					filename := uuidS.String() + p.FileName()
					err = ioutil.WriteFile(filepath.Join(path, filename), slurp, sourceinfo.Mode())
					if err != nil {
						return isNext, err
					}

					fileElement := request.UploadFile{path, filename, p.Header["Content-Type"][0], len(slurp), nil}
					if isBuffer {
						fileElement.Buffer = slurp
					}

					r.Files[p.FormName()] = fileElement
				} else {
					r.PostForm[p.FormName()] = []string{string(slurp)}
				}
			}
		}

		return isNext, nil
	}
}

func cleanUpUpload() {
	err := os.MkdirAll(*uploadPath, sourceinfo.Mode())
	if err != nil {
		log.Fatal(err)
	}
	go func(path *string) {

	}(uploadPath)
}

package router

import (
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"github.com/wayn3h0/go-uuid"
)

var cleanUpSync sync.Once
var uploadPath *string
var sourceinfo os.FileInfo

func init() {
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

	return func(w http.ResponseWriter, r *Request) {
		mediaType, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if err != nil {
			log.Fatal(err)
		}

		r.PostForm = make(map[string][]string)
		r.Files = make(uploads)
		if strings.HasPrefix(mediaType, "multipart/") {
			mr, _ := r.MultipartReader()
			for {
				p, err := mr.NextPart()
				if err == io.EOF {
					return
				}

				if err != nil {
					log.Fatal(err)
				}

				slurp, err := ioutil.ReadAll(p)

				if err != nil {
					log.Fatal(err)
				}

				if len(p.Header["Content-Type"]) > 0 {
					uuidS, _ := uuid.NewV1()
					filename := uuidS.String() + p.FileName()
					err = ioutil.WriteFile(filepath.Join(path, filename), slurp, sourceinfo.Mode())
					if err != nil {
						log.Fatal(err)
					}

					fileElement := UploadFile{path, filename, p.Header["Content-Type"][0], len(slurp), nil}
					if isBuffer {
						fileElement.Buffer = slurp
					}

					r.Files[p.FormName()] = fileElement
				} else {
					r.PostForm[p.FormName()] = []string{string(slurp)}
				}
			}
		}
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
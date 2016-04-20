package view

import (
	"github.com/starmanmartin/simple-fs"
	"html/template"
	"path/filepath"
	"os"
	"regexp"
)

var (
	//ViewPath is the root path of the views
	ViewPath, cwd string
	extendedReg = regexp.MustCompile(`<!--\s*extent:\s*(.*)\s*-->`)
	listReg = regexp.MustCompile(`\s*,\s*`)
	htmlExtReg = regexp.MustCompile(`html$`)
    templateList = make(map[string]*template.Template)
)

func GetTemplate(name string) (tmp *template.Template, ok bool) {
    tmp, ok =  templateList[name]
    return
}

//ParseTemplate parses template
func ParseTemplate(name, filePath string) (tmp *template.Template){
    return ParseTemplateFunc(name, filePath, nil)
}

//ParseTemplateFunc parses template with parse functions
func ParseTemplateFunc(name, filePath string, funcMap map[string]interface{}) (tmp *template.Template){
	cwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	lineAsList, _ := fs.ReadLines(filepath.Join(cwd, ViewPath, filePath), -1)
	parentList := extendedReg.FindAllStringSubmatch(lineAsList[0], -1)[0]
	
	parentList = listReg.Split(parentList[1], -1)
	parentList = append(parentList, filePath)
	for i, pathEnd := range parentList {
		if !htmlExtReg.MatchString(pathEnd) {
			pathEnd = pathEnd + ".html"
		}
		
		parentList[i] = filepath.Join(cwd, ViewPath, pathEnd)
	}
    
    if funcMap == nil {
        funcMap = make(map[string]interface{})
    }
    
    tmp = template.Must(template.New(name).Funcs(funcMap).ParseFiles(parentList...))
	
    templateList[name] = tmp
	
	return
}
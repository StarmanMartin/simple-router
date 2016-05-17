package view

import (
	"html/template"
	"os"
	"path/filepath"
	"regexp"

	"github.com/starmanmartin/simple-fs"
)

var (
	//ViewPath is the root path of the views
	ViewPath, cwd string
	extendedReg   = regexp.MustCompile(`<!--\s*extent:\s*(.*)\s*-->`)
	listReg       = regexp.MustCompile(`\s*,\s*`)
	htmlExtReg    = regexp.MustCompile(`html$`)
	templateList  = make(map[string]*template.Template)
	baseFunctions = make(map[string]interface{})
)

//GetTemplate retruns a template by a given name
func GetTemplate(name string) (tmp *template.Template, ok bool) {
	tmp, ok = templateList[name]
	return
}

func setBaseFunctions(base map[string]interface{}) {
	if base == nil {
		baseFunctions = make(map[string]interface{})
	} else {
		baseFunctions = base
	}
}

//ParseTemplate parses template
func ParseTemplate(name, filePath string) (tmp *template.Template) {
	return ParseTemplateFunc(name, filePath, nil)
}

//ParseTemplateFunc parses template with parse functions
func ParseTemplateFunc(name, filePath string, funcMap map[string]interface{}) (tmp *template.Template) {
	if funcMap != nil {
		funcMap = joinFunctionMaps(baseFunctions, funcMap)
	} else {
		funcMap = baseFunctions
	}

	cwd, _ = filepath.Abs(filepath.Dir(os.Args[0]))
	lineAsList, _ := fs.ReadLines(filepath.Join(cwd, ViewPath, filePath), -1)
	resultList := extendedReg.FindAllStringSubmatch(lineAsList[0], -1)
	var parentList []string
	if len(resultList) > 0 {
		parentList = resultList[0]
		parentList = listReg.Split(parentList[1], -1)
		parentList = append(parentList, filePath)
	} else {
		parentList = []string{filePath}
	}

	for i, pathEnd := range parentList {
		if !htmlExtReg.MatchString(pathEnd) {
			pathEnd = pathEnd + ".html"
		}

		parentList[i] = filepath.Join(cwd, ViewPath, pathEnd)
	}

	tmp = template.Must(template.New(name).Funcs(funcMap).ParseFiles(parentList...))

	templateList[name] = tmp

	return
}

func joinFunctionMaps(maps ...map[string]interface{}) (joined map[string]interface{}) {
	joined = make(map[string]interface{})
	for _, mapIn := range maps {
		for k, v := range mapIn {
			joined[k] = v
		}
	}

	return
}

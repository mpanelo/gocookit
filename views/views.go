package views

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
)

const (
	layoutsDir = "views/layouts"
	fileExt    = ".html"
	contentDir = "views/content"
)

type View struct {
	template *template.Template
}

func NewView(content string) *View {
	files := getLayouts()
	files = append(files, fmt.Sprintf("%s/%s%s", contentDir, content, fileExt))

	t, err := template.ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		template: t,
	}
}

func getLayouts() []string {
	files, err := filepath.Glob(fmt.Sprintf("%s/*%s", layoutsDir, fileExt))
	if err != nil {
		panic(err)
	}
	return files
}

func (v *View) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	v.Render(rw, nil)
}

func (v *View) Render(rw http.ResponseWriter, data interface{}) {
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}

	err := v.template.ExecuteTemplate(rw, "bootstrap", vd)
	if err != nil {
		panic(err)
	}
}

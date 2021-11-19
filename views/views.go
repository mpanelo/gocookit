package views

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/mpanelo/gocookit/context"
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
	v.Render(rw, r, nil)
}

func (v *View) Render(rw http.ResponseWriter, r *http.Request, data interface{}) {
	var vd Data
	var buf bytes.Buffer

	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}

	vd.User = context.User(r.Context())

	err := v.template.ExecuteTemplate(&buf, "bootstrap", vd)
	if err != nil {
		http.Error(rw, AlertGenericMsg, http.StatusInternalServerError)
		return
	}

	io.Copy(rw, &buf)
}

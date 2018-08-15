package views

import (
	"bytes"
	"errors"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"

	"github.com/tjcain/theFieldBiologist/context"
)

const (
	// LayoutDir represents the path to template layouts used during globbing.
	LayoutDir string = "views/layouts/"
	// TemplateExt is the extension of template files used during globbing.
	TemplateExt string = ".gohtml"
	// TemplateDir is the root directory for all template files.
	TemplateDir string = "views/"
)

// View ...comment
type View struct {
	Template *template.Template
	Layout   string
}

// NewView makes it easier to create views. It will append common template files
// to the list of files appended, parse those template files and return a new
// *View.
func NewView(layout string, files ...string) *View {
	addTemplatePath(files)
	addTemplateExt(files)
	files = append(files, layoutFiles()...)
	// t, err := template.ParseFiles(files...)
	t, err := template.New("").Funcs(template.FuncMap{
		"csrfField": func() (template.HTML, error) {
			return "", errors.New("csrf is not implimented")
		},
	}).ParseFiles(files...)
	if err != nil {
		panic(err)
	}

	return &View{
		Template: t,
		Layout:   layout,
	}
}

// ServeHTTP impliments the Handler interface
func (v *View) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	v.Render(w, r, nil)
}

// Render is responsible for rendering the view called by the HandlerFuncs
func (v *View) Render(w http.ResponseWriter, r *http.Request, data interface{}) {
	w.Header().Set("Content-Type", "text/html")
	var vd Data
	switch d := data.(type) {
	case Data:
		vd = d
	default:
		vd = Data{
			Yield: data,
		}
	}
	vd.User = context.User(r.Context())
	var buf bytes.Buffer
	csrfField := csrf.TemplateField(r)
	tpl := v.Template.Funcs(template.FuncMap{
		"csrfField": func() template.HTML {
			return csrfField
		},
	})
	err := tpl.ExecuteTemplate(&buf, v.Layout, vd)
	if err != nil {
		log.Println("RENDER:", err)
		http.Error(w, "Oops! Something went wrong. If the problem persists "+
			"contact us.", http.StatusInternalServerError)
		return
	}
	io.Copy(w, &buf)
}

// other functions:
func layoutFiles() []string {
	files, err := filepath.Glob(LayoutDir + "*" + TemplateExt)
	if err != nil {
		// panic here: if this function errors the application cannot start.
		panic(err)
	}
	return files
}

// addTemplatePath prepends the TemplateDir to each string in the provided slice
func addTemplatePath(files []string) {
	for i, f := range files {
		files[i] = TemplateDir + f
	}
}

// addTemplateExt appends the TemplateExt to each string in the provided slice
func addTemplateExt(files []string) {
	for i, f := range files {
		files[i] = f + TemplateExt
	}
}

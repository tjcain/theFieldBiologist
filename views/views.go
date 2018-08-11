package views

import (
	"html/template"
	"net/http"
	"path/filepath"
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
	t, err := template.ParseFiles(files...)
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
	if err := v.Render(w, nil); err != nil {
		panic(err)
	}
}

// Render is responsible for rendering the view called by the HandlerFuncs
func (v *View) Render(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "text/html")
	switch data.(type) {
	case Data:
		// do nothing
	default:
		data = Data{
			Yield: data,
		}
	}
	return v.Template.ExecuteTemplate(w, v.Layout, data)
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

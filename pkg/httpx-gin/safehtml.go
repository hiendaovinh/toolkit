package httpx

import (
	"net/http"

	"github.com/gin-gonic/gin/render"
	"github.com/google/safehtml/template"
)

var htmlContentType = []string{"text/html; charset=utf-8"}

type SafeHTMLProduction struct {
	Template *template.Template
	Delims   render.Delims
}

func (r SafeHTMLProduction) Instance(name string, data any) render.Render {
	return SafeHTML{
		Template: r.Template,
		Name:     name,
		Data:     data,
	}
}

type SafeHTML struct {
	Template *template.Template
	Name     string
	Data     any
}

func (r SafeHTML) Render(w http.ResponseWriter) error {
	r.WriteContentType(w)

	if r.Name == "" {
		return r.Template.Execute(w, r.Data)
	}
	return r.Template.ExecuteTemplate(w, r.Name, r.Data)
}

// WriteContentType (HTML) writes HTML ContentType.
func (r SafeHTML) WriteContentType(w http.ResponseWriter) {
	writeContentType(w, htmlContentType)
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

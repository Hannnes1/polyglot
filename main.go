package main

import (
	"html/template"
	"io"
	"net/http"
)

var t *Templates

type Templates struct {
	templates *template.Template
}

type Translation struct {
	ID    int
	Key   string
	Value string
}

func newTemplates() *Templates {
	return &Templates{
		templates: template.Must(template.ParseGlob("templates/*.html")),
	}
}

func handleGetTranslations(w http.ResponseWriter, r *http.Request) {
	data := Translation{ID: 1, Key: "hello", Value: "Hello"}

	t.render(w, "translations.html", data)
}

func (t *Templates) render(w io.Writer, name string, data interface{}) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	t = newTemplates()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	http.HandleFunc("/", handleGetTranslations)

	http.ListenAndServe(":8080", nil)
}

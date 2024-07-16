package main

import (
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"net/http"
)

type Handler struct {
	templates *template.Template
}

type Languages struct {
	En interface{}
	Sv interface{}
}

type Language struct {
	Key  string
	Name string
	Data map[string]interface{}
}

func isMap(v interface{}) bool {
	_, ok := v.(map[string]interface{})
	return ok
}

// Helper for passing multiple parameters to template
func dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}
	dict := make(map[string]interface{}, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]
	}
	return dict, nil
}

func newHandler() *Handler {
	t := template.New("").Funcs(template.FuncMap{
		"IsMap": isMap,
		"Dict":  dict,
	})

	t = template.Must(t.ParseGlob("templates/components/*.html"))
	t = template.Must(t.ParseGlob("templates/*.html"))

	return &Handler{
		templates: t,
	}
}

func (h *Handler) render(w io.Writer, name string, data interface{}) error {
	return h.templates.ExecuteTemplate(w, name, data)
}

func (h *Handler) handleTranslations(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.handleGetTranslations(w, r)
		return
	} else if r.Method == http.MethodPost {
		h.handlePostTranslations(w, r)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *Handler) handleGetTranslations(w http.ResponseWriter, r *http.Request) {
	jsonDataEn := `{
        "confirm": "Confirm",
        "cancel": "Cancel",
        "example": {
            "title": "Example",
            "content": "This is an example content.",
			"subExample": {
				"title": "Subexample"
			}
        }
    }`

	jsonDataSv := `{
		"confirm": "Bekräfta",
		"cancel": "Avbryt",
		"example": {
			"title": "Exempel",
			"content": "Detta är ett exempel innehåll.",
			"subExample": {
				"title": "Underexempel"
			}
		}
	}`

	var dataEn map[string]interface{}
	var dataSv map[string]interface{}

	errEn := json.Unmarshal([]byte(jsonDataEn), &dataEn)
	if errEn != nil {
		http.Error(w, errEn.Error(), http.StatusInternalServerError)
		return
	}
	errSv := json.Unmarshal([]byte(jsonDataSv), &dataSv)
	if errSv != nil {
		http.Error(w, errSv.Error(), http.StatusInternalServerError)
		return
	}

	l := Languages{
		En: Language{Key: "en", Name: "English", Data: dataEn},
		Sv: Language{Key: "sv", Name: "Swedish", Data: dataSv},
	}

	h.render(w, "translations.html", l)
}

// Save incoming data.
// Data comes in as form data, where nested keys are separated by dashes.
func (h *Handler) handlePostTranslations(w http.ResponseWriter, r *http.Request) {
	// TODO: Save incoming data
}

func main() {
	renderer := newHandler()

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))

	http.HandleFunc("/", renderer.handleTranslations)

	http.ListenAndServe(":8080", nil)
}

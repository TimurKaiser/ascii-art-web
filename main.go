package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Page struct {
	Valeur string
}

var templates *template.Template

func init() {
	templates = template.Must(template.ParseFiles("index.html", "404.html", "500.html", "400.html"))
}

func main() {
	http.HandleFunc("/execute", executeHandler)
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/404", notFoundHandler)
	http.HandleFunc("/500", internalServerErrorHandler)
	http.HandleFunc("/400", badRequestHandler)

	http.ListenAndServe(":8080", nil)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	_, err := filepath.Abs("404.html")
	if err != nil {
		http.Error(w, "Erreur lors de la résolution du chemin du fichier 404.html", http.StatusNotFound)
		return
	}

	var p Page
	p.Valeur = "404ERROR"

	err = templates.ExecuteTemplate(w, "404.html", p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du template 404 : %s", err), http.StatusInternalServerError)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/404", http.StatusSeeOther)
		return
	}

	var p Page
	p.Valeur = "Ascii Art Web"

	err := templates.ExecuteTemplate(w, "index.html", p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du template index : %s", err), http.StatusInternalServerError)
	}
}

func internalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
	var p Page
	p.Valeur = "500ERROR"

	err := templates.ExecuteTemplate(w, "500.html", p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du template 500 : %s", err), http.StatusInternalServerError)
	}
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	var p Page
	p.Valeur = "400 Bad Request"
	err := templates.ExecuteTemplate(w, "400.html", p)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erreur lors de l'exécution du template 400 : %s", err), http.StatusInternalServerError)
	}
}

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	text := r.FormValue("text")
	action := r.FormValue("action")

	if !isASCII(text) {
		http.Redirect(w, r, "/400", http.StatusBadRequest)
		return
	}

	fileName := getFileName(action)
	if fileName == "" {
		http.Redirect(w, r, "/500", http.StatusInternalServerError)
		return
	}

	var result string
	for _, line := range strings.Split(text, "\\n") {
		lineResult, err := convertInputWithFile(line, fileName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result += lineResult
	}

	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, result)
}

func getFileName(action string) string {
	switch action {
	case "Standard":
		return "standard.txt"
	case "Shadow":
		return "shadow.txt"
	case "Thinkertoy":
		return "thinkertoy.txt"
	default:
		return ""
	}
}

func convertInputWithFile(line string, fileName string) (string, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")

	var result string
	if line != "" {
		for i := 1; i < 9; i++ {
			for _, char := range line {
				ascii := int(char) - 32
				result += lines[ascii*9+i]
			}
			result += "\n"
		}
	} else {
		result += "\n"
	}
	return result, nil
}

func isASCII(s string) bool {
	for _, char := range s {
		if char >= 127 {
			return false
		}
	}
	return true
}

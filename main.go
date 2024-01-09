package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
)

type Page struct {
	Valeur string
}

var templates *template.Template

func main() {
	templates = template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/execute", executeHandler)

	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	var p Page
	p.Valeur = "Ascii Art Web"
	templates.Execute(w, p)
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

func executeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		text := r.FormValue("text")
		action := r.FormValue("action")

		args := strings.Split(text, "\\n")

		var fileName string
		switch action {
		case "Standard":
			fileName = "standard.txt"
		case "Shadow":
			fileName = "shadow.txt"
		case "Thinkertoy":
			fileName = "thinkertoy.txt"
		default:
			http.Error(w, "Action non prise en charge", http.StatusBadRequest)
			return
		}

		var result string
		for _, line := range args {
			lineResult, err := convertInputWithFile(line, fileName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			result += lineResult
		}

		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, result)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

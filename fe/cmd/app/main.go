package main

import (
	"fe/internal/config"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const (
	TEMPLATES_PATH = "./assets/templates"
)

func main() {
	config := config.NewConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, TEMPLATES_PATH+"/test.page.gohtml")
	})

	fmt.Println("Starting front end service on port ", config.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", config.Port), nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(w http.ResponseWriter, t string) {
	partials := []string{
		TEMPLATES_PATH + "/base.layout.gohtml",
		TEMPLATES_PATH + "/header.partial.gohtml",
		TEMPLATES_PATH + "/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, t)
	templateSlice = append(templateSlice, partials...)

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

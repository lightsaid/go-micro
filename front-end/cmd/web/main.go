package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var webPort = ":9000"

func main(){

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./cmd/web/templates/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		render(w, "test.page.tmpl")
	})


	fmt.Println("Starting front end server on port ", webPort)
	if err := http.ListenAndServe(webPort, nil); err != nil {
		log.Panic(err.Error())
	}
}

func render(w http.ResponseWriter, t string) {
	// 基础模板
	partials := []string{
		"./cmd/web/templates/base.layout.tmpl",
		"./cmd/web/templates/header.partial.tmpl",
		"./cmd/web/templates/footer.partial.tmpl",
	}

	// 每一个page都使用基础模板
	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))
	templateSlice = append(templateSlice, partials...)

	// 解释页面模板
	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

}
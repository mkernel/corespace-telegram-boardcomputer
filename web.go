package main

import (
	"bytes"
	"html/template"
	"net/http"
)

//idea: embed the templates and the static html files using https://github.com/mjibson/esc/
var templates *template.Template

func setupHTTP() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("web/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", httpHandler)

	templates, _ = template.New("templates").ParseFiles("templates/page.html", "templates/ships.html", "templates/crew.html", "templates/inventory.html", "templates/transactions.html")

	go http.ListenAndServe(":8181", nil)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		w.WriteHeader(200)
		var tmplbuf bytes.Buffer
		var crews []crew
		database.Preload("Members").Preload("Items").Preload("Contacts").Preload("Transactions").Find(&crews)
		templates.ExecuteTemplate(&tmplbuf, "ships", crews)
		templates.ExecuteTemplate(w, "page", template.HTML(tmplbuf.String()))
	} else {
		w.WriteHeader(404)
	}
}

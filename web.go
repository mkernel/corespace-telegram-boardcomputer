package main

import (
	"bytes"
	"html/template"
	"net/http"

	"github.com/mkernel/corespace-telegram-boardcomputer/web"
	tmpl "github.com/mkernel/corespace-telegram-boardcomputer/web_templates"
)

//This is made to work using https://github.com/mjibson/esc/

//go:generate esc -o web/encoded.go -pkg web web
//go:generate esc -o web_templates/encoded.go -pkg web_templates web_templates

var templates *template.Template

func transactionSum(txs []transaction) float64 {
	var sum float64 = 0
	for _, tx := range txs {
		sum += tx.Value
	}
	return sum
}

func setupHTTP() {
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(web.Dir(false, "/web/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets/"))))
	http.HandleFunc("/", httpHandler)

	var fmap = template.FuncMap{
		"sum": transactionSum,
	}

	templates = template.New("page")
	templates.Funcs(fmap)
	templates.Parse(tmpl.FSMustString(false, "/web_templates/page.html"))
	templates.Parse(tmpl.FSMustString(false, "/web_templates/ships.html"))
	templates.Parse(tmpl.FSMustString(false, "/web_templates/crew.html"))
	templates.Parse(tmpl.FSMustString(false, "/web_templates/inventory.html"))
	templates.Parse(tmpl.FSMustString(false, "/web_templates/transactions.html"))
	templates.Parse(tmpl.FSMustString(false, "/web_templates/contacts.html"))
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

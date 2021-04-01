package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	tmplfs "github.com/mkernel/corespace-telegram-boardcomputer/admin_templates"
)

var admin_templates *template.Template

func SetupAdmin() {
	http.HandleFunc("/admin", AdminIndex)
	http.HandleFunc("/admin/status", AdminStatus)
	http.HandleFunc("/admin/transactions", AdminTransactions)
	http.HandleFunc("/admin/inventory", AdminInventory)
	http.HandleFunc("/admin/inventory/delete", AdminDeleteItem)
	http.HandleFunc("/admin/inventory/edit", AdminEditItem)
	http.HandleFunc("/admin/crew", AdminMembers)
	http.HandleFunc("/admin/crew/delete", AdminMembersDelete)
	http.HandleFunc("/admin/chat", AdminChat)
	http.HandleFunc("/admin/contacts", AdminContacts)
	http.HandleFunc("/admin/contacts/delete", AdminContactsDelete)
	http.HandleFunc("/admin/contacts/chat", AdminContactsChat)

	admin_templates = template.New("page")
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/page.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/crew-list.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/crew-status.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/tx-list.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/inventory-list.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/members-list.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/chats.html"))
	admin_templates.Parse(tmplfs.FSMustString(false, "/admin_templates/contacts-list.html"))
}

func AdminPage(payload bytes.Buffer, w http.ResponseWriter) {
	admin_templates.ExecuteTemplate(w, "page", template.HTML(payload.String()))
}

func AdminRender(template string, data interface{}) bytes.Buffer {
	var buf bytes.Buffer
	admin_templates.ExecuteTemplate(&buf, template, data)
	return buf
}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	var datasets []crew
	database.Find(&datasets)
	AdminPage(AdminRender("crew-list", datasets), w)
}

func AdminStatus(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var dataset crew
	database.First(&dataset, id)
	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		status := r.PostFormValue("status")
		dataset.Description = status
		database.Save(&dataset)
		http.Redirect(w, r, "/admin", 302)
		return
	}
	w.WriteHeader(200)
	AdminPage(AdminRender("crew-status", dataset), w)
}

func AdminTransactions(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var transactions []transaction

	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		desc := r.PostFormValue("Description")
		value, _ := strconv.Atoi(r.PostFormValue("Value"))
		tx := transaction{Date: int(time.Now().Unix()), CrewID: uint(id), Value: float64(value), Description: desc}
		database.Create(&tx)
	}

	filter := transaction{CrewID: uint(id)}
	database.Where(&filter).Order("date asc").Find(&transactions)
	var datasets []struct {
		Tx             transaction
		Rollingbalance float64
	}
	for idx, trans := range transactions {
		if idx == 0 {
			datasets = append(datasets, struct {
				Tx             transaction
				Rollingbalance float64
			}{Tx: trans, Rollingbalance: trans.Value})
		} else {
			datasets = append(datasets, struct {
				Tx             transaction
				Rollingbalance float64
			}{Tx: trans, Rollingbalance: datasets[idx-1].Rollingbalance + trans.Value})
		}
	}

	w.WriteHeader(200)
	AdminPage(AdminRender("tx-list", datasets), w)
}

func AdminInventory(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		name := r.PostFormValue("Name")
		desc := r.PostFormValue("Description")
		newitem := item{Name: name, Description: desc, CrewID: uint(id)}
		database.Create(&newitem)
	}

	var items []item
	filtered := item{CrewID: uint(id)}
	database.Where(&filtered).Find(&items)

	w.WriteHeader(200)
	AdminPage(AdminRender("inventory-list", items), w)
}

func AdminDeleteItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var todel item
	database.First(&todel, id)
	crew := todel.CrewID
	database.Delete(&todel)
	http.Redirect(w, r, fmt.Sprintf("/admin/inventory?id=%d", crew), 302)
}

func AdminEditItem(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var toedit item
	database.First(&toedit, id)
	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		toedit.Name = r.PostFormValue("Name")
		toedit.Description = r.PostFormValue("Description")
		database.Save(&toedit)
		http.Redirect(w, r, fmt.Sprintf("/admin/inventory?id=%d", toedit.CrewID), 302)
		return
	} else {
		w.WriteHeader(200)
		AdminPage(AdminRender("inventory-edit", toedit), w)
	}
}

func AdminMembers(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if r.Method == "POST" {
		r.ParseMultipartForm(268435456)
		operation := r.PostFormValue("operation")
		var dataset member
		if operation == "update" {
			id, _ := strconv.Atoi(r.PostFormValue("id"))
			database.First(&dataset, id)
		} else if operation == "create" {
			name := r.PostFormValue("name")
			dataset = member{Name: name, CrewID: uint(id)}
			database.Save(&dataset)
		}
		file, _, err := r.FormFile("image")
		if err != nil {
			log.Fatal(err)
		}
		filename := dataset.Filename()
		output(func(print printer) {
			print(filename)
			if file == nil {
				print("file is nil")
			}
		})
		os.Remove(filename)
		destination, _ := os.Create(filename)
		defer destination.Close()
		io.Copy(destination, file)
	}

	var members []member
	filter := member{CrewID: uint(id)}
	database.Where(&filter).Find(&members)

	w.WriteHeader(200)
	AdminPage(AdminRender("members-list", members), w)
}

func AdminMembersDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var dataset member
	database.First(&dataset, id)
	crewId := dataset.CrewID
	os.Remove(dataset.Filename())
	database.Delete(&dataset)
	http.Redirect(w, r, fmt.Sprintf("/admin/crew?id=%d", crewId), 302)
}

func AdminChat(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var dataset crew
	database.Preload("Chat").First(&dataset, id)
	//we need the telegram user.
	chat := dataset.Chat

	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		text := r.PostFormValue("message")
		chat.sendMessage(text)
	}

	messages := chat.fetchMessages()
	lastOne := len(messages)
	firstOne := lastOne - 25
	if firstOne < 0 {
		firstOne = 0
	}
	toPrint := messages[firstOne:lastOne]
	w.WriteHeader(200)
	AdminPage(AdminRender("chats", toPrint), w)
}

func AdminContacts(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		name := r.PostFormValue("name")
		description := r.PostFormValue("description")
		newcontact := contact{Name: name, Description: description, OwnerID: uint(id)}
		database.Save(&newcontact)
	}
	var dataset crew
	database.Preload("Contacts").First(&dataset, id)
	contacts := dataset.Contacts
	w.WriteHeader(200)
	AdminPage(AdminRender("contacts-list", contacts), w)
}

func AdminContactsDelete(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var dataset contact
	database.First(&dataset, id)
	crew := dataset.OwnerID
	database.Delete(&dataset)
	http.Redirect(w, r, fmt.Sprintf("/admin/contacts?id=%d", crew), 302)
}

func AdminContactsChat(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	var dataset contact
	database.First(&dataset, id)
	if r.Method == "POST" {
		r.ParseMultipartForm(262144)
		text := r.PostFormValue("message")
		dataset.sendMessageToCrew(text)
	}
	messages := dataset.fetchSpacemail()
	lastOne := len(messages)
	firstOne := lastOne - 25
	if firstOne < 0 {
		firstOne = 0
	}
	toPrint := messages[firstOne:lastOne]
	w.WriteHeader(200)
	data := struct {
		Messages []spacemail
		Crew     uint
		Contact  string
	}{
		Messages: toPrint,
		Crew:     dataset.OwnerID,
		Contact:  dataset.Name,
	}
	AdminPage(AdminRender("contacts-chat", data), w)
}

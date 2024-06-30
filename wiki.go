package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"text/template"
	"web-app/structs"
)

var templatesPath = "./tmpl/"

var templates = template.Must(template.ParseFiles(templatesPath+"edit.html", templatesPath+"view.html"))

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := structs.LoadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := structs.LoadPage(title)
	if err != nil {
		p = &structs.Page{Title: title}
	}

	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")

	p := &structs.Page{Title: title, Body: []byte(body)}
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func makeHandler(fn func(w http.ResponseWriter, r *http.Request, title string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}

		fn(w, r, m[2])
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/view/FrontPage", http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *structs.Page) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
// 	m := validPath.FindStringSubmatch(r.URL.Path)
// 	if m == nil {
// 		http.NotFound(w, r)
// 		return "", errors.New("invalid structs.Page Title")
// 	}

// 	return m[2], nil // the title is the second subexpression
// }

// TODO:
// 	Spruce up the page templates by making them valid HTML and adding some CSS rules
// 	Implement inter-page linking by converting instances of [PageName] to <a href="/view/PageName</a>". (hint: you could use regexp.ReplaceAllFunc to do this)

func main() {
	fmt.Println("server is running")
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/view/", makeHandler(viewHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

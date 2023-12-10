package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
)

type QueryParams struct {
	EditType string
}

type Page struct {
	Title string
	Body  []byte
}

type SpecialPage struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func render(w http.ResponseWriter, name string, page *Page, params url.Values) {
	t, err := template.ParseFiles(name + ".html")

	if err != nil {
		// Handle the error in some way. For example:
		fmt.Println("error when parsing html file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	parsedParams := QueryParams{
		EditType: params.Get("editType"),
	}

	executeErr := t.Execute(w, struct {
		Page   *Page
		Params QueryParams
	}{Page: page, Params: parsedParams})

	if executeErr != nil {
		fmt.Printf(`error executing %s with %s`, name, page)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/"+title+"?editType=new", http.StatusFound)
		return
	}
	render(w, "view", p, r.URL.Query())

	// rawdog
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]

	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	render(w, "edit", p, r.URL.Query())

	// rawdog
	// fmt.Fprintf(w, `
	// <h1>Editing %s</h1>
	// <form action="/save/%s" method="POST">
	// <textarea name="body">%s</textarea><br>
	// <input type="submit" value="Save">
	// </form>
	// `, p.Title, p.Title, p.Body)
}

func saveHandler() {}

func main() {
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	// http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

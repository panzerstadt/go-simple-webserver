package view

import (
	"fmt"
	"html/template"
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

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func render(w http.ResponseWriter, name string, page *Page, params url.Values) {
	t, err := template.ParseFiles("templates/" + name + ".html")

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

func ViewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Query().Get("title")

	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/api/edit?title="+title+"&editType=new", http.StatusFound)
		return
	}
	render(w, "view", p, r.URL.Query())

	// rawdog
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

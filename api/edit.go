package edit

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
	t, err := template.ParseFiles("../public/" + name + ".html")

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

func EditHandler(w http.ResponseWriter, r *http.Request) {
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

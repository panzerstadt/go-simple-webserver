package wiki

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

type SpecialPage struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func LoadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func Render(w http.ResponseWriter, name string, page *Page, params url.Values) {
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

func Serve(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "heeeeeeyyyyy")
}

// func serve() {
// 	http.HandleFunc("/view/", viewHandler)
// 	http.HandleFunc("/edit/", editHandler)
// 	// http.HandleFunc("/save/", saveHandler)
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

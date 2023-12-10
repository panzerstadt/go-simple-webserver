package main

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
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
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]*)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

var templates = template.Must(template.ParseFiles("templates/edit.html", "templates/view.html", "templates/list.html"))

func render(w http.ResponseWriter, name string, page *Page, params url.Values) {
	parsedParams := QueryParams{
		EditType: params.Get("editType"),
	}

	err := templates.ExecuteTemplate(w, name+".html", struct {
		Page   *Page
		Params QueryParams
	}{Page: page, Params: parsedParams})

	if err != nil {
		// Handle the error in some way. For example:
		fmt.Println("error when parsing html file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("trying to view: ", title)
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/edit/?title="+title+"&editType=new", http.StatusFound)
		return
	}
	render(w, "view", p, r.URL.Query())

	// rawdog
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	fmt.Println("trying to edit: ", title)
	p, err := loadPage(title)
	if err != nil {
		titleParam := r.URL.Query().Get("title")
		p = &Page{Title: titleParam}
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

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: strings.ReplaceAll(title, " ", "_"), Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}
func redirectHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/list/", http.StatusFound)
}
func listHandler(w http.ResponseWriter, r *http.Request) {
	dataDir := "data"
	files, err := os.ReadDir(dataDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	filenames := make([]string, len(files))
	for i, file := range files {
		if file.IsDir() {
			continue
		}
		filename := strings.Split(file.Name(), ".txt")
		filenames[i] = filename[0]
	}

	err = templates.ExecuteTemplate(w, "list.html", struct {
		Filenames []string
	}{Filenames: filenames})

	if err != nil {
		// Handle the error in some way. For example:
		fmt.Println("error when parsing html file")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", redirectHandler)
	http.HandleFunc("/list/", listHandler)
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	fmt.Printf("ready and listening at port %s\n", port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

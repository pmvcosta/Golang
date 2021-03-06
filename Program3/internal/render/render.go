package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/justinas/nosurf"
	"github.com/pmvcosta/bookings/internal/config"
	"github.com/pmvcosta/bookings/internal/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	//Add information to be available on every page
	td.CSRFToken = nosurf.Token(r)
	return td
}

//if name of func begins with lowercase letter, it is private
//func renderTemplate(w http.ResponseWriter, tmpl string) {
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		//Get the template cache from the App Config
		tc = app.TemplateCache
	} else {
		//Rebuilds cache every refresh -- useful for development mode
		tc, _ = CreateTemplateCache()
	}

	t, ok := tc[tmpl] //if template exists, t will have a useful value and ok will be True
	// otherwise t will not have a useful value, and ok will be False

	if !ok {
		log.Fatal("Could not get template from template cache.")
	}

	buf := new(bytes.Buffer)

	//Introduce information meant to be available on every page
	td = AddDefaultData(td, r)

	//_ = t.Execute(buf, nil) //execute value of template to buf, no data is passed
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser:", err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}             //a map of strings to pointers to structs (template.Template)
	pages, err := filepath.Glob("./templates/*.page.tmpl") //find all files containing ".page.tmpl"
	if err != nil {
		fmt.Println("An error occurred!")
		return myCache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		//fmt.Println("Page is currently:", page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			fmt.Println("An error occurred!")
			return myCache, err
		}

		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			fmt.Println("An error occurred!")
			return myCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}

		myCache[name] = ts
	}
	return myCache, nil
}

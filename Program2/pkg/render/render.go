package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/pmvcosta/golangPractice1/pkg/config"
	"github.com/pmvcosta/golangPractice1/pkg/models"
)

var functions = template.FuncMap{}

var app *config.AppConfig

// NewTemplates sets the config for the template package
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	//Add information to be available on every page
	return td
}

//if name of func begins with lowercase letter, it is private
//func renderTemplate(w http.ResponseWriter, tmpl string) {
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {
	//Create template cache (should not be done each time we render a page, that is inefficient)
	//HENCE, application-wide configuration (see several lines below)
	/*tc, err := CreateTemplateCache()
	  if err != nil {
			fmt.Println("Error getting template cache:", err)
			log.Fatal(err) //So that the program does not proceed if templates cannot be retrieved
		}*/

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
	td = AddDefaultData(td)

	//_ = t.Execute(buf, nil) //execute value of template to buf, no data is passed
	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)

	if err != nil {
		fmt.Println("Error writing template to browser:", err)
	}

	//Parse the template (without use of layouts, i.e. CreateTemplateCache)
	/*
		parsedTemplate, _ := template.ParseFiles("./templates/" + tmpl)
		err = parsedTemplate.Execute(w, nil) //No data is passed
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return
		}*/
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
		fmt.Println("Page is currently:", page)
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

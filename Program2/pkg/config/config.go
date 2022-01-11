package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
)

//Ensure this package imports only what it needs to, to avoid running
// into an import cycle

// AppConfig holds the application config, i.e. application wide configuration
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool //Used to assign an adequate value to the session Cookies: true when in production, false when in development
	Session       *scs.SessionManager
}

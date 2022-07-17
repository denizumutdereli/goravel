package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
)

type Render struct {
	Renderer   string `json:"renderer"`
	RootPath   string `json:"root_path"`
	Secure     bool   `json:"secure"`
	Port       string `json:"port"`
	ServerName string `json:"server_name"`
	JetViews   *jet.Set
}

type TemplateData struct {
	IsAuthenticated bool                   `json:"is_authenticated"`
	IntMap          map[string]string      `json:"int_map"`
	FloatMap        map[string]float32     `json:"float_map"`
	Data            map[string]interface{} `json:"data"`
	CSRFToken       string                 `json:"csrf_token"`
	Port            string                 `json:"port"`
	Secure          bool                   `json:"secure"`
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		fmt.Println("go")
		return c.GoPage(w, r, view, variables)

	case "jet":
		fmt.Println("jet")
		return c.JetPage(w, r, view, variables, data)

	default:

	}

	return errors.New("no rendering engine specified")
}

//GoPage rendering standard Go pages
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}
	td := &TemplateData{}

	if data != nil {
		td = data.(*TemplateData)
	}

	err = tmpl.Execute(w, &td)

	if err != nil {
		return err
	}
	return nil
}

//JetPage rendering...
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}

	if data != nil {
		td = data.(*TemplateData)
	}

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))

	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

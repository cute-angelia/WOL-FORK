package internal

import (
	"html/template"
	"net/http"
)

func RenderHomePage(w http.ResponseWriter, r *http.Request) {

	tmpl, _ := template.ParseFiles("pages/index.html")
	tmpl.Execute(w, ComputerList)

}

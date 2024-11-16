package controller

import (
	"fmt"
	"github.com/albakov/go-tennis-scoreboard/internal/util"
	"html/template"
	"net/http"
)

const f = "controller.Controller"

type ServerResponse interface {
	ShowResponse(w http.ResponseWriter, statusCode int, templatePath string, data PageData)
	ShowError(w http.ResponseWriter, statusCode int, errorMessage ErrorMessage)
	BackWithError(w http.ResponseWriter, templatePath string, data PageData)
	ShowMethodNotAllowedError(w http.ResponseWriter)
	ShowNotFound(w http.ResponseWriter)
	ShowServerError(w http.ResponseWriter)
}

type Controller struct {
}

type ErrorMessage struct {
	PageTitle    string
	Title        string
	ErrorMessage any
	Code         int
}

type PageData struct {
	PageTitle    string
	ErrorMessage any
	Data         any
}

func New() *Controller {
	return &Controller{}
}

func (c *Controller) ShowResponse(w http.ResponseWriter, statusCode int, templatePath string, data PageData) {
	const op = "ShowResponse"

	c.setHeaders(w)
	w.WriteHeader(statusCode)

	tmpl := template.Must(template.ParseFiles(
		fmt.Sprintf("view/%s", templatePath),
		"view/header.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, data)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) ShowError(w http.ResponseWriter, statusCode int, errorMessage ErrorMessage) {
	const op = "ShowError"

	c.setHeaders(w)
	w.WriteHeader(statusCode)

	tmpl := template.Must(template.ParseFiles(
		"view/error.html",
		"view/header.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, errorMessage)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) ShowServerError(w http.ResponseWriter) {
	const op = "ShowServerError"

	c.setHeaders(w)
	w.WriteHeader(http.StatusBadRequest)

	tmpl := template.Must(template.ParseFiles(
		"view/error.html",
		"view/header.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, ErrorMessage{
		PageTitle:    MessageServerError,
		Title:        MessageServerError,
		ErrorMessage: MessageServerError,
		Code:         http.StatusBadRequest,
	})
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) ShowNotFound(w http.ResponseWriter) {
	const op = "ShowNotFound"

	c.setHeaders(w)
	w.WriteHeader(http.StatusNotFound)

	tmpl := template.Must(template.ParseFiles(
		"view/error.html",
		"view/header.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, ErrorMessage{PageTitle: "Not Found", Title: "Not Found", Code: http.StatusNotFound})
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) BackWithError(w http.ResponseWriter, templatePath string, data PageData) {
	const op = "BackWithError"

	c.setHeaders(w)
	w.WriteHeader(http.StatusUnprocessableEntity)

	tmpl := template.Must(template.ParseFiles(
		fmt.Sprintf("view/%s", templatePath),
		"view/header.html",
		"view/footer.html",
	))
	err := tmpl.Execute(w, data)
	if err != nil {
		util.LogError(f, op, err)
	}
}

func (c *Controller) ShowMethodNotAllowedError(w http.ResponseWriter) {
	c.ShowError(
		w,
		http.StatusMethodNotAllowed,
		ErrorMessage{Title: MessageMethodNotAllowed, Code: http.StatusMethodNotAllowed},
	)
}

func (c *Controller) setHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
}

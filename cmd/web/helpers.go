package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Internal Server Error
func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)
	app.logger.Error(err.Error(), "method", method, "uri", uri, "trace", trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

// Error on Client
func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// status not found error
func (app *application) notFound(w http.ResponseWriter) {
	app.clientError(w, http.StatusNotFound)
}

// Render Helper
func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {
	// Retrieve the appropriate template set from the cache based on the page
	// name (like 'home.tmpl'). If no entry exists in the cache with the
	// provided name, then create a new error and call the serverError() helper
	// method that we made earlier and return.

	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, r, err)
		return
	}

	// Initialize a new buffer.
	buf := new(bytes.Buffer)

	// Write template to buffer instead of http.ResponseWriter
	// If handle if error rendering template
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// If template is written to buffer withour erros, we are safe
	// to write HTTP status code to http.ResponseWriter
	w.WriteHeader(status)

	// Write contents of buffer to http.ResponseWriter
	buf.WriteTo(w)
}

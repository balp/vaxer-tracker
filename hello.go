package vaxer

import (
        "html/template"
        "net/http"
        "time"

        "appengine"
        "appengine/datastore"
        //"appengine/user"
)

type Greeting struct {
        Author  string
        Content string
        Date    time.Time
}

func init() {
    http.Handle("/static/", http.FileServer(http.Dir(".")))
    http.HandleFunc("/", vaxerIndex)
    http.HandleFunc("/plant", vaxerPlant)
}

var (
  templates = template.Must(template.ParseFiles(
     "header.html",
     "footer.html",
     "vaxer.html",
     "error.html",
  ))
)

type Variables struct{
	Title string
	Heading string
}

func vaxerIndex(w http.ResponseWriter, r *http.Request) {
    data := Variables{"Balp", "One"}
    if err := templates.ExecuteTemplate(w, "vaxer.html", data); err != nil {
	writeError(w, r, err)
    }
}

func vaxerPlant(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/", http.StatusFound)
}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	c := appengine.NewContext(r)
	c.Errorf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	if err := templates.ExecuteTemplate(w, "error.html", err); err != nil {
		c.Errorf("templates.ExecuteTemplate: %v", err)
	}
}

func guestbookKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

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
        http.HandleFunc("/", hello)
//        http.HandleFunc("/sign", sign)
}

var (
  templates = template.Must(template.ParseFiles(
     "header.html",
     "footer.html",
     "hello.html",
     "error.html",
  ))
)

type Variables struct{
	Title string
	Heading string
}

func hello(w http.ResponseWriter, r *http.Request) {
    data := Variables{"Balp", "One"}
    if err := templates.ExecuteTemplate(w, "hello.html", data); err != nil {
	writeError(w, r, err)
    }

}

func writeError(w http.ResponseWriter, r *http.Request, err error) {
	c := appengine.NewContext(r)
	c.Errorf("Error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
	if err := templates.ExecuteTemplate(w, "error.html", err); err != nil {
		c.Errorf("templates.ExecuteTemplate: %v", err)
	}
}

// guestbookKey returns the key used for all guestbook entries.
func guestbookKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    return datastore.NewKey(c, "Guestbook", "default_guestbook", 0, nil)
}

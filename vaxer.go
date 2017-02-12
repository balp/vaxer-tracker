package vaxer

import (
        "html/template"
        "net/http"
        "time"
        "strconv"

        "appengine"
        "appengine/datastore"
        //"appengine/log"
        //"appengine/user"
)

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

type Plant struct {
    Name string
    Planted time.Time
    Where string
    Row int
    Colum int
    Next string
    When time.Time
}

type PlantType struct {
    Name string
    Id string
}

type WebContext struct {
    PlantTypes []PlantType
    Plants []Plant
}

func vaxerIndex(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    q := datastore.NewQuery("Plant").Ancestor(plantsKey(c)).Order("-When").Limit(100)
    webContext := new(WebContext)
    webContext.Plants = make([]Plant, 0, 100)
    webContext.PlantTypes = []PlantType{{"Basil", "basil"},{"Coriander", "coriander"}}
    if _, err := q.GetAll(c, &webContext.Plants); err != nil {
        writeError(w, r, err)
    }


    if err := templates.ExecuteTemplate(w, "vaxer.html", webContext); err != nil {
	writeError(w, r, err)
    }
}

func vaxerPlant(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    var row int
    var column int
    var err error
    if row, err = strconv.Atoi(r.FormValue("row")); err != nil {
	writeError(w, r, err)
        return
    }
    if column, err = strconv.Atoi(r.FormValue("column")); err != nil {
	writeError(w, r, err)
        return
    }
    plant := Plant{
        Name: r.FormValue("plantType"),
        Planted: time.Now(),
        Where: r.FormValue("location"),
        Row: row,
        Colum: column,
        Next: "Cultivate",
        When: time.Now(),
    }
    key := datastore.NewIncompleteKey(c, "Plant", plantsKey(c))
    _, serr := datastore.Put(c, key, &plant)
    if serr != nil {
        writeError(w, r, err)
        return
    }
    //log.Infof(c, "Datastore put successful")
    c.Debugf("Storing %v in database %v", plant, key)
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

func plantsKey(c appengine.Context) *datastore.Key {
    return datastore.NewKey(c, "Plants", "default_plants", 0, nil)
}

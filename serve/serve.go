package serve

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"naivebayes/app"
	"naivebayes/data"
)

type Att struct {
	attname  string
	attvalue string
}

func (att *Att) addval(d *data.Data, aname, nval string) {
	for i := range *d {
		if (*d)[i].AttrName == aname {
			(*d)[i].AttrValues = append((*d)[i].AttrValues, nval)
			break
		}
	}
	return
}

func (att *Att) delatt(d *data.Data, aname string) {
	for i := range *d {
		if (*d)[i].AttrName == aname {
			*d = append((*d)[:i], (*d)[i+1:]...)
			break
		}
	}
	return
}

func (att *Att) delval(d *data.Data, aname, avalue string) {
	var p1, p2 int
	for i := range *d {
		if (*d)[i].AttrName == aname {
			p1 = i
			if len((*d)[i].AttrValues) == 1 {
				(*d)[i].AttrValues[0] = "VALUE"
				return
			}
			for j := range (*d)[i].AttrValues {
				if (*d)[i].AttrValues[j] == avalue {
					p2 = j
					break
				}
			}
			break
		}
	}
	(*d)[p1].AttrValues = append((*d)[p1].AttrValues[:p2], (*d)[p1].AttrValues[p2+1:]...)
	return
}

func (att *Att) savatt(d *data.Data, aname, nname string) {
	for i := range *d {
		if (*d)[i].AttrName == aname {
			(*d)[i].AttrName = nname
			break
		}
	}
	return
}

func (att *Att) savval(d *data.Data, aname, aval, nval string) {
	for i := range *d {
		if (*d)[i].AttrName == aname {
			for j := range (*d)[i].AttrValues {
				if (*d)[i].AttrValues[j] == aval {
					(*d)[i].AttrValues[j] = nval
					break
				}
			}
			break
		}
	}
	return
}

func (att *Att) renderTemplate(w http.ResponseWriter, tmpl string, p interface{}) {
	t, err := template.ParseFiles("serve/" + tmpl + ".html")
	if err != nil {
		log.Fatal(err)
	}
	err = t.Execute(w, p)
	if err != nil {
		log.Fatal(err)
	}
}

func (att *Att) attributesHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/Attributes/"):]
	body, err := app.LoadFile(title)
	if err != nil {
		log.Fatal(err)
	}
	p := &data.Page{Title: title, Body: (*body), Attname: att.attname, Attvalue: att.attvalue}
	att.renderTemplate(w, "attributes", p)
}

func (att *Att) saveAttrHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/saveAttributes/"):]
	bval := r.FormValue("submit")
	if bval == "done" {
		http.Redirect(w, r, "/CreateData/"+title, http.StatusFound)
	} else {
		d, err := app.LoadFile(title)
		if err != nil {
			log.Fatal(err)
		}
		if bval == "addatt" {
			val := r.FormValue("newatt")
			if val != "" {
				*d = append(*d, data.Attribute{val, []string{"VALUE"}})
			}
		} else {
			barr := strings.Split(bval, "-")
			if barr[0] == "add" {
				val := r.FormValue("inpval-" + barr[1])
				if val != "" {
					att.addval(d, barr[1], val)
				}
			} else if barr[0] == "del" {
				if len(barr) == 2 {
					att.delatt(d, barr[1])
				} else {
					att.delval(d, barr[1], barr[2])
				}
			} else if barr[0] == "sav" {
				if len(barr) == 2 {
					val := r.FormValue("inp-" + barr[1])
					if val != "" {
						att.savatt(d, barr[1], val)
					}
				} else {
					val := r.FormValue("inp-" + barr[1] + "-" + barr[2])
					if val != "" {
						att.savval(d, barr[1], barr[2], val)
					}
				}
			} else if barr[0] == "edt" {
				att.attname = barr[1]
				if len(barr) == 3 {
					att.attvalue = barr[2]
				}
			}
		}
		err = app.SaveFile(title, d)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.Redirect(w, r, "/Attributes/"+title, http.StatusFound)
}

func (att *Att) createDataHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/CreateData/"):]
	file := true
	_, err := app.LoadRndFile(title)
	if err != nil {
		file = false
	}
	att.renderTemplate(w, "createdata", &data.CrData{Title: title, File: file})
}

func (att *Att) dataHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/Data/"):]
	bval := r.FormValue("submit")
	if bval == "back" {
		http.Redirect(w, r, "/Attributes/"+title, http.StatusFound)
	}
	d, err := app.LoadFile(title)
	if err != nil {
		log.Fatal(err)
	}
	if bval == "newdata" {
		n, _ := strconv.Atoi(r.FormValue("num"))
		err = app.CreateRndFile(title, d, n)
		if err != nil {
			log.Fatal(err)
		}
	}
	body, err := app.LoadRndFile(title)
	if err != nil {
		log.Fatal(err)
	}
	att.renderTemplate(w, "data", &data.Datapage{Title: title, Body: (*body), Adata: (*d)})
}

func (att *Att) saveDataHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/saveData/"):]
	tdata, err := app.LoadRndFile(title)
	if err != nil {
		log.Fatal(err)
	}
	bval := r.FormValue("submit")
	if bval == "back" {
		http.Redirect(w, r, "/Attributes/"+title, http.StatusFound)
	} else if bval == "done" {
		http.Redirect(w, r, "/Query/"+title, http.StatusFound)
	} else if bval == "add" {
		d, err := app.LoadFile(title)
		if err != nil {
			log.Fatal(err)
		}
		var slice []string
		for i := 0; i < len(*d); i++ {
			slice = append(slice, r.FormValue("sel-"+(*d)[i].AttrName))
		}
		*tdata = append(*tdata, slice)
	} else {
		slice := strings.Split(bval, "-")
		pos, _ := strconv.Atoi(slice[1])
		*tdata = append((*tdata)[:pos], (*tdata)[pos+1:]...)
	}
	err = app.SaveRndFile(title, tdata)
	if err != nil {
		log.Fatal(err)
	}
	http.Redirect(w, r, "/Data/"+title, http.StatusFound)
}

func (att *Att) queryHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/Query/"):]
	d, err := app.LoadFile(title)
	if err != nil {
		log.Fatal(err)
	}
	att.renderTemplate(w, "query", &data.Page{Title: title, Body: (*d)})
}

func (att *Att) answerHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/Answer/"):]
	//bval := r.FormValue("submit")
	d, err := app.LoadFile(title)
	if err != nil {
		log.Fatal(err)
	}

	count := 0
	var qslice []string
	for i := 0; i < len((*d)); i++ {
		val := r.FormValue("sel-" + (*d)[i].AttrName)
		if val == "?" {
			count++
		}
		qslice = append(qslice, val)
	}
	if count != 1 {
		http.Redirect(w, r, "/Query/"+title, http.StatusFound)
	} else {
		tdata, err := app.LoadRndFile(title)
		if err != nil {
			log.Fatal(err)
		}
		ans := app.FindAns(tdata, &qslice, d)
		att.renderTemplate(w, "answer", &data.Anspage{Title: title, Body: (*d), Query: qslice, Ans: ans})
	}
}

func Run() {
	att := &Att{
		attname:  "",
		attvalue: "",
	}
	http.HandleFunc("/Attributes/", att.attributesHandler)
	http.HandleFunc("/saveAttributes/", att.saveAttrHandler)
	http.HandleFunc("/CreateData/", att.createDataHandler)
	http.HandleFunc("/Data/", att.dataHandler)
	http.HandleFunc("/saveData/", att.saveDataHandler)
	http.HandleFunc("/Query/", att.queryHandler)
	http.HandleFunc("/Answer/", att.answerHandler)
	http.ListenAndServe(":8000", nil)
}

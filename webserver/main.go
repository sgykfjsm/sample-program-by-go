package main

import (
	"fmt"
	"net/http"

	"os"
	"os/signal"
	"syscall"

	"crypto/rand"
	"encoding/binary"
	"html/template"
	"strconv"

	"log"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/gorilla/mux"
)

var formHTML = `
<!DOCTYPE html>
<html>
<body>
<h1>Get</h1>
<form action="/action" method="get">
  First name: <input type="text" name="fname"><br>
  Last name: <input type="text" name="lname"><br>
  <input type="submit" value="Submit">
</form>
<hr>
<h1>Post</h1>
<form action="/action" method="post">
  First name: <input type="text" name="fname"><br>
  Last name: <input type="text" name="lname"><br>
  <input type="submit" value="Submit">
</form>
</body>
</html>`

var formTempl = template.Must(template.New("form").Parse(formHTML))
var actionTempl = template.Must(template.New("action").ParseFiles("action.html"))

type Person struct {
	FirstName string
	LastName  string
}

func random() string {
	var n uint64
	binary.Read(rand.Reader, binary.BigEndian, &n)
	return strconv.FormatUint(n, 36)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Index Page")
}

func FormHandler(w http.ResponseWriter, r *http.Request) {
	if err := formTempl.Execute(w, nil); err != nil {
		log.Fatalf("Failed to generate HTML: %s", err.Error())
	}
}

func ActionHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "Failed to parse query: %#v", err.Error())
		return
	}

	var p Person
	p.FirstName = r.Form.Get("fname")
	p.LastName = r.Form.Get("lname")

	//h := fmt.Sprintf(html, p.FirstName, p.LastName, random(), random())

	data := struct {
		Person        *Person
		NextFirstName string
		NextLastName  string
	}{
		&p,
		random(),
		random(),
	}

	if err := actionTempl.ExecuteTemplate(w, "action.html", data); err != nil {
		log.Fatalf("Failed to return html: %s", err.Error())
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			s := <-sig
			fmt.Printf("\nCatch the signal: %s\n", s.String())
			done <- true
			break
		}
	}()

	addr := ":8081"
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/param", ActionHandler)
	r.HandleFunc("/form", FormHandler)
	r.HandleFunc("/action", ActionHandler)
	gracehttp.Serve(&http.Server{Addr: addr, Handler: r})

	<-done
}

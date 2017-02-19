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
	"github.com/gorilla/websocket"
	"net/url"
	"time"
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

var (
	formTempl   = template.Must(template.New("form").Parse(formHTML))
	actionTempl = template.Must(template.New("action").ParseFiles("action.html"))
	wsTemple    = template.Must(template.New("websocket").ParseFiles("websocket.html"))
	upgrader    = websocket.Upgrader{}
	addr        = ":8081"
)

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

func WebSocketClient() (*websocket.Conn, url.URL, error) {
	u := url.URL{Scheme: "ws", Host: "localhost" + addr, Path: "/ws"}
	log.Printf("Connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	return c, u, err
}

func WebSocket(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade gorilla/websocket: %s", err.Error())
		return
	}

	for {
		m, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message: %s", err.Error())
			return
		}
		log.Printf("recieve: %s", string(message))
		err = c.WriteMessage(m, message)
		if err != nil {
			log.Printf("Failed to write message: %s", err.Error())
			return
		}
	}
}

func WebSocketHandler(w http.ResponseWriter, r *http.Request) {
	if err := wsTemple.ExecuteTemplate(w, "websocket.html", "ws://"+r.Host+"/ws"); err != nil {
		log.Fatalf("Failed to return html: %s", err.Error())
	}
}

func main() {
	sig := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sig, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)

	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/param", ActionHandler)
	r.HandleFunc("/form", FormHandler)
	r.HandleFunc("/action", ActionHandler)
	r.HandleFunc("/ws", WebSocket)
	r.HandleFunc("/websocket", WebSocketHandler)

	var c *websocket.Conn
	t0 := time.NewTicker(time.Second)
	go func() {
		defer t0.Stop()
		var err error
		var u url.URL
		for {
			select {
			case <-t0.C:
				c, u, err = WebSocketClient()
				if err != nil {
					log.Printf("Failed to connect to %s: %s", u.String(), err.Error())
				}
				if c != nil {
					log.Printf("Connected to %s", u.String())
					return
				}
			}
		}
	}()
	defer c.Close()

	t1 := time.NewTicker(1 * time.Second)
	go func() {
		defer t1.Stop()
		for {
			select {
			case _t := <-t1.C:
				if c != nil {
					err := c.WriteMessage(websocket.TextMessage, []byte(_t.String()))
					if err != nil {
						log.Printf("Failed to write message: %s", err.Error())
						return
					}
				}
			}
		}
	}()

	go func() {
		s := <-sig
		fmt.Printf("\nCatch the signal: %s\n", s.String())
		err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		if err != nil {
			log.Println("Failed to write close message: %s", err.Error())
			return
		}
		done <- true
	}()

	log.Println("Start Servre")
	gracehttp.Serve(&http.Server{Addr: addr, Handler: r})
	<-done
}

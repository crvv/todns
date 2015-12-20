package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/crvv/todns/config"
	"github.com/crvv/todns/data"
	"github.com/miekg/dns"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

func StartHttpServer() {
	http.HandleFunc("/", file)
	http.HandleFunc("/api/", api)
	go func() {
		err := http.Serve(newTlsListener(), nil)
		log.Fatal(err)
	}()
}

type ApiBody struct {
	Upstream string `json:"upstream"`
	Ttl      int    `json:"ttl"`
	Name     string `json:"name"`
	Address  string `json:"address"`
}

func parseBody(b io.Reader) (body ApiBody, err error) {
	buffer, err := ioutil.ReadAll(b)
	if err != nil {
		return
	}
	err = json.Unmarshal(buffer, &body)
	if err != nil {
		return
	}
	log.Printf("Request body: %v", body)
	return
}
func api(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	log.Printf("API Request, RemoteAddr: %v, RequestPath: %v\n", r.RemoteAddr, path)
	var body ApiBody
	var err error
	if r.Method != "GET" {
		body, err = parseBody(r.Body)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
	}
	switch r.Method {
	case "PUT":
		switch path {
		case "api/ttl":
			if body.Ttl < 1 {
				http.Error(w, fmt.Sprintf("ttl %v is too small", body.Ttl), 400)
				return
			}
			config.SetTtl(body.Ttl)
		case "api/upstream":
			msg := &dns.Msg{}
			msg.SetQuestion("google.com.", dns.TypeA)
			_, err := dns.Exchange(msg, body.Upstream)
			if err != nil {
				http.Error(w, err.Error(), 400)
				return
			}
			config.SetUpstream(body.Upstream)
		default:
			http.NotFound(w, r)
			return
		}
	case "POST":
		switch path {
		case "api/record":
			if net.ParseIP(body.Address) == nil {
				http.Error(w, fmt.Sprintf("%v is not a valid IP address", body.Address), 400)
				return
			}
			config.AddRecord(body.Name, body.Address)
		default:
			http.NotFound(w, r)
			return
		}
	case "DELETE":
		switch path {
		case "api/record":
			config.RemoveRecord(body.Name)
		default:
			http.NotFound(w, r)
			return
		}
	case "GET":
		switch path {
		case "api":
		default:
			http.NotFound(w, r)
			return
		}
	}
	data, _ := json.Marshal(config.GetConfig())
	w.Header().Add("Content-Type", "application/json")
	w.Write(data)
}
func file(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")
	log.Printf("HTTP Request, RemoteAddr: %v, RequestPath: %v\n", r.RemoteAddr, path)
	if path == "" {
		path = "index.html"
	}

	page, err := data.Asset(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	http.ServeContent(w, r, path, time.Now(), bytes.NewReader(page))
}

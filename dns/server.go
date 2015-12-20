package dns

import (
	"fmt"
	"github.com/crvv/todns/config"
	"github.com/miekg/dns"
	"log"
	"os"
	"os/signal"
)

func Start() {
	var err error

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan)

	log.Println("DNS Server start")
	server := &dns.Server{Addr: config.GetDnsListen(), Net: "udp", Handler: &Handler{}}
	go func() {
		err = server.ListenAndServe()
		signalChan <- os.Interrupt
	}()
	<-signalChan
	if err != nil {
		log.Println(err)
	}
	log.Println("DNS Server stop")
	server.Shutdown()
}

type Handler struct{}

func (h *Handler) ServeDNS(w dns.ResponseWriter, request *dns.Msg) {
	log.Printf("Receive dns request from %v\n", w.RemoteAddr())

	response := get(request)

	if response.Answer != nil && len(response.Answer) > 0 {
		log.Println("Answer: ")
		for _, ans := range response.Answer {
			fmt.Println("    ", ans)
		}
	}
	w.WriteMsg(response)
}

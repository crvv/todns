package dns

import (
	"fmt"
	"github.com/crvv/todns/config"
	"github.com/miekg/dns"
	"log"
	"net"
)

func get(request *dns.Msg) (response *dns.Msg) {
	if len(request.Question) < 1 {
		log.Println("no question")
		response = empty(request)
		return
	}
	question := request.Question[0]
	log.Printf("Question: name is '%v', type is '%v'\n", question.Name, dns.TypeToString[question.Qtype])
	if question.Qclass != dns.ClassINET {
		log.Println("query class is not INET")
		response = empty(request)
		return
	}
	if response = local(request); response == nil {
		response = proxy(request)
	}
	return
}

func empty(request *dns.Msg) (response *dns.Msg) {
	response = &dns.Msg{}
	response.SetReply(request)
	return
}
func proxy(request *dns.Msg) (response *dns.Msg) {
	var err error
	response, err = dns.Exchange(request, config.GetUpstream())
	log.Println("Query from upstream dns server", config.GetUpstream())
	if err != nil {
		response = empty(request)
	}
	return
}

func local(request *dns.Msg) (response *dns.Msg) {
	records := config.GetRecord(request.Question[0].Name)
	if records == nil || len(records) == 0 {
		return
	}
	rrs := addrsToRRs(request.Question[0].Name, records, request.Question[0].Qtype)
	response = empty(request)
	response.Answer = rrs
	return response
}
func addrsToRRs(name string, addrs []string, queryType uint16) (rrs []dns.RR) {
	addrs = filterAddr(addrs, queryType)
	log.Printf("find %v address\n", len(addrs))
	for _, addr := range addrs {
		rrString := fmt.Sprintf("%v %v %v %v", name, config.GetTtl(), dns.TypeToString[queryType], addr)
		log.Println(rrString)
		rr, err := dns.NewRR(rrString)
		if err != nil {
			log.Println(err)
			continue
		}
		rrs = append(rrs, rr)
	}
	return
}
func filterAddr(addrs []string, queryType uint16) (result []string) {
	for _, addr := range addrs {
		var ipAddr net.IP
		if ipAddr = net.ParseIP(addr); ipAddr == nil {
			log.Println("net.ParseIP failed", addr)
			continue
		}
		if (ipAddr.To4() != nil && queryType == dns.TypeA) || (ipAddr.To4() == nil && queryType == dns.TypeAAAA) {
			result = append(result, addr)
		}
	}
	return
}

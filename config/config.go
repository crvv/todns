package config

import (
	"encoding/json"
	"github.com/miekg/dns"
	"io/ioutil"
	"sync"
)

type Config struct {
	Crt         string              `json:"crt"`
	Key         string              `json:"key"`
	ClientCrt   string              `json:"clientCrt"`
	HttpsListen string              `json:"httpsListen"`
	DnsListen   string              `json:"dnsListen"`
	Upstream    string              `json:"upstream"`
	Ttl         int                 `json:"ttl"`
	Record      map[string][]string `json:"record"`
}

const ConfigFilename = "./config.json"

var config *Config
var lock sync.RWMutex
var defaultConfig Config = Config{
	Crt:         "./crt",
	Key:         "./key",
	ClientCrt:   "./client",
	HttpsListen: ":443",
	DnsListen:   ":53",
	Upstream:    "8.8.8.8:53",
	Ttl:         300,
	Record:      map[string][]string{"localhost.": []string{"127.0.0.1", "::1", "127.0.0.64"}},
}

func SetUpstream(upstream string) {
	lock.Lock()
	config.Upstream = upstream
	lock.Unlock()
	writeToFile()
}
func SetTtl(ttl int) {
	lock.Lock()
	config.Ttl = ttl
	lock.Unlock()
	writeToFile()
}
func AddRecord(name, addr string) {
	name = dns.Fqdn(name)
	lock.Lock()
	defer writeToFile()
	defer lock.Unlock()
	if _, ok := config.Record[name]; ok {
		config.Record[name] = append(config.Record[name], addr)
		return
	}
	config.Record[name] = []string{addr}
}
func RemoveRecord(name string) {
	name = dns.Fqdn(name)
	lock.Lock()
	delete(config.Record, name)
	lock.Unlock()
	writeToFile()
}
func GetConfig() Config {
	lock.RLock()
	defer lock.RUnlock()
	return *config
}
func GetCrt() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.Crt
}
func GetKey() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.Key
}
func GetClientCrt() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.ClientCrt
}
func GetHttpsListen() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.HttpsListen
}
func GetDnsListen() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.DnsListen
}
func GetRecord(name string) []string {
	lock.RLock()
	defer lock.RUnlock()
	return config.Record[name]
}
func GetTtl() int {
	lock.RLock()
	defer lock.RUnlock()
	return config.Ttl
}
func GetUpstream() string {
	lock.RLock()
	defer lock.RUnlock()
	return config.Upstream
}

func writeToFile() {
	data, err := json.MarshalIndent(config, "", "    ")
	if err == nil {
		ioutil.WriteFile(ConfigFilename, data, 0600)
	}
}
func init() {
	data, err := ioutil.ReadFile(ConfigFilename)
	if err != nil {
		config = &defaultConfig
		data, err = json.MarshalIndent(config, "", "    ")
		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(ConfigFilename, data, 0600)
		if err != nil {
			panic(err)
		}
		return
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}

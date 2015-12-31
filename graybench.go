package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

var (
	pemData []byte
	shortLen int
	fullLen int
	customLen int
)

type LogItem struct {
	Version      string `json:"version"`
	Host         string `json:"host"`
	ShortMessage string `json:"short_message"`
	FullMessage  string `json:"full_message"`
	TimeStamp    int64  `json:"timestamp"`
	Level        string `json:"level"`
	Custom1      string `json:"_custom1,omitempty"`
	Custom2      string `json:"_custom2,omitempty"`
	Custom3      string `json:"_custom3,omitempty"`
	Custom4      string `json:"_custom4,omitempty"`
	Custom5      string `json:"_custom5,omitempty"`
}

func init() {
	rand.Seed(time.Now().UnixNano())
	// In case go < 1.5
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func RandStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func randomItem() LogItem {
	item := LogItem{}
	item.Version = "1.0"
	item.Host = "benchmarkhost"
	item.ShortMessage = RandStringRunes(shortLen)
	item.FullMessage = RandStringRunes(fullLen)
	item.TimeStamp = time.Now().Unix()
	item.Level = "1"
	if (customLen > 0) {
		item.Custom1 = RandStringRunes(customLen)
	    item.Custom2 = RandStringRunes(customLen)
	    item.Custom3 = RandStringRunes(customLen)
	    item.Custom4 = RandStringRunes(customLen)
	    item.Custom5 = RandStringRunes(customLen)
	}
	return item
}
	
func buildClient() *http.Client {
	tlsConfig := &tls.Config{RootCAs: x509.NewCertPool()}
	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client := &http.Client{Transport: transport}
	ok := tlsConfig.RootCAs.AppendCertsFromPEM(pemData)
	if !ok {
		panic("Unable to load CA data")
	}
	return client
}

func sendone(target string, client *http.Client) {
	item := randomItem()
	buf, err := json.Marshal(item)
	if err != nil {
		fmt.Printf("%s: %s\n", time.Now().Format(time.RFC3339), err)
		return
	}

	reader := bytes.NewReader(buf)
	//client := buildClient()

	resp, err := client.Post(target, "text/plain", reader)
	if err != nil {
		fmt.Printf("%s: %s\n", time.Now().Format(time.RFC3339), err)
		return
	}
	defer resp.Body.Close()
}

func benchthread(events int, target string, wg *sync.WaitGroup) {
	client := buildClient()

	for i := 0; i < events; i++ {
		sendone(target, client)
	}
	wg.Done()
}

func main() {
	threads := flag.Int("threads", 10, "amount of threads")
	events := flag.Int("events", 100000, "events per thread")
	shortLenPtr := flag.Int("shortlen", 20, "length of random short message")
	fullLenPtr := flag.Int("fulllen", 200, "length of random full message")
	customLenPtr := flag.Int("customlen", 20, "length of random custom message")
	var ca string
	flag.StringVar(&ca, "ca", "cert.pem", "file with ca certificate chain")
	var target string
	flag.StringVar(&target, "target", "https://graylog.local:12201/gelf", "target HTTP Gelf service")
	
	flag.Parse()
	// Lengths
	shortLen = *shortLenPtr
	fullLen = *fullLenPtr
	customLen = *customLenPtr

	var err error
	pemData, err = ioutil.ReadFile(ca)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s: Threads: %d, Events: %d, CA certificate chain: %s, Target: %s\n", time.Now().Format(time.RFC3339), *threads, *events, ca, target)

	start := time.Now()
	fmt.Printf("%s: Launching threads\n", start.Format(time.RFC3339))
	wg := sync.WaitGroup{}
	for i := 0; i < *threads; i++ {
		wg.Add(1)
		go benchthread(*events, target, &wg)
	}

	wg.Wait()
	end := time.Now()
	fmt.Printf("%s: Threads finished\n", end.Format(time.RFC3339))

	duration := end.Sub(start)
	duration_seconds := duration.Seconds()

	total := ((float64(*threads)) * (float64(*events)))
	eps := total / duration_seconds
	fmt.Printf("%s: Total events: %d, Total time: %ds, EPS: %d\n", time.Now().Format(time.RFC3339), int(total), int(duration_seconds), int(eps))
}

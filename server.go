package main

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"time"
)

const (
	srvSocket   = "localhost:1234"                 //Server socket
	srvProtocol = "tcp4"                           //Connection srvProtocol
	urlProverb  = "https://go-proverbs.github.io/" //Source of Go proverbs
)

type Proverb struct {
	Addrs string
	Text  string
}

func newProverbs() []Proverb {
	html, err := getProverbs()
	//Handle empty content
	if err != nil {
		log.Fatal(err)
	}
	p := parseProverbs(html)
	if len(p) < 1 {
		log.Fatal("There aren't any proverbs!")
	}
	return p
}

func getProverbs() (string, error) {
	resp, err := http.Get(urlProverb) //Get https://go-proverbs.github.io/

	if err != nil {
		return "No any content!", err
	}
	//Close connection
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	//Check 2xx
	if resp.StatusCode != http.StatusOK {
		return "", errors.New(resp.Status)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

func parseProverbs(html string) (ar []Proverb) {
	re := regexp.MustCompile(`<h3><a href="(.+)">(.+)</a></h3>`)
	find := re.FindAllStringSubmatch(html, -1)

	for _, s := range find {
		ar = append(ar, Proverb{Addrs: s[1], Text: s[2]})
	}
	return
}

func handleConn(conn net.Conn, proverb *[]Proverb) {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	log.Println("Connection to the:", conn.RemoteAddr())

	for {
		p := getRandomProverb(*proverb)
		_, _ = conn.Write([]byte(p.Text + "\n\r"))
		_, _ = conn.Write([]byte(p.Addrs + "\n\r"))
		time.Sleep(time.Second * 3)
	}
}

func getRandomProverb(ar []Proverb) Proverb {
	rand.Seed(time.Now().UnixNano())
	return ar[rand.Intn(len(ar))]
}

func main() {
	//Init proverbs
	proverbs := newProverbs()

	//Run server
	listener, err := net.Listen(srvProtocol, srvSocket)
	if err != nil {
		log.Fatal(err)
	}
	defer func(listener net.Listener) {
		_ = listener.Close()
	}(listener)

	//Handle connections
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
		}
		go handleConn(conn, &proverbs)
	}
}

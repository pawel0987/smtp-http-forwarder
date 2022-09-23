package main

import (
    "errors"
    "os"
    "io"
    //"io/ioutil"
    "log"
    "time"
    "strings"
    "net/http"
    "crypto/tls"

    "github.com/emersion/go-smtp"
)

var smtp_username = os.Getenv("SMTP_USERNAME")
var smtp_password = os.Getenv("SMTP_PASSWORD")
var server_domain = os.Getenv("SERVER_DOMAIN")

type HttpEndpoint struct {
    Email string
    Endpoint string
}

var http_endpoints = []HttpEndpoint{}

// The Backend implements SMTP server methods.
type Backend struct{}

func (bkd *Backend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
    log.Println("Trying without auth...")
    //return &Session{}, nil
    return nil, smtp.ErrAuthRequired
}

func (bkd *Backend) Login(state *smtp.ConnectionState, username string, password string) (smtp.Session, error) {
    log.Println("Trying with auth...")
    if username != smtp_username || password != smtp_password {
        log.Println("Invalid credentials")
        return nil, errors.New("Invalid username or password")
    }
    return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct{}

func (s *Session) Mail(from string, opts smtp.MailOptions) error {
    log.Println("Mail from:", from)
    for _, http_endpoint := range http_endpoints {
        if http_endpoint.Email == from {
            transCfg := &http.Transport{
                TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
            }
            client := &http.Client{Transport: transCfg}
            log.Println("Found matching http endpoint. Sending request...")
            _, err := client.Get(http_endpoint.Endpoint)
            if err != nil {
               log.Println(err)
            }
            break
        }
    }
    return nil
}

func (s *Session) Rcpt(to string) error {
    log.Println("Rcpt to:", to)
    return nil
}

func (s *Session) Data(r io.Reader) error {
    //if b, err := ioutil.ReadAll(r); err != nil {
    //    return err
    //} else {
    //    log.Println("Data:", string(b))
    //}
    return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
    return nil
}

func main() {
    if smtp_username == "" || smtp_password == "" {
        panic(errors.New("both SMTP_USERNAME and SMTP_PASSWORD must be set"))
    }
	
    for _, e := range os.Environ() {
	pair := strings.SplitN(e, "=", 2)
        if strings.HasPrefix(pair[0], "HTTP_ENDPOINT") {
            entry_pair := strings.SplitN(pair[1], " ", 2)
	        http_endpoints = append(http_endpoints, HttpEndpoint{
		        Email: entry_pair[0],
	            Endpoint: entry_pair[1],
            })
        }
    }
    log.Println("HTTP Endpoints: ", http_endpoints)

    be := &Backend{}

    s := smtp.NewServer(be)
    s.Addr = ":25"
    s.Domain = server_domain
    s.ReadTimeout = 120 * time.Second
    s.WriteTimeout = 120 * time.Second
    s.MaxMessageBytes = 1024 * 1024
    s.MaxRecipients = 50
    s.AllowInsecureAuth = true

    log.Println("Starting server at", s.Addr)
    if err := s.ListenAndServe(); err != nil {
        log.Fatal(err)
    }
}

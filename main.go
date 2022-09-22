package main

import (
    "errors"
    "io"
    "io/ioutil"
    "log"
    "time"
    "strings"

    "github.com/emersion/go-smtp"
)

var smtp_username = os.Getenv("SMTP_USERNAME")
var smtp_password = os.Getenv("SMTP_PASSWORD")

tyle HttpEndpoint struct {}

var http_endpoints = []HttpEndpoint

// The Backend implements SMTP server methods.
type Backend struct{}

func (bkd *Backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
    return &Session{}, nil
}

// A Session is returned after EHLO.
type Session struct{}

func (s *Session) AuthPlain(username, password string) error {
    if username != smtp_username || password != smtp_password {
        return errors.New("Invalid username or password")
    }
    return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
    log.Println("Mail from:", from)
    for _, http_endpoint := range http_endpoints {
        if http_endpoint.email == from {
            // send request to http_endpoint.endpoint
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
    if b, err := ioutil.ReadAll(r); err != nil {
        return err
    } else {
        log.Println("Data:", string(b))
    }
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
            enrty_pair = strings.SplitN(pair[1], " ", 2)
	    http_endpoints = append(http_endpoints, HttpEndpoint{
		email: entry_pair[0],
	        endpoint: entry_pair[1]
        })
    }
}

log.Println("HTTP Endpoints: ", http_endpoints)
	be := &Backend{}

	s := smtp.NewServer(be)
	s.Addr = ":25"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

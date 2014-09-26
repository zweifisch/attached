package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/alexcesaro/mail/gomail"
)

type Config struct {
	Account   string `toml:"account"`
	Password  string `toml:"password"`
	Signature string `toml:"signature"`
	From      string `toml:"from"`
	Smtp      Smtp   `toml:"smtp"`
}

type Smtp struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Client struct {
	config Config
}

func (client Client) send(receiver string, subject string, message string, attachments []string) {
	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", client.config.From, client.config.Signature)
	msg.SetHeader("To", receiver)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/plain", message)
	for _, attachment := range attachments {
		if err := msg.Attach(attachment); err != nil {
			log.Fatal(err)
		}
	}

	mailer := gomail.NewMailer(client.config.Smtp.Host, client.config.Account, client.config.Password, client.config.Smtp.Port)
	if err := mailer.Send(msg); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var config Config

	usr, _ := user.Current()
	configPath := path.Join(usr.HomeDir, ".attachedrc")

	_, err := toml.DecodeFile(configPath, &config)
	if err != nil {
		log.Fatal(err)
	}

	client := Client{config}

	var message = flag.String("message", "", "the message to send")
	flag.StringVar(message, "m", "", "ailas for message")

	var to = flag.String("to", "", "where to send")
	flag.StringVar(to, "t", "", "ailas for to")

	flag.Parse()

	attachments := flag.Args()

	if len(attachments) == 0 {
		fmt.Printf("usage: %s -t receiver@somemail.com -m message attachment1 attachment2\n", os.Args[0])
		os.Exit(1)
	}

	client.send(*to, filepath.Base(attachments[0]), *message, attachments)

	fmt.Println("sent")
}

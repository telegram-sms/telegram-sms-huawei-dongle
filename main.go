package main

import (
	"bytes"
	"log"
	"regexp"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		URL: "https://api.telegram.org",
		Token:  "1076016207:AAHZSJVSXXA80AkmxUtxeje_A2PlQW7npos",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	log.Println("Configuration file loaded.")
	if err != nil {
		log.Fatal(err)
		return
	}
	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			log.Println("")
			return
		}

		b.Send(m.Sender, "[System information]\nAvailable Commands:\n/getinfo - Get system information\n/sendsms - Send SMS")
	})

	b.Handle("/sendsms",func(m *tb.Message){
		head := "[Send SMS]"
		if !m.Private() {
			return
		}
		command := m.Text
		commandList := strings.Split(command,"\n")
		log.Println(len(commandList))
		if len(commandList)<=2 {
			log.Println("Error Long")
			b.Send(m.Sender,head+"\nFail to get information.")
			return
		}
		if ! isPhoneNumber(commandList[1]){
			log.Println("This is not a legal phone number.")
			b.Send(m.Sender,head+"\nThis is not a legal phone number.")
			return
		}
		var buffer bytes.Buffer
		for i := 3; i <= len(commandList); i++{
			if i != 3 {
				buffer.WriteString("\n")
			}
			buffer.WriteString(commandList[i-1])
		}
		Content := buffer.String()
		log.Println(head+" To: "+commandList[1]+" Content: "+Content)
		b.Send(m.Sender,head+"\nTo: "+commandList[1]+"\nContent: "+Content)
	})

	b.Handle("/getinfo", func(m *tb.Message) {
		//b.Send(m.Sender, "Hello World!")
	})

	b.Start()
}
func isPhoneNumber(number string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-\.\ \\\/]?)?((?:\(?\d{1,}\)?[\-\.\ \\\/]?){0,})(?:[\-\.\ \\\/]?(?:#|ext\.?|extension|x)[\-\.\ \\\/]?(\d+))?$`)
	return re.MatchString(number)
}
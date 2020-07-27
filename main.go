package main

import (
	"bytes"
	"fmt"
	"github.com/telegram-sms/telegram-sms-huawei-dongle/client"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/json-iterator/go"

	"gopkg.in/tucnak/telebot.v2"
)

const SYSTEMHEAD = "[System Information]"

type ConfigObj struct {
	ChatId        int    `json:"chat_id"`
	BotToken      string `json:"bot_token"`
	DongleURL     string `json:"dongle_url"`
	AdminPassword string `json:"password"`
}

func main() {

	jsoniterObj := jsoniter.ConfigCompatibleWithStandardLibrary
	var SystemConfig ConfigObj
	err2 := jsoniterObj.Unmarshal(openFile("config.json"), &SystemConfig)
	if err2 != nil {
		log.Fatal(err2)
		return
	}

	log.Println("Configuration file loaded.")

	var botHandle, err = telebot.NewBot(telebot.Settings{
		URL:    "https://api.telegram.org",
		Token:  SystemConfig.BotToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
		return
	}

	adminClient := getAdminClient(SystemConfig.DongleURL, SystemConfig.AdminPassword)

	go receiveSMS(adminClient, botHandle, SystemConfig)

	botCommand(adminClient, botHandle, SystemConfig)
}

func receiveSMS(clientOBJ *client.Client, botHandle *telebot.Bot, SystemConfig ConfigObj) {
	for {
		if !checkLoginStatus(clientOBJ) {
			log.Println("logout")
			return
		}
		result, err := clientOBJ.SMSCount()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Unread: %s\n", strconv.Itoa(result.InboxUnread))
		if result.InboxUnread > 0 {
			response, err := clientOBJ.SMSList(1, 50)
			if err != nil {
				log.Fatal(err)
			}
			for _, item := range response.Messages {
				if item.Status == client.SMS_UNREAD_STATUS {
					message := fmt.Sprintf("[Receive SMS]\nFrom: %s\nContent: %s\nDate: %s", item.Phone, item.Content, item.Date)
					botHandle.Send(telebot.ChatID(SystemConfig.ChatId), message)
					messageID, _ := strconv.ParseInt(item.MessageID, 10, 64)
					clientOBJ.SetRead(messageID)
				} else {
					log.Println("The message has been read, skip it.")
				}
			}
		}
		time.Sleep(5 * time.Second)
	}

}

func botCommand(clientOBJ *client.Client, botHandle *telebot.Bot, SystemConfig ConfigObj) {
	botHandle.Handle("/start", func(m *telebot.Message) {
		log.Println("/start")
		if !checkChatState(SystemConfig.ChatId, m) {
			return
		}
		botHandle.Send(m.Sender, SYSTEMHEAD+"\nAvailable Commands:\n/getinfo - Get system information\n/sendsms - Send SMS")
	})

	botHandle.Handle("/sendsms", func(m *telebot.Message) {
		if !checkChatState(SystemConfig.ChatId, m) {
			return
		}
		if !checkLoginStatus(clientOBJ) {
			log.Println("Login status check failed")
			err := clientOBJ.UpdateSession()
			if err != nil {
				botHandle.Send(m.Sender, "Unable to update login session information.")
				log.Fatal(err)
			}
		}
		head := "[Send SMS]"
		command := m.Text
		commandList := strings.Split(command, "\n")
		log.Println(len(commandList))
		if len(commandList) <= 2 {
			log.Println("Command format error.")
			botHandle.Send(m.Sender, head+"\nFail to get information.")
			return
		}
		if !isPhoneNumber(commandList[1]) {
			log.Println("This is not a legal phone number.")
			botHandle.Send(m.Sender, head+"\nThis is not a legal phone number.")
			return
		}
		PhoneNumber := commandList[1]
		log.Println(PhoneNumber)
		var buffer bytes.Buffer
		for i := 3; i <= len(commandList); i++ {
			if i != 3 {
				buffer.WriteString("\n")
			}
			buffer.WriteString(commandList[i-1])
		}
		Content := buffer.String()
		log.Println(fmt.Sprintf("%s To: %s Content: %s", head, PhoneNumber, Content))
		botHandle.Send(m.Sender, fmt.Sprintf("%s\nTo: %s\nContent: %s", head, PhoneNumber, Content))
		_, err := clientOBJ.SendSMS(PhoneNumber, Content)
		if err != nil {
			log.Fatal(err)
		}
	})

	botHandle.Handle("/getinfo", func(m *telebot.Message) {
		if !checkChatState(SystemConfig.ChatId, m) {
			return
		}
		unavailable := "Not available"
		response := fmt.Sprintf("%s\nBattery Level: %s\nNetwork status: %s\nSIM: %s", SYSTEMHEAD, unavailable, unavailable, unavailable)
		botHandle.Send(m.Sender, response)
	})
	botHandle.Start()

}

func openFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func checkChatState(chatId int, m *telebot.Message) bool {
	if !m.Private() {
		log.Println("Request type is not allowed by security policy.")
		return false
	}
	if chatId != m.Sender.ID {
		log.Printf("Chat ID[%s] not allow.\n", chatId)
		return false
	}
	return true
}
func isPhoneNumber(number string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-. \\/]?)?((?:\(?\d+\)?[\-. \\/]?)*)(?:[\-. \\/]?(?:#|ext\.?|extension|x)[\-. \\/]?(\d+))?$`)
	return re.MatchString(number)
}

func checkLoginStatus(dongleClient *client.Client) bool {
	login, err := dongleClient.GetLoginState()
	if err != nil {
		log.Fatal(err)
		return false
	}
	if login.IsLoggedIn() {
		log.Println("Huawei dongle login detected.")
		return true
	}
	return false
}

func getAdminClient(dongleURL string, password string) *client.Client {
	log.Println("logging in...")
	c := &client.Client{}
	_ = c.Init(dongleURL)
	_, err := c.Login("admin", password)
	if err != nil {
		log.Fatal(err)
	}
	_, _ = c.GetSessionTokenInfo()
	return c
}

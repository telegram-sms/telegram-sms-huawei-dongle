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

//goland:noinspection GoSnakeCaseUsage
var G_adminClient *client.Client

type ConfigObj struct {
	ChatID        int64  `json:"chat_id"`
	BotToken      string `json:"bot_token"`
	DongleURL     string `json:"dongle_url"`
	AdminPassword string `json:"password"`
}

func main() {

	jsoniterObj := jsoniter.ConfigCompatibleWithStandardLibrary
	var SystemConfig ConfigObj
	errLoadingJson := jsoniterObj.Unmarshal(openFile("config.json"), &SystemConfig)
	if errLoadingJson != nil {
		log.Fatal(errLoadingJson)
	}

	log.Println("Configuration file loaded.")
	var botHandle, err = telebot.NewBot(telebot.Settings{
		URL:    "https://api.telegram.org",
		Token:  SystemConfig.BotToken,
		Poller: &telebot.LongPoller{Timeout: 50 * time.Second},
	})
	if err != nil {
		log.Fatal(err)
	}

	G_adminClient = getAdminClient(SystemConfig.DongleURL, SystemConfig.AdminPassword)

	login, err := G_adminClient.GetLoginState()
	if err != nil {
		log.Fatal("Login to Huawei Dongle failed. Please check the connection.")
	}
	if login.IsLoggedIn() {
		log.Println("Login to Huawei Dongle successfully.")
	}
	go receiveSMS(botHandle, SystemConfig)

	botCommand(botHandle, SystemConfig)
}

func receiveSMS(botHandle *telebot.Bot, SystemConfig ConfigObj) {
	for {
		renewAdminClient(SystemConfig)
		result, err := G_adminClient.SMSCount()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Unread: %s\n", strconv.Itoa(result.InboxUnread))
		if result.InboxUnread > 0 {
			response, err := G_adminClient.SMSList(1, 50)
			if err != nil {
				log.Println(err)
				log.Println(response)
			}
			for _, item := range response.Messages {
				if item.Status == client.SMS_UNREAD_STATUS {
					message := fmt.Sprintf("[Receive SMS]\nFrom: %s\nContent: %s\nDate: %s", item.Phone, item.Content, item.Date)
					botHandle.Send(telebot.ChatID(SystemConfig.ChatID), message, &telebot.SendOptions{DisableWebPagePreview: true})
					messageID, _ := strconv.ParseInt(item.MessageID, 10, 64)
					G_adminClient.SetRead(messageID)
				} else {
					//log.Println("The message has been read, skip it.")
				}
			}
		}
		time.Sleep(60 * time.Second)
	}

}

func botCommand(botHandle *telebot.Bot, SystemConfig ConfigObj) {
	var SMSSendInfoNextStatus = -1
	var SMSSendPhoneNumber = ""
	//goland:noinspection GoUnusedConst,GoSnakeCaseUsage
	const (
		SMS_SEND_INFO_STANDBY_STATUS       = -1
		SMS_SEND_INFO_PHONE_INPUT_STATUS   = 0
		SMS_SEND_INFO_MESSAGE_INPUT_STATUS = 1
	)

	botHandle.Handle("/start", func(m *telebot.Message) {
		SMSSendInfoNextStatus = SMS_SEND_INFO_STANDBY_STATUS
		if !checkChatState(SystemConfig.ChatID, m) {
			return
		}
		botHandle.Send(telebot.ChatID(SystemConfig.ChatID), SYSTEMHEAD+"\nAvailable Commands:\n/getinfo - Get system information\n/sendsms - Send SMS")
	})

	botHandle.Handle("/sendsms", func(m *telebot.Message) {
		SMSSendInfoNextStatus = SMS_SEND_INFO_STANDBY_STATUS
		if !checkChatState(SystemConfig.ChatID, m) {
			return
		}
		renewAdminClient(SystemConfig)
		head := "[Send SMS]\n"
		command := m.Text
		commandList := strings.Split(command, "\n")
		log.Println(len(commandList))
		if len(commandList) <= 2 {
			SMSSendInfoNextStatus = SMS_SEND_INFO_PHONE_INPUT_STATUS
			botHandle.Send(telebot.ChatID(SystemConfig.ChatID), head+"Please enter the receiver's number.")
			return
		}
		if !isPhoneNumber(commandList[1]) {
			log.Println("This is not a legal phone number.")
			botHandle.Send(telebot.ChatID(SystemConfig.ChatID), head+"This is not a legal phone number.")
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
		doSendSMS(botHandle, G_adminClient, SystemConfig.ChatID, PhoneNumber, Content)
	})

	botHandle.Handle("/getinfo", func(m *telebot.Message) {
		SMSSendInfoNextStatus = SMS_SEND_INFO_STANDBY_STATUS
		if !checkChatState(SystemConfig.ChatID, m) {
			return
		}
		renewAdminClient(SystemConfig)
		unavailable := "Not available"
		response := fmt.Sprintf("%s\nBattery Level: %s\nNetwork status: %s\nSIM: %s", SYSTEMHEAD, unavailable, unavailable, unavailable)
		botHandle.Send(m.Chat, response)
	})

	botHandle.Handle(telebot.OnText, func(m *telebot.Message) {
		log.Println(m.Text)
		head := "[Send SMS]\n"
		switch SMSSendInfoNextStatus {
		case SMS_SEND_INFO_STANDBY_STATUS:
			return
		case SMS_SEND_INFO_PHONE_INPUT_STATUS:
			if !isPhoneNumber(m.Text) {
				botHandle.Send(telebot.ChatID(SystemConfig.ChatID), head+"This phone number is invalid. Please enter it again.")
				break
			}
			SMSSendPhoneNumber = m.Text
			SMSSendInfoNextStatus = SMS_SEND_INFO_MESSAGE_INPUT_STATUS
			botHandle.Send(telebot.ChatID(SystemConfig.ChatID), head+"Please enter the message to be sent.")
			break
		case SMS_SEND_INFO_MESSAGE_INPUT_STATUS:
			doSendSMS(botHandle, G_adminClient, SystemConfig.ChatID, SMSSendPhoneNumber, m.Text)
			break
		}
		return
	})
	botHandle.Start()

}

func doSendSMS(botHandle *telebot.Bot, clientOBJ *client.Client, chatID int64, PhoneNumber string, Content string) {
	head := "[Send SMS]"
	botHandle.Send(telebot.ChatID(chatID), fmt.Sprintf("%s\nTo: %s\nContent: %s", head, PhoneNumber, Content))
	_, err := clientOBJ.SendSMS(PhoneNumber, Content)
	if err != nil {
		log.Fatal(err)
	}

}

func openFile(filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return b
}

func checkChatState(chatId int64, m *telebot.Message) bool {
	//if !m.Private() {
	//log.Println("Request type is not allowed by security policy.")
	//return false
	//}
	if chatId != m.Chat.ID {
		log.Printf("Chat ID[%s] not allow.\n", strconv.FormatInt(m.Chat.ID, 10))
		return false
	}
	return true
}
func isPhoneNumber(number string) bool {
	re := regexp.MustCompile(`^(?:(?:\(?(?:00|\+)([1-4]\d\d|[1-9]\d?)\)?)?[\-. \\/]?)?((?:\(?\d+\)?[\-. \\/]?)*)(?:[\-. \\/]?(?:#|ext\.?|extension|x)[\-. \\/]?(\d+))?$`)
	return re.MatchString(number)
}

func renewAdminClient(SystemConfig ConfigObj) {
	login, err := G_adminClient.GetLoginState()
	if err != nil {
		log.Print(err)
	}
	if !login.IsLoggedIn() {
		G_adminClient = getAdminClient(SystemConfig.DongleURL, SystemConfig.AdminPassword)
	}
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

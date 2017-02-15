package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/maddevsio/simple-config"
	"strings"
)

type Session struct {
	Phone string
	Address string
	FareId int
	OrderStarted bool
	OrderCreated bool
}

var (
	config = simple_config.NewSimpleConfig("config", "yml")
	sessions = make(map[int64]*Session)
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.GetString("bot_token"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		chatStateMachine(update, bot)
	}
}

func chatStateMachine (update tgbotapi.Update, bot *tgbotapi.BotAPI) {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Быстрый заказ такси"),
		),
		//tgbotapi.NewKeyboardButtonRow(
		//	tgbotapi.NewKeyboardButton("Заказ такси"),
		//	tgbotapi.NewKeyboardButton("Машины рядом"),
		//),
		//tgbotapi.NewKeyboardButtonRow(
		//	tgbotapi.NewKeyboardButton("Тарифы"),
		//	tgbotapi.NewKeyboardButton("Помощь"),
		//),
	)

	keyboard.OneTimeKeyboard = true

	if session := sessions[update.Message.Chat.ID]; session != nil {
		if session.OrderStarted {
			if strings.HasPrefix(update.Message.Text, "+996") {
				session.Phone = update.Message.Text
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Телефон сохранен. Теперь укажите адрес")
				bot.Send(msg)
				return
			} else if session.Phone != "" && session.Address == "" {
				session.Address = update.Message.Text
				session.OrderCreated = true
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Заказ создан!")
				bot.Send(msg)
				return
			} else if session.OrderCreated {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Машина скоро будет")
				bot.Send(msg)
				return
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Заказ не открыт. Откройте заново")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = keyboard
			bot.Send(msg)
			return
		}
	}

	if update.Message.Text == "Быстрый заказ такси" {
		sessions[update.Message.Chat.ID] = &Session{}
		sessions[update.Message.Chat.ID].OrderStarted = true
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите ваш телефон. Например: +996555112233")
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что-что?")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = keyboard
	bot.Send(msg)
	return
}
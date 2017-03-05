package main

import (
	"log"
	"gopkg.in/telegram-bot-api.v4"
	"github.com/maddevsio/simple-config"
	"github.com/maddevsio/nambataxi-telegram-bot/api"
	"fmt"
	"strings"
	"strconv"
	"github.com/maddevsio/nambataxi-telegram-bot/storage"
	"github.com/maddevsio/nambataxi-telegram-bot/chat"
)

var (
	config = simple_config.NewSimpleConfig("config", "yml")
	sessions = storage.GetAllSessions()
	nambaTaxiApi api.NambaTaxiApi
	db = storage.GetGormDB("namba-taxi-bot.db")
)

func main() {
	storage.MigrateAll(db)

	nambaTaxiApi = api.NewNambaTaxiApi(
		config.GetString("partner_id"),
		config.GetString("server_token"),
		config.GetString("url"),
		config.GetString("version"),
	)

	chat.NambaTaxiApi = nambaTaxiApi //init this for keyboards

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
	basicKeyboard := chat.GetBasicKeyboard()
	orderKeyboard := chat.GetOrderKeyboard()


	// TODO: we do not need to use all sessions here, need to change this code to sqlite quering
	if session := sessions[update.Message.Chat.ID]; session != nil {
		switch session.State {

		case storage.STATE_NEED_PHONE:
			if !strings.HasPrefix(update.Message.Text, "+996") {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Телефон должен начинаться с +996")
				bot.Send(msg)
				return
			}
			session.Phone = update.Message.Text
			session.State = storage.STATE_NEED_FARE
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Телефон сохранен. Теперь укажите тариф")
			msg.ReplyMarkup = chat.GetFaresKeyboard()
			bot.Send(msg)
			return

		case storage.STATE_NEED_FARE:
			fareId, err := chat.GetFareIdByName(update.Message.Text)
			if (err != nil) {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка! Не удалось получить тариф по имени. Попробуйте еще раз")
				msg.ReplyMarkup = basicKeyboard
				bot.Send(msg)
				return
			}
			session.FareId = fareId
			session.State = storage.STATE_NEED_ADDRESS
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите ваш адрес. Куда подать машину?")
			bot.Send(msg)
			return

		case storage.STATE_NEED_ADDRESS:
			session.Address = update.Message.Text
			orderOptions := map[string][]string{
				"phone_number": {session.Phone},
				"address":      {session.Address},
				"fare":         {strconv.Itoa(session.FareId)},
			}

			order, err := nambaTaxiApi.MakeOrder(orderOptions)
			if err != nil {
				delete(sessions, update.Message.Chat.ID)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка создания заказа. Попробуйте еще раз")
				bot.Send(msg)
				return
			}
			session.State = storage.STATE_ORDER_CREATED
			session.OrderId = order.OrderId
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Заказ создан! Номер заказа %v", order.OrderId))
			msg.ReplyMarkup = chat.GetOrderKeyboard()
			bot.Send(msg)
			return

		case storage.STATE_ORDER_CREATED:
			if update.Message.Text == "Отменить мой заказ" {
				var message string
				var keyboard = orderKeyboard

				cancel, err := nambaTaxiApi.CancelOrder(session.OrderId)
				if err != nil {
					message = "Произошла системная ошибка. Попробуйте еще раз"
					log.Printf("Error canceling order %v", err)
				}

				if cancel.Status == "200" {
					message = "Ваш заказ отменен"
					keyboard = basicKeyboard
					delete(sessions, update.Message.Chat.ID)
				}
				if cancel.Status == "400" {
					message = "Ваш заказ уже нельзя отменить, он передан водителю"
				}

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				msg.ReplyMarkup = keyboard
				bot.Send(msg)
				return
			}

			order, err := nambaTaxiApi.GetOrder(session.OrderId)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Ошибка получения заказа: %v", err))
				msg.ReplyMarkup = orderKeyboard
				bot.Send(msg)
				return
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, fmt.Sprintf("Машина скоро будет. Статус вашего заказа: %v", order.Status))
				msg.ReplyMarkup = orderKeyboard
				bot.Send(msg)
				return
			}

		default:
			delete(sessions, update.Message.Chat.ID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Заказ не открыт. Откройте заново")
			msg.ReplyToMessageID = update.Message.MessageID
			msg.ReplyMarkup = basicKeyboard
			bot.Send(msg)
			return
		}
	}

	// messages reactions while out of session scope

	if update.Message.Text == "Быстрый заказ такси" {
		sessions[update.Message.Chat.ID] = &storage.Session{}
		sessions[update.Message.Chat.ID].State = storage.STATE_NEED_PHONE
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Укажите ваш телефон. Например: +996555112233")
		bot.Send(msg)
		return
	}

	if update.Message.Text == "Тарифы" {
		fares, err := nambaTaxiApi.GetFares()
		if err != nil {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ошибка. Не удалось получить тарифы. Попробуйте еще раз")
			msg.ReplyMarkup = basicKeyboard
			bot.Send(msg)
			return
		}

		var faresText string
		for _, fare := range fares.Fare {
			faresText = faresText + fmt.Sprintf("Тариф: %v. Стоимость посадки: %.2f. Стоимость за километр: %.2f.\n\n",
				fare.Name,
				fare.Flagfall,
				fare.Cost_per_kilometer,
			)
		}

		faresText = faresText + "Для получения подробной информации посетите https://nambataxi.kg/ru/tariffs/"

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, faresText)
		msg.ReplyMarkup = basicKeyboard
		bot.Send(msg)
		return
	}

	if update.Message.Text == "Узнать статус моего заказа" {
		delete(sessions, update.Message.Chat.ID)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "К сожалению у вас нет заказа")
		msg.ReplyMarkup = basicKeyboard
		bot.Send(msg)
		return
	}

	if update.Message.Text == "/start" {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Вас приветствует бот Намба Такси для мессенджера Телеграм")
		msg.ReplyMarkup = basicKeyboard
		bot.Send(msg)
		return
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Что-что?")
	msg.ReplyToMessageID = update.Message.MessageID
	msg.ReplyMarkup = basicKeyboard
	bot.Send(msg)
	return
}
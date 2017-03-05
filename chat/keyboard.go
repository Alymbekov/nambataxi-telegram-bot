package chat

import (
	"gopkg.in/telegram-bot-api.v4"
	"fmt"
	"log"
	"errors"
	"github.com/maddevsio/nambataxi-telegram-bot/api"
)

var NambaTaxiApi api.NambaTaxiApi

func GetBasicKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Быстрый заказ такси"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Тарифы"),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func GetOrderKeyboard() tgbotapi.ReplyKeyboardMarkup {
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Узнать статус моего заказа"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Отменить мой заказ"),
		),
	)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func GetFaresKeyboard() tgbotapi.ReplyKeyboardMarkup {
	fares, err := NambaTaxiApi.GetFares()
	if err != nil {
		log.Printf("error getting fares: %v", err)
		return tgbotapi.NewReplyKeyboard()
	}

	var rows []tgbotapi.KeyboardButton
	for _, fare := range fares.Fare {
		rows = append(rows, tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton(fare.Name))...)
	}

	keyboard := tgbotapi.NewReplyKeyboard(rows)
	keyboard.OneTimeKeyboard = true
	return keyboard
}

func GetFareIdByName(fareName string) (int, error) {
	fares, err := NambaTaxiApi.GetFares()
	if err != nil {
		log.Printf("error getting fares: %v", err)
		return 0, err
	}
	for _, fare := range fares.Fare {
		if fare.Name == fareName {
			return fare.Id, nil
		}
	}
	return 0, errors.New(fmt.Sprintf("Cannot find fare with name %v", fareName))
}

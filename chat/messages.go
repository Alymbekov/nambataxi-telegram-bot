package chat

const (
	BOT_WELCOME_MESSAGE =
`Вас приветствует Telegram-бот для вызова NambaTaxi!

Заказывайте такси с моей помощью за считанные секунды:
1. Вызовите команду "Быстрый заказ такси" кнопкой снизу, либо вводом текста
2. Введите телефон
3. Выберите тариф
4. Введите адрес

Я умный (и продолжаю учиться), поэтому после первой поездки запоминаю введенные номера и адреса, и при повторном заказе вам не придется вводить их заново. В процессе обработки заказа, вы можете узнать его статус.

Приятного использования!`

	BOT_ORDER_DONE =
BOT_EMOJI_PARTY + ` Ваш заказ выполнен!
Ваш счет составил *%v сом* ` + BOT_EMOJI_MONEY + `

Спасибо, что воспользовались услугами Намба Такси. Если вдруг что-то не так, то телефон Отдела Контроля Качества к вашим услугам:
+996 (312) 97-90-60
+996 (701) 97-67-03
+996 (550) 97-60-23`
)

const (
	BOT_FARES_LINK = "Для получения более подробной информации, перейдите по ссылке: https://nambataxi.kg/ru/tariffs/"
	BOT_NO_ORDERS = "К сожалению у вас нет заказа"
	BOT_ERROR_GET_FARES = "Ошибка. Не удалось получить тарифы. Попробуйте еще раз"
	BOT_ASK_PHONE = BOT_EMOJI_PHONE + " Укажите ваш телефон. Например: +996555112233\n\n/Cancel"
	BOT_PHONE_START_996 = "Телефон должен начинаться с +996"
	BOT_ASK_FARE = BOT_EMOJI_FARE + " Телефон сохранен. Теперь укажите тариф\n\n/Cancel"
	BOT_ASK_ADDRESS = BOT_EMOJI_ADDR + " Укажите ваш адрес. Куда подать машину?\n\n/Cancel"
	BOT_ERROR_GET_1_FARE = "Ошибка! Не удалось получить тариф по имени. Попробуйте еще раз"
	BOT_ERROR_EARLY_NEAREST_DRIVERS = "Для начала нужно создать заказ"
	BOT_FARE_INFO = "Тариф: %v. Стоимость посадки: %.2f. Стоимость за километр: %.2f.\n\n"
	BOT_ORDER_NOT_CREATED = "Заказ не открыт. Откройте заново"
	BOT_ERROR_GET_ORDER = "Ошибка получения заказа: %v"
	BOT_ERROR_GET_NEAREST_DRIVERS = "Ошибка получения машин рядом: %v"
	BOT_NEAREST_DRIVERS = BOT_EMOJI_CAR + " Свободных машин рядом с вами: %v"
	BOT_ORDER_THANKS = "Спасибо за ваш заказ. Он находится в обработке. Вы можете узнать сколько рядом с вами машин нажав на кнопку 'Машины рядом'. Совсем скоро водитель возьмет ваш заказ"
	BOT_ORDER_ACCEPTED = "Ура! Ваш заказ принят ближайшим водителем!\nНомер борта: %v\nВодитель: %v\nТелефон: %v\nГосномер: %v\nМарка машины: %v"
	BOT_DRIVER_LOCATION = "Текущее местоположение водителя"
	BOT_ORDER_CANCELED_BY_OPERATOR = BOT_EMOJI_CRY + " Извините, но Ваш заказ был отклонен оператором. Возможно в вашем районе нет машин"
	BOT_ORDER_CANCELED_BY_USER = BOT_EMOJI_CRY + " Ваш заказ отменен"
	BOT_ORDER_CREATED = BOT_EMOJI_THUMUP + " Заказ создан! Номер заказа %v"
	BOT_SYSTEM_ERROR = "Произошла системная ошибка. Попробуйте еще раз"
	BOT_ORDER_CANCEL_ERROR = "Ваш заказ уже нельзя отменить, он передан водителю"
	BOT_ERROR_ORDER_CREATION = "Ошибка создания заказа. Попробуйте еще раз"

	BOT_MESSAGE_START_COMMAND = "/start"
	BOT_MESSAGE_CANCEL_COMMAND = "/Cancel"
	BOT_MESSAGE_FARES = "Тарифы"
	BOT_MESSAGE_MY_ORDER_STATUS = "Узнать статус моего заказа"
	BOT_MESSAGE_ORDER_FAST_START = "Быстрый заказ такси"
	BOT_MESSAGE_NEAREST_CARS = "Машины рядом"
	BOT_MESSAGE_CANCEL = "Отменить мой заказ"
	BOT_MESSAGE_SEND_MY_PHONE = "Отправить ваш номер телефона"
)

const (
	BOT_EMOJI_CAR = "🚗"
	BOT_EMOJI_PHONE = "📞"
	BOT_EMOJI_FARE = "📄"
	BOT_EMOJI_ADDR = "🏠"
	BOT_EMOJI_THUMUP = "👍"
	BOT_EMOJI_CRY = "😢"
	BOT_EMOJI_PARTY = "🎉"
	BOT_EMOJI_MONEY = "💰"
)

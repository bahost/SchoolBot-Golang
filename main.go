package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const token = "token"

const mainGreetingsMessage = "Привет,\nМеня зовут StudentBot, я помогу тебе устроиться в нашу школу иностранных языков"
const mainSchoolInfoMessage = "Мы - LevelUP Language School\nМы занимаемся обучением студентов англискому и тайскому языкам\nНаша школа расположена тут:\nhttps://maps.app.goo.gl/X4aR9ceH58iFPTVcA?g_st=ic"
const mainVisaInfoMessage = "А также, помогаем студентам получить студенческие визы\nПодробнее на нашем сайте:\nhttps://schoolphuket.com"
const mainManagerInfoMessage = "Наши контакты:\n Line/Whatsup/Viber/Telegram +66953617452\nEmail: Info@schoolphuket.com\nFacebook: https://www.facebook.com/LevelUpPhuket/"
const visaDetailedInfoMessage = "Срок подачи документов на ED-Visa - от 3 недель до предполагаемого получения тайской визы.\nСтуденческую визу можно получить в любом посольстве Королевства Таиланд."
const visaDetailedDocumentsMessage = "Документы, необходимые для получения студенческой визы:\nПаспорт\nФото 12 шт. (3x4 см)"
const mainNearestDateMessage = "Ближайшая свободная дата для записи:\n"
const mainRegistrationInfoMessage = "Для регистрации потребуется:\n1) Ваше ФИО\n2) Планируемая дата визита в Школу"
const mainRegistrationAskForFioMessage = "Пожалуйста, укажите ваше ФИО"
const mainRegistrationAskForDateMessage = "Укажите планируемую дата визита в школу"

const backToMainMenuMessage = "Возвращаемся"

var mainCloseDate = "2023-01-01"

var photoBytes = getYandexCloudFile("https://storage.yandexcloud.net/for-cloud-student-bot/phuket_language_school.jpg")

var photoFileBytes = tgbotapi.FileBytes{
	Name:  "schoolPicture",
	Bytes: photoBytes,
}

var greetingsKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Регистрация"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Визовые услуги"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("О нас"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Ближайшая свободная дата для записи"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Связаться с менеджером"),
	),
)

var detailedVisaInfoKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Документы"),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Главное меню"),
	),
)

func main() {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
			switch update.Message.Text {
			case "/start":
				msg.Text = mainGreetingsMessage
				msg.ReplyMarkup = greetingsKeyboard
			case "О нас":
				bot.Send(tgbotapi.NewPhoto(update.Message.Chat.ID, photoFileBytes))
				msg.Text = mainSchoolInfoMessage
				bot.Send(msg)
				msg.Text = mainVisaInfoMessage
			case "Визовые услуги":
				msg.Text = visaDetailedInfoMessage
				msg.ReplyMarkup = detailedVisaInfoKeyboard
			case "Документы":
				msg.Text = visaDetailedDocumentsMessage
			case "Главное меню":
				msg.ReplyMarkup = greetingsKeyboard
				msg.Text = backToMainMenuMessage
			case "Связаться с менеджером":
				msg.Text = mainManagerInfoMessage
			case "Ближайшая свободная дата для записи":
				msg.Text = getClosestDate(mainNearestDateMessage, mainCloseDate)
			case "Регистрация":
				msg.Text = mainRegistrationInfoMessage
				bot.Send(msg)
				msg.Text = mainRegistrationAskForFioMessage
				bot.Send(msg)
				// get fio from message func
				msg.Text = mainRegistrationAskForDateMessage
				// get date from message func
			default:
				continue
			}

			_, err := bot.Send(msg)
			if err != nil {
				panic(err)
			}
		}

	}
}

func getClosestDate(msg string, dt string) string {
	msg_slice := []string{msg, dt}
	result := strings.Join(msg_slice, "\n")
	return result
}

func getYandexCloudFile(link string) []byte {
	resp, _ := http.Get(link)
	data, _ := ioutil.ReadAll(resp.Body)
	return data
}

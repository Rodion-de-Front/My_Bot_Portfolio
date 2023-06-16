package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ResponseT struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID int `json:"update_id"`
		Message  struct {
			MessageID int `json:"message_id"`
			From      struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				LastName     string `json:"last_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Chat struct {
				ID        int    `json:"id"`
				FirstName string `json:"first_name"`
				LastName  string `json:"last_name"`
				Username  string `json:"username"`
				Type      string `json:"type"`
			} `json:"chat"`
			Date    int `json:"date"`
			Contact struct {
				PhoneNumber string `json:"phone_number"`
			} `json:"contact"`
			Text string `json:"text"`
			Data string `json:"data"`
		} `json:"message"`
	} `json:"result"`
}

type InlineButton struct {
	Ok     bool `json:"ok"`
	Result []struct {
		UpdateID      int `json:"update_id"`
		CallbackQuery struct {
			ID   string `json:"id"`
			From struct {
				ID           int    `json:"id"`
				IsBot        bool   `json:"is_bot"`
				FirstName    string `json:"first_name"`
				Username     string `json:"username"`
				LanguageCode string `json:"language_code"`
			} `json:"from"`
			Message struct {
				MessageID int `json:"message_id"`
				From      struct {
					ID        int64  `json:"id"`
					IsBot     bool   `json:"is_bot"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
				} `json:"from"`
				Chat struct {
					ID        int    `json:"id"`
					FirstName string `json:"first_name"`
					Username  string `json:"username"`
					Type      string `json:"type"`
				} `json:"chat"`
				Date        int    `json:"date"`
				Text        string `json:"text"`
				ReplyMarkup struct {
					InlineKeyboard [][]struct {
						Text         string `json:"text"`
						CallbackData string `json:"callback_data"`
					} `json:"inline_keyboard"`
				} `json:"reply_markup"`
			} `json:"message"`
			ChatInstance string `json:"chat_instance"`
			Data         string `json:"data"`
		} `json:"callback_query"`
	} `json:"result"`
}

type UserT struct {
	ID          string
	FirstName   string
	LastName    string
	RegDate     int
	PhoneNumber string
}

var host string = "https://api.telegram.org/bot"
var token string = "6037986461:AAG-K9ROPEYYXEGMBZIgs_w8Qlvm44bnOzw"

func main() {

	lastMessage := 0

	for range time.Tick(time.Second * 1) {

		//отправляем запрос к Telegram API на получение сообщений
		var url string = host + token + "/getUpdates?offset=" + strconv.Itoa(lastMessage)
		response, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
		}
		data, _ := ioutil.ReadAll(response.Body)

		//посмотреть данные
		fmt.Println(string(data))

		// var responseObj ResponseT
		//парсим данные из json
		var responseObj ResponseT
		json.Unmarshal(data, &responseObj)

		var need InlineButton
		json.Unmarshal(data, &need)
		//fmt.Println(responseObj)

		//считаем количество новых сообщений
		number := len(responseObj.Result)

		//если сообщений нет - то дальше код не выполняем
		if number < 1 {
			continue
		}

		//в цикле доставать инормацию по каждому сообщению
		for i := 0; i < number; i++ {

			text := responseObj.Result[i].Message.Text
			chatId := responseObj.Result[i].Message.From.ID
			messageTime := responseObj.Result[i].Message.Date
			firstName := responseObj.Result[i].Message.From.FirstName
			button := need.Result[i].CallbackQuery.Data
			id := need.Result[i].CallbackQuery.From.ID
			mesId := need.Result[i].CallbackQuery.Message.MessageID

			//пишем бизнес логику ----------- мозги

			//отвечаем пользователю на его сообщение
			go sendMessage(chatId, id, mesId, messageTime, text, firstName, button)

		}

		//запоминаем update_id  последнего сообщения
		lastMessage = responseObj.Result[number-1].UpdateID + 1

	}
}

func sendMessage(chatId int, id int, mesId int, messageTime int, text string, firstName string, button string) {

	if text == "/start" {

		buttons := [][]map[string]interface{}{
			{{"text": "Обо мне 👨🏻‍💻", "callback_data": "about"}},
			{{"text": "Мои работы 🎯", "callback_data": "works"}},
			{{"text": "Моё резюме 📋", "callback_data": "resume"}},
			{{"text": "Контакты 📱", "callback_data": "contacts"}},
		}

		inlineKeyboard := map[string]interface{}{
			"inline_keyboard": buttons,
		}

		inlineKeyboardJSON, _ := json.Marshal(inlineKeyboard)

		http.Get(host + token + "/deleteMessage?chat_id=" + strconv.Itoa(id) + "&message_id=" + strconv.Itoa(mesId))
		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=Здравствуйте, " + firstName + "! Рад приветствовать Вас в моём боте. Уверен, что он поможет получить ответ на интересующие вас вопросы обо мне&reply_markup=" + string(inlineKeyboardJSON))

	}

	if button == "back" {
		buttons := [][]map[string]interface{}{
			{{"text": "Обо мне 👨🏻‍💻", "callback_data": "about"}},
			{{"text": "Мои работы 🎯", "callback_data": "works"}},
			{{"text": "Моё резюме 📋", "callback_data": "resume"}},
			{{"text": "Контакты 📱", "callback_data": "contacts"}},
		}

		inlineKeyboard := map[string]interface{}{
			"inline_keyboard": buttons,
		}

		inlineKeyboardJSON, _ := json.Marshal(inlineKeyboard)

		http.Get(host + token + "/deleteMessage?chat_id=" + strconv.Itoa(id) + "&message_id=" + strconv.Itoa(mesId))
		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(id) + "&text=Что ещё вы хотите узнать?&reply_markup=" + string(inlineKeyboardJSON))
	}

	if button == "about" {

		buttons := [][]map[string]interface{}{
			{{"text": "Назад 🔙", "callback_data": "back"}},
		}

		inlineKeyboard := map[string]interface{}{
			"inline_keyboard": buttons,
		}

		inlineKeyboardJSON, _ := json.Marshal(inlineKeyboard)

		imagePath := "backend/me.jpg"
		// Создание буфера для запроса с изображением
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		// Открытие файла изображения
		file, err := os.Open(imagePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// Создание формы для файла
		fileWriter, err := bodyWriter.CreateFormFile("photo", filepath.Base(imagePath))
		if err != nil {
			log.Fatal(err)
		}

		// Копирование содержимого файла в форму
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			log.Fatal(err)
		}

		// Закрытие формы
		contentType := bodyWriter.FormDataContentType()
		bodyWriter.Close()

		// Создание URL запроса
		apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendPhoto?chat_id=%s&caption=Я Full Stack разработчик - специалист, обладающий навыками и опытом в разработке как на стороне клиента (frontend), так и на стороне сервера (backend). Я занимаюсь созданием веб-приложений и владею широким спектром технологий и инструментов. На стороне клиента я работаю с языками программирования, такими как HTML, CSS и JavaScript. На стороне сервера я занимаюсь разработкой бэкенд-логики и взаимодействием с базами данных. Я работаю с языками программирования, такими как Golang или PHP. Кроме того, я знаком с базой данных MySQL и умею эффективно работать с ней. Я также разбираюсь в системе контроля версий Git. Как Full Stack разработчик, я способен охватывать полный цикл разработки приложений - от проектирования и разработки до развертывания. Буду рад помочь воплотить ваши идеи в реальность и достичь поставленных целей в разработке&reply_markup="+string(inlineKeyboardJSON), token, strconv.Itoa(id))
		requestURL, err := url.Parse(apiURL)
		if err != nil {
			log.Fatal(err)
		}

		// Создание HTTP POST-запроса с изображением
		request, err := http.NewRequest("POST", requestURL.String(), bodyBuf)
		if err != nil {
			log.Fatal(err)
		}
		request.Header.Set("Content-Type", contentType)

		// Отправка запроса
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Fatal(err)
		}
		defer response.Body.Close()

		// Чтение ответа
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Вывод конечной ссылки запроса
		finalURL := request.URL.String()
		fmt.Println("Final URL:", finalURL)

		// Вывод ответа от сервера
		fmt.Println("Response:", string(responseData))

		http.Get(host + token + "/deleteMessage?chat_id=" + strconv.Itoa(id) + "&message_id=" + strconv.Itoa(mesId))
	}

	if button == "works" {
		buttons := [][]map[string]interface{}{
			{{"text": "Проекты 👀", "callback_data": "resume"}},
			{{"text": "Назад 🔙", "callback_data": "back"}},
		}

		inlineKeyboard := map[string]interface{}{
			"inline_keyboard": buttons,
		}

		inlineKeyboardJSON, _ := json.Marshal(inlineKeyboard)

		http.Get(host + token + "/deleteMessage?chat_id=" + strconv.Itoa(id) + "&message_id=" + strconv.Itoa(mesId))
		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(id) + "&text=Я имею большой опыт по работе с CRM системами, лендингом страниц и написанию ботов. Весь перечень моих крупных проектов Вы можете увидеть в моём резюме. &reply_markup=" + string(inlineKeyboardJSON))
	}

	if button == "resume" {
		buttons := [][]map[string]interface{}{
			{{"text": "Назад 🔙", "callback_data": "back"}},
		}

		inlineKeyboard := map[string]interface{}{
			"inline_keyboard": buttons,
		}

		inlineKeyboardJSON, _ := json.Marshal(inlineKeyboard)

		http.Get(host + token + "/deleteMessage?chat_id=" + strconv.Itoa(id) + "&message_id=" + strconv.Itoa(mesId))
		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(id) + "&text=Вот ссылка на моё подробное резюме: https://rodion-de-front.github.io/rodionka.site/&reply_markup=" + string(inlineKeyboardJSON))
	}

	if button == "contacts" {
		buttons := [][]map[string]interface{}{
			{{"text": "Назад 🔙", "callback_data": "back"}},
		}

		inlineKeyboard := map[string]interface{}{
			"inline_keyboard": buttons,
		}

		inlineKeyboardJSON, _ := json.Marshal(inlineKeyboard)

		http.Get(host + token + "/deleteMessage?chat_id=" + strconv.Itoa(id) + "&message_id=" + strconv.Itoa(mesId))
		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(id) + "&text=Свяжитесь со мной. VK: https://vk.com/fantom_uk Telegram: @rodionaka Телефон: +7 (916) 762-53-03&reply_markup=" + string(inlineKeyboardJSON))
	}

}

// func sendMessage(chatId int, messageTime int, text, firstName string) {

// 	fmt.Println(text)

// 	if text == "/start" {

// 		keys := [][]map[string]interface{}{
// 			{{"text": "Обо мне"}},
// 			{{"text": "Мои работы"}},
// 			{{"text": "Моё резюме"}},
// 			{{"text": "Контакты"}},
// 		}
// 		replyKeyboard := map[string]interface{}{
// 			"keyboard":          keys,
// 			"resize_keyboard":   true,
// 			"one_time_keyboard": true,
// 		}
// 		keyboardJson, _ := json.Marshal(replyKeyboard)
// 		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=Здравствуйте, " + firstName + "! Рад приветствовать Вас в моём боте. Уверен, что он поможет Вам получить ответ на интересующие вас воросы обо мне&reply_markup=" + string(keyboardJson))

// 	}

// 	if text == "Назад" {

// 		keys := [][]map[string]interface{}{
// 			{{"text": "Обо мне"}},
// 			{{"text": "Мои работы"}},
// 			{{"text": "Моё резюме"}},
// 			{{"text": "Контакты"}},
// 		}
// 		replyKeyboard := map[string]interface{}{
// 			"keyboard":          keys,
// 			"resize_keyboard":   true,
// 			"one_time_keyboard": true,
// 		}
// 		keyboardJson, _ := json.Marshal(replyKeyboard)
// 		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=Что вы ещё хотите узнать обо мне?&reply_markup=" + string(keyboardJson))

// 	}

// 	if text == "Обо мне" {
// 		// Если нажата любая кнопка
// 		backButton := [][]map[string]interface{}{
// 			{{"text": "Назад"}},
// 		}
// 		backKeyboard := map[string]interface{}{
// 			"keyboard":          backButton,
// 			"resize_keyboard":   true,
// 			"one_time_keyboard": true,
// 		}
// 		backKeyboardJson, _ := json.Marshal(backKeyboard)
// 		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text= Я Full Stack разработчик - специалист, обладающий навыками и опытом в разработке как на стороне клиента (frontend), так и на стороне сервера (backend). Я занимаюсь созданием веб-приложений и владею широким спектром технологий и инструментов. На стороне клиента я работаю с языками программирования, такими как HTML, CSS и JavaScript. На стороне сервера я занимаюсь разработкой бэкенд-логики и взаимодействием с базами данных. Я работаю с языками программирования, такими как Golang или PHP. Кроме того, я знаком с базой данных MySQL и умею эффективно работать с ней. Я также разбираюсь в системе контроля версий Git. Как Full Stack разработчик, я способен охватывать полный цикл разработки приложений - от проектирования и разработки до развертывания. Буду рад помочь воплотить ваши идеи в реальность и достичь поставленных целей в разработке." + "&reply_markup=" + string(backKeyboardJson))

// 	}

// 	if text == "Мои работы" {

// 		// Если нажата любая кнопка
// 		backButton := [][]map[string]interface{}{
// 			{{"text": "Назад"}},
// 		}
// 		backKeyboard := map[string]interface{}{
// 			"keyboard":          backButton,
// 			"resize_keyboard":   true,
// 			"one_time_keyboard": true,
// 		}
// 		backKeyboardJson, _ := json.Marshal(backKeyboard)
// 		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=Вот ссылка на моё подробное резюме: https://rodion-de-front.github.io/rodionka.site/" + "&reply_markup=" + string(backKeyboardJson))
// 	}

// 	if text == "Моё резюме" {
// 		// Если нажата любая кнопка
// 		backButton := [][]map[string]interface{}{
// 			{{"text": "Назад"}},
// 		}
// 		backKeyboard := map[string]interface{}{
// 			"keyboard":          backButton,
// 			"resize_keyboard":   true,
// 			"one_time_keyboard": true,
// 		}
// 		backKeyboardJson, _ := json.Marshal(backKeyboard)
// 		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=Вот ссылка на моё подробное резюме: https://rodion-de-front.github.io/rodionka.site/" + "&reply_markup=" + string(backKeyboardJson))

// 	}

// 	if text == "Контакты" {
// 		// Если нажата любая кнопка
// 		backButton := [][]map[string]interface{}{
// 			{{"text": "Назад"}},
// 		}
// 		backKeyboard := map[string]interface{}{
// 			"keyboard":          backButton,
// 			"resize_keyboard":   true,
// 			"one_time_keyboard": true,
// 		}
// 		backKeyboardJson, _ := json.Marshal(backKeyboard)
// 		http.Get(host + token + "/sendMessage?chat_id=" + strconv.Itoa(chatId) + "&text=Свяжитесь со мной. VK: https://vk.com/fantom_uk Telegram: @rodionaka Телефон: +7 (916) 762-53-03" + "&reply_markup=" + string(backKeyboardJson))

// 	}
// }

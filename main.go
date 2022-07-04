package main

import (
	"fmt"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func getMess(id int) string {
	xlsx, err := excelize.OpenFile("./stat.xlsx")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	sheet := "Лист1"

	line_number := map[int]int{
		111111: 2, // Иванов
		392676: 3, // Сидоров
		242422: 4, // Смирнов
	}

	if line_number[id] != 0 {
		cell1 := xlsx.GetCellValue(sheet, "B"+strconv.Itoa(line_number[id]))
		str1 := "Продажи: " + cell1 + "\n"
		cell2 := xlsx.GetCellValue(sheet, "C"+strconv.Itoa(line_number[id]))
		str2 := "Вежливость: " + cell2 + "\n"
		cell3 := xlsx.GetCellValue(sheet, "D"+strconv.Itoa(line_number[id]))
		str3 := "Интеллект: " + cell3
		return str1 + str2 + str3
	} else {
		return "Тебя нет в списке, чувак"
	}
}

func main() {
	botToken := os.Getenv("TOKEN")
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Smth went wrong: ", err.Error())
		}
		for _, update := range updates {
			respond(botUrl, update)
			offset = update.UpdateId + 1
		}
		//fmt.Println(updates)
	}
}

// запрос обновлений
func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

// ответ на обновления
func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	//что будем слать
	botMessage.Text = getMess(update.Message.FromUser.Id)

	// запаковываем сообщение в формат json
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	// и отсылаем
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return err
	}
	return nil
}

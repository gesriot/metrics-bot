package main

import (
	"fmt"
	"os"

	"github.com/360EntSecGroup-Skylar/excelize"

	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
)

func getMessage(id int) string {
	xlsx, err := excelize.OpenFile("./stat.xlsx")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	sheet := "Лист1"

	line_number := map[int]int{
		26078024: 2, // Иванов
		10359264: 3, // Сидоров
		30962007: 4, // Смирнов
	}

	metric := [3]string{"Продажи: ", "Вежливость: ", "Интеллект: "}

	var metric_plus_cell [3]string
	var cell [3]string
	column_letter := [3]string{"B", "C", "D"}

	if line_number[id] != 0 {
		for i := 0; i < 3; i++ {
			cell[i] = xlsx.GetCellValue(sheet, column_letter[i]+strconv.Itoa(line_number[id]))

			if cell[i] == "" {
				metric_plus_cell[i] = metric[i] + cell[i] + "–\n\n"
			} else {
				if s, err := strconv.ParseFloat(cell[i], 32); err == nil {
					cell[i] = fmt.Sprintf("%g", round(s, 2))
				}
				metric_plus_cell[i] = metric[i] + cell[i] + "\n\n"
			}
		}
		return strings.Join(metric_plus_cell[:], "")
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
	botMessage.Text = getMessage(update.Message.FromUser.Id)

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

func round(number float64, presision int) float64 {
	exp10 := math.Pow10(presision)
	return math.Round(number*exp10) / exp10
}

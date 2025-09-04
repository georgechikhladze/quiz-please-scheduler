package gameprovider

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HttpProvider struct {
}

const DOMAIN string = "https://quizplease.ru"

func New() GameProvider {
	return &HttpProvider{}
}

func (w *HttpProvider) GetGamesList() map[int][]Game {
	openGamesUrl := DOMAIN + "/schedule?QpGameSearch%5BcityId%5D=9&QpGameSearch%5Bdates%5D=&QpGameSearch%5Bstatus%5D%5B%5D=1&QpGameSearch%5Bformat%5D%5B%5D=0&QpGameSearch%5Btype%5D%5B%5D=1&QpGameSearch%5Bgame_difficulty%5D%5B%5D=all&QpGameSearch%5Bbars%5D%5B%5D=46&QpGameSearch%5Bbars%5D%5B%5D=85&QpGameSearch%5Bbars%5D%5B%5D=1563&QpGameSearch%5Bbars%5D%5B%5D=1797"
	reserveGamesUrl := DOMAIN + "/schedule?QpGameSearch%5BcityId%5D=9&QpGameSearch%5Bdates%5D=&QpGameSearch%5Bstatus%5D%5B%5D=2&QpGameSearch%5Bformat%5D%5B%5D=0&QpGameSearch%5Btype%5D%5B%5D=1&QpGameSearch%5Bgame_difficulty%5D%5B%5D=all&QpGameSearch%5Bbars%5D%5B%5D=46&QpGameSearch%5Bbars%5D%5B%5D=85&QpGameSearch%5Bbars%5D%5B%5D=1563&QpGameSearch%5Bbars%5D%5B%5D=1797"

	openGames := getGames(openGamesUrl)
	reserveGames := getGames(reserveGamesUrl)

	return map[int][]Game{
		1: openGames,
		2: reserveGames,
	}
}

func getGames(url string) []Game {
	// Загружаем HTML страницу
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Ошибка загрузки страницы %s: %v", url, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Статус код ошибки: %d %s", res.StatusCode, res.Status)
		return nil
	}

	// Парсим HTML
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Ошибка парсинга HTML: %v", err)
		return nil
	}

	var games []Game

	// Ищем элементы с классом "schedule-column"
	doc.Find(".schedule-column").Each(func(_ int, s *goquery.Selection) {
		// Извлекаем данные из различных элементов
		dateNode := s.Find(".h3").First()
		numberNode := s.Find(".h2-game-card").Last()
		placeNode := s.Find(".schedule-block-info-bar").Last()
		timeNode := s.Find(".techtext").Last()
		linkNode := s.Find(".schedule-block-head").First()

		if dateNode == nil || numberNode == nil || placeNode == nil || timeNode == nil || linkNode == nil {
			return
		}

		date := cleanText(dateNode.Text())

		// Пропускаем будние дни
		if containsWeekday(date) {
			return
		}

		number := cleanText(numberNode.Text())
		place := cleanText(placeNode.Nodes[0].FirstChild.Data)
		time := cleanText(timeNode.Text())

		// Получаем ссылку
		link, exists := s.Find(".schedule-block-head").First().Attr("href")
		if !exists {
			return
		}
		fullLink := DOMAIN + link + "#play"

		games = append(games, Game{
			Date:   date,
			Link:   fullLink,
			Number: number,
			Place:  place,
			Time:   time,
		})
	})

	return games
}

// Функция для очистки текста от лишних символов
func cleanText(text string) string {
	// Удаляем табы, новые строки и лишние пробелы
	re := regexp.MustCompile(`\t|\n|\r`)
	cleaned := re.ReplaceAllString(text, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ") // Убираем множественные пробелы
	return cleaned
}

// Функция проверки на будние дни
func containsWeekday(text string) bool {
	weekdays := []string{"понедельник", "вторник", "среда", "четверг"}
	lowerText := strings.ToLower(text)

	for _, weekday := range weekdays {
		if strings.Contains(lowerText, weekday) {
			return true
		}
	}
	return false
}

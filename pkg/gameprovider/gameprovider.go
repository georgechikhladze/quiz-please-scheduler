package gameprovider

import (
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GameProvider struct {
}

const Domain string = "https://quizplease.ru"

func NewInstance() Provider {
	return &GameProvider{}
}

func (w *GameProvider) GetGamesList() map[int][]Game {
	openGamesUrl := Domain + "/schedule?QpGameSearch%5BcityId%5D=9&QpGameSearch%5Bdates%5D=&QpGameSearch%5Bstatus%5D%5B%5D=1&QpGameSearch%5Bgame_difficulty%5D%5B%5D=2&QpGameSearch%5Bgame_difficulty%5D%5B%5D=3&QpGameSearch%5Bgame_difficulty%5D%5B%5D=4&QpGameSearch%5Btype%5D%5B%5D=1&QpGameSearch%5Bformat%5D%5B%5D=0&QpGameSearch%5Bbars%5D%5B%5D=all"
	reserveGamesUrl := Domain + "/schedule?QpGameSearch%5BcityId%5D=9&QpGameSearch%5Bdates%5D=&QpGameSearch%5Bstatus%5D%5B%5D=2&QpGameSearch%5Bgame_difficulty%5D%5B%5D=2&QpGameSearch%5Bgame_difficulty%5D%5B%5D=3&QpGameSearch%5Bgame_difficulty%5D%5B%5D=4&QpGameSearch%5Btype%5D%5B%5D=1&QpGameSearch%5Bformat%5D%5B%5D=0&QpGameSearch%5Bbars%5D%5B%5D=all"

	openGames := getGames(openGamesUrl)
	reserveGames := getGames(reserveGamesUrl)

	return map[int][]Game{
		1: openGames,
		2: reserveGames,
	}
}

func getGames(url string) []Game {
	res, err := http.Get(url)
	if err != nil {
		log.Printf("Error HTML page loading %s: %v", url, err)
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Printf("Error, statis code: %d %s", res.StatusCode, res.Status)
		return nil
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Printf("Error HTML parse: %v", err)
		return nil
	}

	var games []Game

	doc.Find(".schedule-column").Each(func(_ int, s *goquery.Selection) {
		dateNode := s.Find(".h3").First()
		numberNode := s.Find(".h2-game-card").Last()
		placeNode := s.Find(".schedule-block-info-bar").Last()
		timeNode := s.Find(".techtext").Last()
		linkNode := s.Find(".schedule-block-head").First()

		if dateNode == nil || numberNode == nil || placeNode == nil || timeNode == nil || linkNode == nil {
			return
		}

		date := cleanText(dateNode.Text())

		if containsWeekday(date) {
			return
		}

		number := cleanText(numberNode.Text())
		place := cleanText(placeNode.Nodes[0].FirstChild.Data)
		time := cleanText(timeNode.Text())

		link, exists := s.Find(".schedule-block-head").First().Attr("href")
		if !exists {
			return
		}
		fullLink := Domain + link + "#play"

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

func cleanText(text string) string {
	re := regexp.MustCompile(`\t|\n|\r`)
	cleaned := re.ReplaceAllString(text, "")
	cleaned = strings.Join(strings.Fields(cleaned), " ")
	return cleaned
}

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

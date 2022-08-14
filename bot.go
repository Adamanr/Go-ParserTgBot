package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/geziyor/geziyor"
	"github.com/geziyor/geziyor/client"
	"github.com/geziyor/geziyor/export"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io/ioutil"
	"log"
	"os"
)

type Post struct {
	text   string `json:"text"`
	author string `json:"author"`
	url    string `json:"url"`
}

func main() {
	bot, err := tgbotapi.NewBotAPI("You token")
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			if err := os.Truncate("./out.json", 0); err != nil {
				log.Println("Err")
				os.WriteFile("out.json", []byte("He"), 0755)
			}
			if update.Message.Text == "/bot" {
				geziyor.NewGeziyor(&geziyor.Options{
					StartURLs: []string{"https://habr.com/ru/hub/go/"},
					ParseFunc: quotesParse,
					Exporters: []export.Exporter{&export.JSON{}},
				}).Start()
				jsonFile, err := os.Open("out.json")
				if err != nil {
					fmt.Println(err)
				}
				defer jsonFile.Close()
				byteValue, _ := ioutil.ReadAll(jsonFile)
				var users []Post
				json.Unmarshal(byteValue[1:], &users)
				for i := 0; i < len(users); i++ {
					fmt.Println(users[i].text)
				}
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, lines[1])
				//msg.ReplyToMessageID = update.Message.MessageID
				//bot.Send(msg)
			}
		}
	}
}

func quotesParse(g *geziyor.Geziyor, r *client.Response) {
	r.HTMLDoc.Find("div.tm-article-snippet").Each(func(i int, s *goquery.Selection) {
		attr_href, _ := s.Find("a.tm-article-snippet__title-link").Attr("href")
		url := "https://habr.com" + attr_href
		g.Exports <- map[string]interface{}{

			"text":   s.Find("a.tm-article-snippet__title-link").Text(),
			"author": s.Find("a.tm-user-info__username").Text(),
			"url":    url,
		}
	})
	if href, ok := r.HTMLDoc.Find("a.tm-pagination__navigation-link > a").Attr("href"); ok {
		g.Get(r.JoinURL(href), quotesParse)
	}
}

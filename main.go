package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	urlPath = "https://od-api.oxforddictionaries.com:443/api/v1/entries/en/"
	appID   = ""
	appKey  = ""
)

type inlineQueryResult struct {
	inlineQueryResultArticle tgbotapi.InlineQueryResultArticle
}

var Error *log.Logger

func init() {
	Error = log.New(
		os.Stdout, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile|log.LstdFlags)
}

func main() {
	err := os.Setenv("token", "")
	if err != nil {
		Error.Println(err)
		return
	}

	var token string
	token = os.Getenv("token")

	if token == "" {
		log.Println("No token for authorization")
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		Error.Println(err)
		return
	}

	conf := tgbotapi.NewUpdate(0)
	conf.Timeout = 30
	updates, err := bot.GetUpdatesChan(conf)
	if err != nil {
		Error.Println(err)
		return
	}

	fmt.Fprintln(os.Stdout, "Authorized on account", bot.Self.UserName)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")

	for update := range updates {
		if update.InlineQuery == nil {
			continue
		} else if update.InlineQuery.Query == "" {
			continue
		} else if update.InlineQuery.Query == "?" {
			err := sendInformation(bot, update)
			if err != nil {
				Error.Println(err)
				return
			}

		} else {
			if ok := checkWord(update.InlineQuery.Query); !ok {
				sendError(bot, update)
			} else {
				err := sendAnswers(bot, update)
				if err != nil {
					if strings.HasSuffix(err.Error(),
						"invalid character '<'"+
							" looking for beginning of value",
					) {
						continue
					}

					Error.Println(err)
					return
				}
			}
		}
	}
}

func sendError(bot *tgbotapi.BotAPI, u tgbotapi.Update) {
	title := "Error"
	description := "Could not find word. " +
		"Please, check whether the word is correct and try later."
	article := tgbotapi.NewInlineQueryResultArticle("0", title, description)
	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: u.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{article},
	}

	_, err := bot.AnswerInlineQuery(inlineConf)
	if err != nil {
		Error.Println(err)
		return
	}
}

func sendInformation(bot *tgbotapi.BotAPI, u tgbotapi.Update) error {
	title := "Info"
	description := "itmeansbot is inline bot.\n" +
		"It can help you to find meaning of any English word " +
		"and send in your Telegram chat.\n" +
		"Type word in field and see the definition.\n" +
		"Results depends on what part of speech the word is.\n"
	article := tgbotapi.NewInlineQueryResultArticle("0", title, description)
	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: u.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       []interface{}{article},
	}

	_, err := bot.AnswerInlineQuery(inlineConf)
	if err != nil {
		Error.Println(err)
		return err
	}

	return nil
}

func sendAnswers(bot *tgbotapi.BotAPI, update tgbotapi.Update) error {
	results, err := getResultsForSending(bot, update)
	if err != nil {
		return err
	}

	inlineConf := tgbotapi.InlineConfig{
		InlineQueryID: update.InlineQuery.ID,
		IsPersonal:    true,
		CacheTime:     0,
		Results:       results,
	}

	_, err = bot.AnswerInlineQuery(inlineConf)
	if err != nil {
		return err
	}

	return nil
}

func request(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("app_id", appID)
	req.Header.Set("app_key", appKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return resp, err
	}

	return resp, nil
}

func getJSON(url string) (*OxfordDictionary, error) {
	var oxfordDictionaryStruct OxfordDictionary

	resp, err := request(url)
	if err != nil {
		Error.Println(err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	err = json.Unmarshal(data, &oxfordDictionaryStruct)
	if err != nil {
		return nil, err
	}

	return &oxfordDictionaryStruct, nil
}

func getWordInformation(entries OxfordDictionary) word {
	var wordStruct word
	var lexicalEntryStruct lexical
	var lexicalEntriesStruct []lexical
	var grammaticalFeaturesSlice []string

	for _, result := range entries.Results {
		for _, lexicalEntry := range result.LexicalEntries {
			for _, wordEntry := range lexicalEntry.Entries {

				var definitionsStruct definition
				var definitionsSlice []definition

				for _, value := range wordEntry.GrammaticalFeatures {
					grammaticalFeaturesSlice = append(
						grammaticalFeaturesSlice, value.Text+" "+value.Type)
				}

				if len(wordEntry.Senses) == 0 {
					for _, deriative := range lexicalEntry.DerivativeOf {
						var seeAnother []string
						seeAnother = append(seeAnother, "See: "+deriative.Text)
						definitionsStruct = definition{
							definition: seeAnother,
						}

						definitionsSlice = append(
							definitionsSlice, definitionsStruct)
					}

				} else {
					for _, sense := range wordEntry.Senses {
						if len(sense.Definitions) == 0 {
							definitionsStruct = definition{
								definition: sense.CrossReferenceMarkers,
							}

							definitionsSlice = append(
								definitionsSlice, definitionsStruct)
						}

						definitionsStruct = definition{
							definition: sense.Definitions,
						}

						definitionsSlice = append(
							definitionsSlice, definitionsStruct)
					}
				}

				lexicalEntryStruct = lexical{
					category:            lexicalEntry.LexicalCategory,
					grammaticalFeatures: grammaticalFeaturesSlice,
					definitions:         definitionsSlice,
				}

				lexicalEntriesStruct = append(
					lexicalEntriesStruct, lexicalEntryStruct)

				wordStruct = word{
					ID:              result.ID,
					lexicalCategory: lexicalEntriesStruct,
				}

				grammaticalFeaturesSlice = nil
			}
		}
	}

	return wordStruct
}

func getMeaning(incomingWord string) (*word, error) {
	entries, err := getJSON(urlPath + url.QueryEscape(incomingWord))
	if err != nil {
		return nil, err
	}

	wordStruct := getWordInformation(*entries)

	return &wordStruct, nil
}

func makeAnswerAndDefinition(
	category lexical, wordID string,
) (string, string) {

	definitionsAnswer := ""
	grammaticalFeaturesAnswer := ""

	for _, grammaticalFeature := range category.grammaticalFeatures {
		grammaticalFeaturesAnswer += strings.Trim(
			fmt.Sprint(grammaticalFeature), "[].") + "; "
		grammaticalFeaturesAnswer += "\n"
	}

	for _, definition := range category.definitions {
		for _, d := range definition.definition {
			definitionsAnswer += strings.Trim(fmt.Sprint(d), "[].") + "; "
			definitionsAnswer += "\n"
		}
	}

	answer := "Word: " + wordID + "\n\n" + "Part of speech: " +
		category.category + "\n\n" + "Definitions: " + "\n" +
		definitionsAnswer + "\n"

	return answer, grammaticalFeaturesAnswer
}

func getResultsForSending(
	bot *tgbotapi.BotAPI, update tgbotapi.Update,
) ([]interface{}, error) {

	var inlineQueryStructAll []inlineQueryResult

	meaning, err := getMeaning(update.InlineQuery.Query)
	if err != nil {
		sendError(bot, update)
		return nil, err
	}

	var inlineQueryStructResult inlineQueryResult

	for _, category := range meaning.lexicalCategory {
		wordID := meaning.ID
		articleID := strconv.Itoa(rand.Int())

		answer, grammaticalFeaturesAnswer := makeAnswerAndDefinition(
			category, wordID)

		article := tgbotapi.NewInlineQueryResultArticle(
			articleID, wordID, answer)
		article.Description = category.category + "; " +
			grammaticalFeaturesAnswer

		inlineQueryStructResult = inlineQueryResult{
			inlineQueryResultArticle: article,
		}

		inlineQueryStructAll = append(
			inlineQueryStructAll, inlineQueryStructResult)
	}

	results := make([]interface{}, len(meaning.lexicalCategory))
	for index, value := range inlineQueryStructAll {
		results[index] = value.inlineQueryResultArticle
	}

	return results, nil
}

func checkWord(word string) bool {
	for _, letter := range word {
		if (letter < 'a' || letter > 'z') && (letter < 'A' || letter > 'Z') {
			return false
		}
	}

	return true
}

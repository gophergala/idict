package main

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"bufio"
	"strconv"
)

type LinguaResp struct {
	ErrorMsg string   `json:"error_msg"`
	Count    uint         `json:"count_words"`
	ShowMore bool      `json:"show_more"`
	Userdict []Userdict `json:"userdict3"`
}

type Userdict struct {
	Name  string `json:"name"`
	Count uint    `json:"count"`
	Words []Word `json:"words"`
}

type Word struct {
	Id          uint            `json:"word_id"`
	Value       string          `json:"word_value"`
	Transcript  string          `json:"transcription"`
	//	Created     time.Time       `json:"created_at"`
	//	LastUpdated time.Time       `json:"last_updated_at"`
	Translates   []UserTranslate `json:"user_translates"`
	SoundUrl     string          `json:"sound_url"`
	PictureUrl   string          `json:"picture_url"`
}

type UserTranslate struct {
	Value string `json:"translate_value"`
}

const (
	linguaDictUrl  = "http://lingualeo.com/userdict/json"
	linguaLoginUrl = "http://api.lingualeo.com/api/login"
	pageCount      = 2 // restriction policy, you have 116
	httpTimeout    = 15
)

var client *gorequest.SuperAgent
var showMore = true

type Config struct {
	Email    string
	Password string
}

func readConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	config := Config{}
	configScanner := bufio.NewScanner(file)
	if configScanner.Scan() {
		config.Email = configScanner.Text()
	}
	if configScanner.Scan() {
		config.Password = configScanner.Text()
	}
	//	fmt.Printf("%v \n", config)
	return config
}

func main() {
	config := readConfig("settings.txt")
	authLeo(config.Email, config.Password)
	fmt.Printf("I'll grab %v pages", pageCount)
	for i := 4; i <= 6; i++ {
		leoAskPage(i)
	}
}

/**
 * json requested by pages, @see showMore flag
 */
func leoAskPage(page int) {
	pageNumber := strconv.Itoa(page)
	fmt.Printf("Url " + linguaDictUrl + "sortBy=date&wordType=0&filter=all&page=" + pageNumber)
	_, body, errs := client.Post(linguaDictUrl).
									Query("page=" + pageNumber).
									Query("sortBy=date").Query("wordType=0").Query("filter=all").
									End()
	if errs != nil {
		log.Fatalf("Error %v \n", errs)
		os.Exit(1)
	}
//	fmt.Printf("\nPage %v Body %v\n", pageNumber, body)
	var linguaResp LinguaResp
	json.NewDecoder(strings.NewReader(body)).Decode(&linguaResp)
	showMore = linguaResp.ShowMore
	fmt.Printf("\n ShowMore %v \n", showMore)
	userdicts := linguaResp.Userdict
	fmt.Printf("\n === %v User Dictionaries \n", len(userdicts))
	for i := 0; i < len(userdicts); i++ {
//		userdicts[i].Print()
	}
	//	fmt.Printf("Decoded %v \n", linguaResp)
}

func (d *Userdict) Print() {
	fmt.Printf("\n === Dictionary '%v' [ %v words] \n", d.Name, d.Count)
	words := d.Words
	for i := 0; i < len(words); i++ {
		words[i].Print()
	}
}

func (w *Word) Print() {
	fmt.Printf("= %v [ %v ] \n", w.Value, w.Transcript)
	//	fmt.Printf("  picture %v, sound %v \n", w.PictureUrl, w.SoundUrl)
	wordTranslates := w.Translates
	for i := 0; i < len(wordTranslates); i++ {
		fmt.Printf("  - %v \n", wordTranslates[i])
	}
}

func authLeo(email, password string) {
	client = gorequest.New().Timeout(httpTimeout*time.Second)
	resp, body, errs := client.Get(linguaLoginUrl).Query("email=" + email).Query("password=" + password).End()
	if errs != nil {
		log.Fatalf("%v \n", errs)
	}
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("Failed login %v", resp.Status)
	}
	var loginResp LoginResp
	err := json.NewDecoder(strings.NewReader(body)).Decode(&loginResp)
	if err != nil {
		log.Fatalf("Failed decode %v", body)
	}
	//	log.Printf("login parse %v \n", loginResp)
}

func linguaLeoAPI() string {
	result := fmt.Sprintf(linguaDictUrl)
	log.Println("Api " + result)
	return result
}

type LoginResp struct {
	ErrorMsg   string `json:"error_msg"`
	User       User   `json:"user"`
}

type User struct {
	Username     string `json:"nickname"`
	Id           int    `json:"user_id"`
	AutologinKey string `json:"autologin_key"`
}

func addWord(word, tword string) {
	url := "http://api.lingualeo.com/addword"
	params := fmt.Sprintf("?word=%v&tword=%v", word, tword)
	fmt.Printf("%v %v", url, params)
}

func getTranslates(word string) {
	url := "http://api.lingualeo.com/gettranslates"
	params := "?word=" + word
	fmt.Printf("%v %v", url, params)
}


package main

import (
	"database/sql"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	bot       *tgbotapi.BotAPI
	botToken  = "1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8"
	baseURL   = "https://sosialbot.herokuapp.com/"
	templates *template.Template
)
//https://api.telegram.org/bot1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8/setWebhook?url=https://sosialbot.herokuapp.com/1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8

var db *sql.DB
var err error
func telegram() {
	/*
		   heroku consoleda icra run etmak lazimdir
		  curl -F "url=https://calm-garden-87183.herokuapp.com" -F
		"certificate=@/etc/ssl/certs/bot.pem"
		https://api.telegram.org/bot1280195263:AAH3ASJo92XYnYdw8psjIoD9rfsB0eG-Zbk/setWebhook
	*/
	type Tdata struct {
		Chatid  int
		Pincode int
	}

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	//log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 0

	updates, err := bot.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		var pincode int
		pincode, _ = strconv.Atoi(update.Message.Text)
		var chatid int
		chatid = int(update.Message.Chat.ID)
		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		//tdata=Tdata{int(update.Message.Chat.ID),pincode}

		log.Printf("[%s] %s", chatid, pincode)
		//var urc int
		//urc=0

		//logFatal(err)
		//b := urc != 0
		//if b {
		//
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Siz qeydiyyatdan keçdiniz")
		//	msg.ReplyToMessageID = update.Message.MessageID
		//
		//	bot.Send(msg)
		//}
	}
}

func initTelegram() {
	var err error

	bot, err = tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Println(err)
		return
	}

	// this perhaps should be conditional on GetWebhookInfo()
	// only set webhook if it is not set properly
	url := baseURL + bot.Token
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(url))
	if err != nil {
		log.Println(err)
	}
}

func webhookHandler( /*c *gin.Context*/ w http.ResponseWriter, r *http.Request) {
	defer /*c.Request.Body.Close() */ r.Body.Close()

	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}

	var update tgbotapi.Update
	err = json.Unmarshal(bytes, &update)
	if err != nil {
		log.Println(err)
		return
	}

	// to monitor changes run: heroku logs --tail
	log.Printf("FromID: %+v  From: %+v Text: %+v\n", update.Message.Chat.ID, update.Message.From, update.Message.Text)

	//var chatid int
	//chatid := int(update.Message.Chat.ID)
	//var fio string
	//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
	//tdata=Tdata{int(update.Message.Chat.ID),pincode}

	//log.Printf("[%s] %s", chatid, pincode)
	//var urc int
	//urc=0

	//logFatal(err)
	//b := urc != 0
	//if b {



	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Salam")
	msg.ReplyToMessageID = update.Message.MessageID

	bot.Send(msg)
	//}
}

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	port := os.Getenv("PORT")

	initTelegram()
	//telegram()
	//var DB_URL = "postgres://nyrdyxoc:r4lOIZWMIoHImjb16U3u6XBQEe1Fdd7Q@queenie.db.elephantsql.com:5432/nyrdyxoc"
	//db, err := pgx.Connect(context.Background(), DB_URL)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer db.Close(context.Background())

	db,err=sql.Open("postgres","postgres://nyrdyxoc:r4lOIZWMIoHImjb16U3u6XBQEe1Fdd7Q@queenie.db.elephantsql.com:5432/nyrdyxoc")
	if err!=nil{
		panic(err.Error())
	}
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5*time.Minute)


	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/"+bot.Token, webhookHandler).Methods("POST")
	router.HandleFunc("/login", loginGetHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))

}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

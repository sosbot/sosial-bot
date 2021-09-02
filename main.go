package main

import (
	"database/sql"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	//_"github.com/go-sql-driver/mysql"
)

var (
	bot       *tgbotapi.BotAPI
	botToken  = "1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8"
	baseURL   = "https://sosialbot.herokuapp.com/"
	templates *template.Template
)

var mainMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🏠 Müraciət et"),
		tgbotapi.NewKeyboardButton("📧 Müraciətlərim")),

	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📌 Ünvanımı paylaş"), tgbotapi.NewKeyboardButton("☑ Agentlik haqda")),
	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("🔘 Rəhbərlik")),
)

var reqMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⤴Geriyə")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Haqqı ödənilən ictimai işlərə cəlb olunma")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("İşsizlikdən sığorta ödənişinin təyin edilməsi")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Peşə hazırlığına cəlb olunma")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Peşə hazırlığına cəlb olunma")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Özünüməşğulluğun təşkili")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("İşə düzəltmə")),
)

type req1 struct {
	State int
	Email string
	Fin   string
	Phone string
}

var req1Map map[int]*req1

//https://api.telegram.org/bot1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8/setWebhook?url=https://sosialbot.herokuapp.com/1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8

/*
no such file or directory
Run go mod vendor and commit the updated vendor/ directory.
Remove the vendor directory and commit the removal.


*/

var db *sql.DB
var err error
var cmdLine string

func init() {
	req1Map = make(map[int]*req1)
	cmdLine = ""
}

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
	var id int
	err = db.QueryRow("insert into public.messages(text,sent,sentby,tel_chat_id,tel_message_id) values($1,$2,$3,$4,$5) returning id;", update.Message.Text, time.Now(), update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID).Scan(&id)
	if err != nil {
		log.Println(err)
		return
	}
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

	//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Salam3")
	//msg.ReplyToMessageID = update.Message.MessageID

	//bot.Send(msg)

	//msg1 := tgbotapi.NewMessage(820987449, "From-"+update.Message.From.UserName+"_"+update.Message.From.FirstName+update.Message.From.LastName+":"+update.Message.Text)
	//msg1.ReplyToMessageID = update.Message.MessageID
	//msg1.ReplyMarkup = mainMenu
	//bot.Send(msg1)
	////u:=tgbotapi.NewUpdate(0)
	////msg,err:=bot.GetUpdatesChan(u)
	cmdText := ""

	if update.Message != nil {

		if update.Message.IsCommand() {
			cmdText = update.Message.Command()
			if cmdText == "start" {
				//message := "Xoş gəlmişsiniz!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🇦🇿 Dövlət Məşğulluq Agentliyinin telegram kanalına,xoş gəlmişsiniz!")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
			if cmdText == "stop" {
				message := "Müraciət etdiyiniz üçün, təşəkkür edirik! 🤝"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				bot.Send(msg)
			}
			if cmdText == "menu" {
				message := "Əsas səhifəyə keçid edildi"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
		} else {
			switch update.Message.Text {
			case mainMenu.Keyboard[0][0].Text:
				cmdLine = mainMenu.Keyboard[0][0].Text
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Müraciət növünü seçiniz:")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			case mainMenu.Keyboard[0][1].Text:
				cmdLine = mainMenu.Keyboard[0][1].Text
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			case mainMenu.Keyboard[1][1].Text:
				cmdLine = mainMenu.Keyboard[1][1].Text
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://dma.gov.az/agentlik/haqqimizda")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			case mainMenu.Keyboard[2][0].Text:
				cmdLine = mainMenu.Keyboard[2][0].Text
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://dma.gov.az/agentlik/idare-heyeti/idare-heyetinin-sedri/abbasbeyli-mustafa-aslan-oglu")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			case reqMenu.Keyboard[0][0].Text: //"⤴Geriyə":
				cmdLine = reqMenu.Keyboard[0][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Əsas menyuya keçid edildi")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			case mainMenu.Keyboard[1][0].Text: //"🏠 Müraciət ünvanı":
				cmdLine = mainMenu.Keyboard[1][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Zəhmət olmasa, paylaşmağa razılıq verin")
				btn := tgbotapi.KeyboardButton{
					RequestLocation: true,
					Text:            "🗺Paylaşmağa razılıq verirəm",
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
				//msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			case reqMenu.Keyboard[1][0].Text: //"Müraciət növü 1":
				cmdLine = reqMenu.Keyboard[1][0].Text
				req1Map[update.Message.From.ID] = new(req1)
				req1Map[update.Message.From.ID].State = 0
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Fin-i daxil edin:")
				msg.ReplyMarkup = tgbotapi.NewHideKeyboard(true)
				bot.Send(msg)
			case reqMenu.Keyboard[2][0].Text:
				cmdLine = reqMenu.Keyboard[2][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			case reqMenu.Keyboard[3][0].Text:
				cmdLine = reqMenu.Keyboard[3][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			case reqMenu.Keyboard[4][0].Text:
				cmdLine = reqMenu.Keyboard[4][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			case reqMenu.Keyboard[5][0].Text:
				cmdLine = reqMenu.Keyboard[5][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			case reqMenu.Keyboard[6][0].Text:
				cmdLine = reqMenu.Keyboard[6][0].Text
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			default:
				cs, ok := req1Map[update.Message.From.ID]
				if ok && cmdLine == reqMenu.Keyboard[1][0].Text {
					switch cs.State {
					case 0:
						if checkFin(update.Message.Text) == false {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Fin yanlışdır. Xahiş edirik, doğru FİN-i daxil edin:")
							bot.Send(msg)
						} else {
							cs.Fin = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Əlaqə nömrəsini daxil edin:")
							req1Map[update.Message.From.ID].State = 1
							bot.Send(msg)
						}
					case 1:
						if validPhoneFormat(update.Message.Text) == false {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Nömrə düzgün qaydada yığılmayıbdır.Misal olaraq, 9940551010101 olaraq yığılmalıdır.")
							bot.Send(msg)
						} else {
							cs.Phone = update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Email-i  daxil edin:")
							req1Map[update.Message.From.ID].State = 2
							bot.Send(msg)
						}
					case 2:
						if validEmail(update.Message.Text) == false {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Email yanlışdır.Xahiş edirik, doğru Email-i daxil edin:")
							bot.Send(msg)
						} else {
							cs.Email = update.Message.Text
							//values := req1Map[update.Message.From.ID].Phone + " " + req1Map[update.Message.From.ID].Email + " " + req1Map[update.Message.From.ID].Fin
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Müraciətiniz qəbul olundu. Müraciət nömrəsi: "+strconv.Itoa(rand.Intn(1000000)))
							msg.ReplyMarkup = mainMenu
							bot.Send(msg)
							cs.State = -1
						}

					}

				}
			}

		}

	}

	//}
}

func checkFin(value string) bool {
	if utf8.RuneCountInString(value) != 7 {
		return false
	}
	return true
}

func validPhoneFormat(value string) bool {
	re := regexp.MustCompile(`^[0-9]+$`)
	if re.MatchString(value) == true {
		return true
	} else {
		return false
	}
}

func validEmail(value string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if re.MatchString(value) == true {
		return true
	} else {
		return false
	}
}

func main() {
	templates = template.Must(template.ParseGlob("templates/*.html"))
	port := os.Getenv("PORT")

	db, err = sql.Open("postgres", "postgres://nyrdyxoc:r4lOIZWMIoHImjb16U3u6XBQEe1Fdd7Q@queenie.db.elephantsql.com:5432/nyrdyxoc")
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	initTelegram()
	//telegram()
	//var DB_URL = "postgres://nyrdyxoc:r4lOIZWMIoHImjb16U3u6XBQEe1Fdd7Q@queenie.db.elephantsql.com:5432/nyrdyxoc"
	//db, err := pgx.Connect(context.Background(), DB_URL)
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	//	os.Exit(1)
	//}
	//defer db.Close(context.Background())

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("./static/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/"+bot.Token, webhookHandler).Methods("POST")
	router.HandleFunc("/", loginGetHandler).Methods("GET")
	router.HandleFunc("/login", loginGetHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":"+port, router))

}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

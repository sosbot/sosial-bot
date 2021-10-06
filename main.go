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
		tgbotapi.NewKeyboardButton("🏠 Müraciət et."),
		tgbotapi.NewKeyboardButton("📧 Müraciətlərim")),

	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📌 Ünvanımı paylaş"), tgbotapi.NewKeyboardButton("☑ Agentlik haqda")),
	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("🔘 Rəhbərlik")),
	tgbotapi.NewKeyboardButtonRow(tgbotapi.NewKeyboardButton("📌 Məşğulluq Mərkəzləri")),
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

var branchesMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⤴Geriyə")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("🔘 Digər şəhərlər"),
		tgbotapi.NewKeyboardButton("🔘 Bakı")),
)

var capitalBranchesMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⤴Geriyə")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("3 saylı Məşğulluq mərkəzi")),
)

var regionBranchesMenu = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("⤴Geriyə")),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("Sumqayıt Məşğulluq Mərkəzi")),
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
var cmdLineMenu string
var back_clicked_once bool
var reqNumber int

func init() {
	req1Map = make(map[int]*req1)
	cmdLine = ""
	cmdLineMenu = ""
	back_clicked_once = false
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
		// if update.Message.From.ID != 820987449 {
		// 	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "texniki işlər səbəbiylə sistem müvəqqəti olaraq dayandırılıb")
		// 	bot.Send(msg)
		// 	return
		// }

		if update.Message.From.ID != 820987449 {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SosialBot-un funksionallığını daha da yaxşılaşdırmaq məqsədilə komanda olaraq, gecə-gündüz işləyirik. Hal-hazırda yeni dəyişikliklərimizi tətbiq etməyə çalışırıq. Bu səbəbdən botun funksionallığını müvəqqəti olaraq dayandırmışıq. Az sonra, son yeniliklərlə, bot fəaliyyətini davam etidərəcək. Anlayışınız üçün təşəkkür edirik.")
			bot.Send(msg)
			msg1 := tgbotapi.NewPhotoShare(update.Message.Chat.ID, `https://telegram.org.ru/uploads/posts/2018-05/1525167580_photo_2018-04-07_12-02-10.jpg`)
			bot.Send(msg1)
			return
		}
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

			if update.Message.Text == mainMenu.Keyboard[0][0].Text {
				cmdLine = mainMenu.Keyboard[0][0].Text
				cmdLineMenu = "mainMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Müraciət növünü seçiniz:")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			}
			if update.Message.Text == mainMenu.Keyboard[0][1].Text {
				cmdLine = mainMenu.Keyboard[0][1].Text
				cmdLineMenu = "mainMenu"
				rows, err := db.Query("SELECT reqnumber,reqtype,reqtext FROM public.requests WHERE reqfrom = " + strconv.Itoa(update.Message.From.ID))
				if err != nil {
					log.Println(err)
				}
				defer rows.Close()
				var reqFound bool = false
				for rows.Next() {
					reqFound = true
					var reqType string
					var reqText string
					var reqNumber int

					_ = rows.Scan(&reqNumber, &reqType, &reqText)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, reqType+"\n"+"Müraciət nömrəsi:"+strconv.Itoa(reqNumber)+"\n"+reqText+"\n"+"Status: Baxılmaqdadır")
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)

				}
				if reqFound == false {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Müraciət tapılmadı.")
					msg.ReplyMarkup = mainMenu
					bot.Send(msg)
				}

			}
			if update.Message.Text == mainMenu.Keyboard[1][1].Text {
				cmdLine = mainMenu.Keyboard[1][1].Text
				cmdLineMenu = "mainMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://dma.gov.az/agentlik/haqqimizda")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
			if update.Message.Text == mainMenu.Keyboard[2][0].Text {
				cmdLine = mainMenu.Keyboard[2][0].Text
				cmdLineMenu = "mainMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://dma.gov.az/agentlik/idare-heyeti/idare-heyetinin-sedri/abbasbeyli-mustafa-aslan-oglu")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
			if update.Message.Text == mainMenu.Keyboard[3][0].Text { //📌 Filiallar
				cmdLine = mainMenu.Keyboard[3][0].Text
				cmdLineMenu = "mainMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ərazi növünü seçiniz")
				msg.ReplyMarkup = branchesMenu
				bot.Send(msg)
			}
			if update.Message.Text == branchesMenu.Keyboard[0][0].Text && cmdLine == "" { //"⤴Geriyə":
				back_clicked_once = true
				cmdLine = branchesMenu.Keyboard[0][0].Text
				cmdLineMenu = "branchesMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Əsas səhifəyə keçid edildi")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)

			}
			if update.Message.Text == branchesMenu.Keyboard[1][0].Text { //🔘 Rayonlar üzrə
				cmdLine = branchesMenu.Keyboard[1][0].Text
				cmdLineMenu = "branchesMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Məşğulluq mərkəzini seçiniz")
				msg.ReplyMarkup = regionBranchesMenu
				bot.Send(msg)
			}
			if update.Message.Text == branchesMenu.Keyboard[1][1].Text { //🔘 Bakı üzrə
				cmdLine = branchesMenu.Keyboard[1][1].Text
				cmdLineMenu = "branchesMenu"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Məşğulluq mərkəzini seçiniz")
				msg.ReplyMarkup = capitalBranchesMenu
				bot.Send(msg)
			}
			if update.Message.Text == capitalBranchesMenu.Keyboard[0][0].Text && back_clicked_once == false && (cmdLine == branchesMenu.Keyboard[1][0].Text || cmdLineMenu == "capitalBranchesMenu") { //"⤴Geriyə":
				back_clicked_once = true
				cmdLine = capitalBranchesMenu.Keyboard[0][0].Text

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ərazi növünü seçiniz")
				msg.ReplyMarkup = branchesMenu
				cmdLineMenu = "branchesMenu"
				bot.Send(msg)

			}
			if update.Message.Text == capitalBranchesMenu.Keyboard[1][0].Text {
				cmdLine = capitalBranchesMenu.Keyboard[1][0].Text
				cmdLineMenu = "capitalBranchesMenu"
				pnt := tgbotapi.NewLocation(update.Message.Chat.ID, 40.420349239282245, 49.996552114612854)
				bot.Send(pnt)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "3 saylı Ərazi Məşğulluq Mərkəzi \n Tel:+994124525134 \n Iş saatları:  09:00–13:00,14:00–18:00")
				msg.ReplyMarkup = capitalBranchesMenu
				bot.Send(msg)

			}
			if update.Message.Text == regionBranchesMenu.Keyboard[0][0].Text && back_clicked_once == false && (cmdLine == branchesMenu.Keyboard[1][1].Text || cmdLineMenu == "regionBranchesMenu") { //"⤴Geriyə":
				back_clicked_once = true
				cmdLine = regionBranchesMenu.Keyboard[0][0].Text

				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Ərazi növünü seçiniz")
				msg.ReplyMarkup = branchesMenu
				cmdLineMenu = "branchesMenu"
				bot.Send(msg)

			}
			if update.Message.Text == regionBranchesMenu.Keyboard[1][0].Text {
				cmdLine = regionBranchesMenu.Keyboard[1][0].Text
				cmdLineMenu = "regionBranchesMenu"
				pnt := tgbotapi.NewLocation(update.Message.Chat.ID, 40.575157484113916, 49.687489343855006)
				bot.Send(pnt)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sumqayıt Məşğulluq Mərkəzi \n Tel:+994186420257 \n Iş saatları:  09:00–13:00,14:00–18:00 \n 71 Z. Hajiyev, Sumqayit 5001, Азербайджан")
				msg.ReplyMarkup = regionBranchesMenu
				bot.Send(msg)

			}
			if update.Message.Text == reqMenu.Keyboard[0][0].Text && back_clicked_once == false && (cmdLineMenu == "mainMenu" || cmdLineMenu == "reqMenu" || cmdLineMenu == "branchesMenu") { //"⤴Geriyə":
				cmdLine = reqMenu.Keyboard[0][0].Text
				cmdLineMenu = "reqMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Əsas menyuya keçid edildi")
				msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
			back_clicked_once = false
			if update.Message.Text == mainMenu.Keyboard[1][0].Text { //"🏠 Müraciət ünvanı":
				cmdLine = mainMenu.Keyboard[1][0].Text
				cmdLineMenu = "mainMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Zəhmət olmasa, paylaşmağa razılıq verin")
				btn := tgbotapi.KeyboardButton{
					RequestLocation: true,
					Text:            "🗺Paylaşmağa razılıq verirəm",
				}
				msg.ReplyMarkup = tgbotapi.NewReplyKeyboard([]tgbotapi.KeyboardButton{btn})
				//msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
			if update.Message.Text == reqMenu.Keyboard[1][0].Text { //"Müraciət növü 1":
				cmdLine = reqMenu.Keyboard[1][0].Text
				cmdLineMenu = "reqMenu"
				req1Map[update.Message.From.ID] = new(req1)
				req1Map[update.Message.From.ID].State = 999
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Fin-i daxil edin:")
				msg.ReplyMarkup = tgbotapi.NewHideKeyboard(true)
				bot.Send(msg)
			}
			if update.Message.Text == reqMenu.Keyboard[2][0].Text {
				cmdLine = reqMenu.Keyboard[2][0].Text
				cmdLineMenu = "reqMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			}
			if update.Message.Text == reqMenu.Keyboard[3][0].Text {
				cmdLine = reqMenu.Keyboard[3][0].Text
				cmdLineMenu = "reqMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			}
			if update.Message.Text == reqMenu.Keyboard[4][0].Text {
				cmdLine = reqMenu.Keyboard[4][0].Text
				cmdLineMenu = "reqMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			}
			if update.Message.Text == reqMenu.Keyboard[5][0].Text {
				cmdLine = reqMenu.Keyboard[5][0].Text
				cmdLineMenu = "reqMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			}
			if update.Message.Text == reqMenu.Keyboard[6][0].Text {
				cmdLine = reqMenu.Keyboard[6][0].Text
				cmdLineMenu = "reqMenu"
				//msg.ReplyToMessageID = update.Message.MessageID
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hörmətli Vətəndaş, Bu bölmə üzrə hal-hazırda texniki işlər aparılır. Qısa zamanda aktivləşəcək")
				msg.ReplyMarkup = reqMenu
				bot.Send(msg)
			} else {

				cs, ok := req1Map[update.Message.From.ID]

				if ok && cmdLine == reqMenu.Keyboard[1][0].Text {

					switch cs.State {
					case 0:
						if checkFin(update.Message.Text) == false {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Fin yanlışdır. Xahiş edirik, doğru FİN-i daxil edin:")
							bot.Send(msg)
						} else {
							cs.Fin = "Fin-i daxil edin:" + update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mobil nömrəni daxil edin:")
							req1Map[update.Message.From.ID].State = 1
							bot.Send(msg)
						}
					case 1:
						if validPhoneFormat(update.Message.Text) == false {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mobil nömrə yanlışdır. Düzgün qayda: +9940XXXXXXXXX")
							bot.Send(msg)
						} else {
							cs.Phone = "Mobil nömrəni daxil edin:" + update.Message.Text
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Email-i  daxil edin:")
							req1Map[update.Message.From.ID].State = 2
							bot.Send(msg)
						}
					case 2:
						if validEmail(update.Message.Text) == false {
							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Email yanlışdır.Xahiş edirik, doğru Email-i daxil edin:")
							bot.Send(msg)
						} else {
							cs.Email = "Email-i  daxil edin:" + update.Message.Text
							reqText := req1Map[update.Message.From.ID].Phone + "\n" + req1Map[update.Message.From.ID].Email + "\n" + req1Map[update.Message.From.ID].Fin
							rand.Seed(time.Now().UTC().UnixNano())
							reqNumber = rand.Intn(10000000)

							err = db.QueryRow("insert into public.requests(reqnumber,reqfrom,reqtype,reqtext) values($1,$2,$3,$4) returning reqnumber;", reqNumber, update.Message.From.ID, cmdLine, reqText).Scan(&id)
							if err != nil {
								log.Println(err)
								return
							}

							msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Müraciətiniz qəbul olundu. Müraciət nömrəsi: "+strconv.Itoa(reqNumber))
							msg.ReplyMarkup = mainMenu
							bot.Send(msg)
							cs.State = -1
						}
					case 999:
						cs.State = 0

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
	re := regexp.MustCompile("^[\\+]{1}[0-9]{3}[0]{1}[1-9]{2}[0-9]{7}$")
	if re.MatchString(value) == true {
		return true
	} else {
		return false
	}
}

func validEmail(value string) bool {
	//re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	re := regexp.MustCompile("^[a-zA-Z_\\-\\.]+[@][a-zA-Z]+[\\.][a-zA-z]{2,3}$")
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

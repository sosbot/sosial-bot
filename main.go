package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	//_"github.com/go-sql-driver/mysql"
)

const (
	LogError      = "Error"
	LogInfo       = "Info"
	LogWarning    = "Warning"
	LogAppError   = "AppError"
	LogAppInfo    = "AppInfo"
	LogAppWarning = "AppWarning"
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

type questionsArr struct {
	QuestionTypeName       string
	State                  int
	RequestText            string
	RequestErrorText       string
	ResponseValidationType string
}

type repositoryMessage struct {
	Id           string
	Text         string
	Sent         string
	SentBy       string
	TelChatId    string
	TelMessageId string
	MessageType  string
	ViewedBy     string
	ViewedAt     string
	RepltyTo     string
	VoiceText    string
	Duration     string
}

type repositoryMesages struct {
	Repos []repositoryMessage
}

type repositoryMessagesCount struct {
	Count string
}

type repositoryMessagesCountArr struct {
	Repos []repositoryMessagesCount
}

type repositoryRequestType struct {
	Id   int64
	Name string
}

type repositoryRequestTypeArr struct {
	Repos []repositoryRequestType
}

type repositoryServiceRequest struct {
	Id   int64
	Name string
}

type repositoryServiceRequestArr struct {
	Repos []repositoryServiceRequest
}

type repositoryServiceRequestToClient struct {
	TypeName    string
	ServiceName string
}
type repositoryServiceRequestToClientArr struct {
	Repos []repositoryServiceRequestToClient
}

type repositoryRequests struct {
	Name string
}
type repositoryRequestsArr struct {
	Repos []repositoryRequests
}
type MyString struct {
	val string
}

type repositoryUser struct {
	User        string
	UnviewedCnt int
}

type repositoryUsers struct {
	Repos []repositoryUser
}

type repositoryServiceRequestReq struct {
	Id int64
}

type repositoryServiceRequestReqArr struct {
	Repos []repositoryServiceRequestReq
}

type RepoRequest struct {
	ReqTypeName    string
	ReqSubTypeId   string
	ReqSubTypeName string
	ReqNumber      string
	ReqDate        string
	Status         string
}

type RepoRequestArr struct {
	Repos []RepoRequest
}

type RepoComponent struct {
	Description string
	Value       string
}

type RepoComponentArr struct {
	Repos []RepoComponent
}

var questionsArrMap map[int64]*questionsArr
var questionArrMapCurrentState int
var req1Map map[int]*req1
var CurrentState int

//https://api.telegram.org/bot1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8/setWebhook?url=https://sosialbot.herokuapp.com/1563958753:AAFNwjzp_Kvgqw0SIzHeJlxXjZnOYp2rNz8

/*
no such file or directory
Run go mod vendor and commit the updated vendor/ directory.
Remove the vendor directory and commit the removal.


*/

var db *sql.DB
var err error
var cmdLine string
var cmdLineArch string
var cmdLineMenu string
var back_clicked_once bool
var reqNumber int

var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonURL("1.com", "http://1.com"),
		tgbotapi.NewInlineKeyboardButtonData("2", "2"),
		tgbotapi.NewInlineKeyboardButtonData("3", "3"),
	),
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("4", "4"),
		tgbotapi.NewInlineKeyboardButtonData("5", "5"),
		tgbotapi.NewInlineKeyboardButtonData("6", "6"),
	),
)

func init() {
	req1Map = make(map[int]*req1)
	questionsArrMap = make(map[int64]*questionsArr)
	cmdLine = ""
	cmdLineArch = ""
	cmdLineMenu = ""
	back_clicked_once = false
	CurrentState = 0
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
	//temporary disabled for the reason callback
	// to monitor changes run: heroku logs --tail
	// log.Printf("FromID: %+v  From: %+v Text: %+v\n", update.Message.Chat.ID, update.Message.From, update.Message.Text)
	var id int64
	var text = strings.Replace(update.Message.Text, "%", "", -1)
	if len(text) > 0 {
		err = db.QueryRow("insert into public.messages(text,sent,sentby,tel_chat_id,tel_message_id,message_type) values($1,$2,$3,$4,$5,$6) returning id;", update.Message.Text, time.Now(), update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID, 1).Scan(&id)
		if err != nil {
			log.Println(err)
			return
		}
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

	if update.CallbackQuery != nil {
		logger(123, "not nil", LogAppInfo)
		logger(update.CallbackQuery.Message.Chat.ID, "Top message chat id  "+fmt.Sprint(update.CallbackQuery.Message.Chat.ID), LogAppInfo)
		logger(update.CallbackQuery.Message.Chat.ID, "Top message  id  "+fmt.Sprint(update.CallbackQuery.Message.MessageID), LogAppInfo)
		if update.CallbackQuery.Data != "nextButton" {
			execQuestionsAnswer(&update, cmdLine, update.CallbackQuery.Message.Chat.ID, CurrentState, update.CallbackQuery.Data)
		} else {
			execQuestions(cmdLine, update.CallbackQuery.Message.Chat.ID, CurrentState)
		}

		// Respond to the callback query, telling Telegram to show the user
		// a message with the data received.
		//callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)

		// And finally, send a message containing the data received.
		//msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
		//if _, err := bot.Send(msg); err != nil {
		//	panic(err)
		//}
	} else if update.Message != nil {

		//if update.Message.From.ID != 820987449 {
		//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "SosialBot-un funksionallığını daha da yaxşılaşdırmaq məqsədilə komanda olaraq, gecə-gündüz işləyirik. Hal-hazırda yeni dəyişikliklərimizi tətbiq etməyə çalışırıq. Bu səbəbdən botun funksionallığını müvəqqəti olaraq dayandırmışıq. Az sonra, son yeniliklərlə, bot fəaliyyətini davam etidərəcək. Anlayışınız üçün təşəkkür edirik.")
		//	bot.Send(msg)
		//	msg1 := tgbotapi.NewPhotoShare(update.Message.Chat.ID, `https://fins.az/file/articles/2021/04/30/1619774456_dovlet-mesgulluq-agentliyi.jpg`)
		//	bot.Send(msg1)
		//	return
		//}
		if update.Message.IsCommand() {
			cmdText = update.Message.Command()
			if cmdText == "start" {
				//message := "Xoş gəlmişsiniz!"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "🇦🇿 Dövlət Məşğulluq Agentliyinin telegram kanalına,xoş gəlmişsiniz!")
				//msg.ReplyMarkup = mainMenu
				bot.Send(msg)
			}
			if cmdText == "stop" {
				//message := "Müraciət etdiyiniz üçün, təşəkkür edirik! 🤝"
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				//msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
				//bot.Send(msg)
			}
			if cmdText == "menu" {
				//message := "Əsas səhifəyə keçid edildi"
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
				//msg.ReplyMarkup = mainMenu
				//bot.Send(msg)
			}
		} else if update.Message.Voice != nil {

			voice := *update.Message.Voice
			resp, _ := bot.GetFile(tgbotapi.FileConfig{voice.FileID})
			r, _ := http.Get("https://api.telegram.org/file/bot" + botToken + "/" + resp.FilePath)
			fmt.Println(r)
			defer r.Body.Close()
			//msg := tgbotapi.NewAudioShare(update.Message.Chat.ID, voice.FileID)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Operator səsli müraciətinizi qısa zamanda cavabalandıracaq.")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)

			vc, _ := ioutil.ReadAll(r.Body)
			err := ioutil.WriteFile("test.ogg", vc, 0644)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(vc)

			sqlStatement := `insert into voices(voice,chatid,messageid,voicesize,duration,sentdate,messages_id) values($1,$2,$3,$4,$5,$6,$7)`
			var messageid int64
			err = db.QueryRow("insert into public.messages(text,sent,sentby,tel_chat_id,tel_message_id,message_type) values($1,$2,$3,$4,$5,$6) returning id;", update.Message.Text, time.Now(), update.Message.From.ID, update.Message.Chat.ID, update.Message.MessageID, 2 /* Voice Type*/).Scan(&messageid)

			_, err = db.Exec(sqlStatement, vc, update.Message.Chat.ID, update.Message.MessageID, voice.FileSize, voice.Duration, time.Now(), messageid)

			if err != nil {
				panic(err)
			}
		} else if 1 == 2 {

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
				//rows, err := db.Query("SELECT reqnumber,reqtype,reqtext FROM public.requests WHERE reqfrom = " + strconv.Itoa(update.Message.From.ID))
				rows, err := db.Query(`SELECT name AS question_type_name,
				string_agg(answer, chr(10)) AS answer,
				CAST(request_date AS DATE) AS request_date,
				request_number,
				status
		 FROM
		   (SELECT qt."name",
				   '		*'||q.request_text||'*' || ' : ' || chr(10)||'			_'||qa.value||'_' AS answer,
		 
			  (SELECT min(qa1."timestamp"::date)
			   FROM question_answers qa1
			   WHERE qa1.request_number = qa.request_number
			   GROUP BY qa1.request_number) AS request_date,
				   qa.request_number ,
		 
			  (SELECT status
			   FROM request_statuses rs
			   WHERE id =
				   (SELECT max(id)
					FROM request_statuses rs2
					WHERE rs2.request_id =
						(SELECT id
						 FROM requests r
						 WHERE r.reqnumber = qa.request_number))) AS status
			FROM question_answers qa,
				 questions q,
				 question_type qt,
				 requests r2
			WHERE qa.questions_id = q.id
			  AND q.question_type_id = qt.id
			  and r2.reqnumber=qa.request_number
			  and r2.status=1
			  AND qa.chat_id = $1) tt
		 GROUP BY tt.name,
				  tt.request_date,
				  tt.request_number,
				  tt.status
		 ORDER BY tt.request_number ASC;	`, update.Message.Chat.ID)
				if err != nil {
					log.Println(err)
				}
				defer rows.Close()
				var reqFound bool = false
				for rows.Next() {
					reqFound = true
					var questionTypeName string
					var answer string
					var requestDate string
					var requestNumber string
					var status string
					_ = rows.Scan(&questionTypeName, &answer, &requestDate, &requestNumber, &status)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "*Sorğu nömrəsi:*"+"\n	_"+requestNumber+"_\n"+"*Müraciət mövzusu:*  "+"\n	_"+questionTypeName+"_\n"+"*Sorğu və Cavab:*"+"\n"+answer+"\n"+"*Müraciət Tarixi:*"+"\n	_"+requestDate+"_\n"+"*Müraciətin statusu:*	"+"\n_"+status+"_")
					//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Sorğu nömrəsi:"+"\n"+requestNumber+"\n"+"Müraciət mövzusu:"+"\n"+questionTypeName+"\n"+"Sorğu və Cavab:"+"\n"+answer+"\n"+"Müraciət Tarixi:"+"\n"+requestDate+"\n"+"Müraciətin statusu:"+"\n"+status)

					msg.ParseMode = "markdown"
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
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, "https://dma.gov.az/agentlik/idare-heyeti/idare-heyetinin-sedri/abbasbeyli-mustafa-aslan-oglu")
				//msg.ReplyMarkup = mainMenu
				//bot.Send(msg)
				//file := tgbotapi.File{FilePath: "https://sosialbot.eu-central-1.linodeobjects.com/audio_2022-05-07_17-49-21.ogg"}
				msg := tgbotapi.NewAudioShare(update.Message.Chat.ID, "https://sosialbot.eu-central-1.linodeobjects.com/audio_2022-05-07_17-49-21.ogg")

				msg.ReplyToMessageID = update.Message.MessageID
				_, err := bot.Send(msg)
				if err != nil {
					panic(err)
				}

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
				cmdLineArch = cmdLine
				cmdLineMenu = "reqMenu"
				req1Map[update.Message.From.ID] = new(req1)
				req1Map[update.Message.From.ID].State = 999
				//rand.Seed(time.Now().UTC().UnixNano())
				//reqNumber = rand.Intn(10000000)
				reqNumber = getNewRequestNumber()
				createNewRequest(reqNumber, update.Message.Chat.ID, cmdLine)
				setNewStatusToRequest(reqNumber, "Bot", "Gözləmədə")
				execQuestions(cmdLine, update.Message.Chat.ID, CurrentState)
				CurrentState = 999
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, execQuestions(cmdLine, update.Message.Chat.ID, CurrentState))

				//msg.ReplyMarkup = tgbotapi.NewHideKeyboard(true)
				//bot.Send(msg)
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
				if cmdLine != cmdLineArch {
					CurrentState = 0
				} else {

					if CurrentState == 999 {
						CurrentState = 1
					} else {
						execQuestionsAnswer(&update, cmdLine, update.Message.Chat.ID, CurrentState, update.Message.Text)
					}
				}
				//cs, ok := req1Map[update.Message.From.ID]
				//ok && cmdLine == reqMenu.Keyboard[1][0].Text
				if 1 == 3 {

					// switch cs.State {
					// case 0:
					// 	if checkFin(update.Message.Text) == false {
					// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Fin yanlışdır. Xahiş edirik, doğru FİN-i daxil edin:")
					// 		logger(update.Message.Chat.ID, "Fin yanlışdır. Xahiş edirik, doğru FİN-i daxil edin:", LogError)
					// 		bot.Send(msg)
					// 	} else {
					// 		cs.Fin = "Fin-i daxil edin:" + update.Message.Text
					// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mobil nömrəni(+9940XXXXXXXXX) daxil edin:")
					// 		req1Map[update.Message.From.ID].State = 1
					// 		bot.Send(msg)
					// 	}
					// case 1:
					// 	if validPhoneFormat(update.Message.Text) == false {
					// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Mobil nömrə yanlışdır. Düzgün format: +9940XXXXXXXXX")
					// 		bot.Send(msg)
					// 	} else {
					// 		cs.Phone = "Mobil nömrəni daxil edin:" + update.Message.Text
					// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Email-i  daxil edin:")
					// 		req1Map[update.Message.From.ID].State = 2
					// 		bot.Send(msg)
					// 	}
					// case 2:
					// 	if validEmail(update.Message.Text) == false {
					// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Email yanlışdır.Xahiş edirik, doğru Email-i daxil edin:")
					// 		bot.Send(msg)
					// 	} else {
					// 		cs.Email = "Email-i  daxil edin:" + update.Message.Text
					// 		reqText := req1Map[update.Message.From.ID].Phone + "\n" + req1Map[update.Message.From.ID].Email + "\n" + req1Map[update.Message.From.ID].Fin
					// 		rand.Seed(time.Now().UTC().UnixNano())
					// 		reqNumber = rand.Intn(10000000)

					// 		err = db.QueryRow("insert into public.requests(reqnumber,reqfrom,reqtype,reqtext) values($1,$2,$3,$4) returning reqnumber;", reqNumber, update.Message.From.ID, cmdLine, reqText).Scan(&id)
					// 		if err != nil {
					// 			log.Println(err)
					// 			return
					// 		}

					// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Müraciətiniz qəbul olundu. Müraciət nömrəsi: "+strconv.Itoa(reqNumber))
					// 		msg.ReplyMarkup = mainMenu
					// 		bot.Send(msg)
					// 		cs.State = -1
					// 	}
					// case 999:
					// 	cs.State = 0

					// }

				}
			}

		}

	}

	//}
}

func execQuestionsAnswer(update *tgbotapi.Update, QuestionTypeName string, chat_id int64, currentState int, answer string) {
	logger(123, QuestionTypeName, LogAppInfo)

	if reqNumber == 0 {
		return
	}

	rows, err := db.Query(`SELECT q.id,qt.name,q.state,q.request_text,q.request_error_text,coalesce(q.response_validation_type,'') as response_validation_type,q.response_type from public.questions q,public.question_type qt  where qt.id=q.question_type_id and qt.name=$1 and q.state=$2;`, QuestionTypeName, currentState)
	checkErr(err)
	defer rows.Close()
	var sequence int = 0
	var cs int
	logger(123, "ok1", LogAppInfo)
	//_, _ = questionsArrMap[chat_id]
	var questionId int
	var questionTypeName string
	var state int
	var requestText string
	var requestErrorText string
	var responseValidationType string = ""
	var responseErrorText string
	var response_type int = 0
	//var response_type_list_count int

	for rows.Next() {
		logger(123, "seq_"+strconv.Itoa(sequence), LogAppInfo)
		sequence = sequence + 1

		err = rows.Scan(&questionId, &questionTypeName, &state, &requestText, &requestErrorText, &responseValidationType, &response_type)
		checkErr(err)
	}

	switch response_type {
	case 3:
		logger(update.CallbackQuery.Message.Chat.ID, "inside message chat id  "+fmt.Sprint(update.CallbackQuery.Message.Chat.ID), LogAppInfo)
		logger(update.CallbackQuery.Message.Chat.ID, "inside message  id  "+fmt.Sprint(update.CallbackQuery.Message.MessageID), LogAppInfo)

		res, err := db.Exec(` delete from  public.question_answers where questions_id=$1 and chat_id=$2 and request_number=$3 and value=$4`, questionId, chat_id, reqNumber, answer)
		checkErr(err)
		count, err := res.RowsAffected()
		checkErr(err)
		if count == 0 {
			_, err = db.Exec(`insert into public.question_answers(questions_id,value,chat_id,request_number) values($1,$2,$3,$4);`, questionId, answer, chat_id, reqNumber)
			checkErr(err)
		}
		qlistCount := 0
		rows, err = db.Query(`SELECT count(*) as cnt  from public.question_list ql where ql.question_id=$1;`, questionId)
		checkErr(err)
		for rows.Next() {
			err = rows.Scan(&qlistCount)
		}

		rows, err = db.Query(` select ql.value,qa.value as answer_value from question_list ql  left join question_answers qa  on ql.question_id =qa.questions_id and ql.value=qa.value and qa.chat_id = $1 and qa.request_number = $2 ;`, chat_id, reqNumber)
		checkErr(err)
		defer rows.Close()
		InlineButtons := make([][]tgbotapi.InlineKeyboardButton, qlistCount)
		index := 0
		value := ""
		answerValue := ""
		for rows.Next() {
			err = rows.Scan(&value, &answerValue)
			if answerValue == value {
				InlineButtons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("☑ "+value, value))
			} else {
				InlineButtons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("☒ "+value, value))
			}
			index++
		}
		logger(123, "lenInlineButtons_"+strconv.Itoa(len(InlineButtons)), LogAppInfo)
		//msg := tgbotapi.NewMessage(chat_id, requestText)
		//msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineButtons...)
		//bot.Send(msg)

		// markup := tgbotapi.NewInlineKeyboardMarkup(
		// 	tgbotapi.NewInlineKeyboardRow(
		// 		tgbotapi.NewInlineKeyboardButtonData("❌ "+update.CallbackQuery.Data, update.CallbackQuery.Data),
		// 	))
		markup := tgbotapi.NewInlineKeyboardMarkup(InlineButtons...)
		edit := tgbotapi.NewEditMessageReplyMarkup(
			update.CallbackQuery.Message.Chat.ID,
			update.CallbackQuery.Message.MessageID,
			markup,
		)
		edit.ReplyMarkup = &markup
		bot.Send(edit)

	default:
		switch responseValidationType {
		case "FIN":
			if checkFin(answer) == false {
				responseErrorText = requestErrorText
			}
		case "MOBIL":
			if validPhoneFormat(answer) == false {
				responseErrorText = requestErrorText
			}
		case "EMAIL":
			if validEmail(answer) == false {
				responseErrorText = requestErrorText
			}
		default:
			responseErrorText = ""
		}
		logger(123, "ok2", LogAppInfo)
		logger(123, strconv.Itoa(sequence), LogAppInfo)

		if responseErrorText == "" {
			cs = currentState
			CurrentState = cs
			//_, err = db.Exec(`insert into public.question_answers(questions_id,value,chat_id,request_number) values($1,$2,$3,$4);`, questionId, answer, chat_id, reqNumber)
			//checkErr(err)
			logger(123, "responseErrorText==null", LogAppInfo)
			execQuestions(QuestionTypeName, chat_id, CurrentState)
		} else {
			logger(123, "responseErrorText not null", LogAppInfo)
			msg := tgbotapi.NewMessage(chat_id, responseErrorText)
			msg.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			bot.Send(msg)

		}
	}

}

func execQuestions(QuestionTypeName string, chat_id int64, currentState int) {
	logger(123, QuestionTypeName, LogAppInfo)
	cs := currentState + 1
	rows, err := db.Query(`SELECT  qt.name,q.state,q.request_text,q.request_error_text,q.response_type,q.id from public.questions q,public.question_type qt  where qt.id=q.question_type_id and qt.name=$1 and q.state=$2;`, QuestionTypeName, cs)
	checkErr(err)
	defer rows.Close()
	var sequence int = 0

	logger(123, "ok1", LogAppInfo)
	//_, _ = questionsArrMap[chat_id]
	var questionTypeName string
	var state int
	var requestText string
	var requestErrorText string
	//var responseValidationType string
	var response_type int = 0
	var response_type_list_count int
	var questionId int

	for rows.Next() {
		logger(123, "seq_"+strconv.Itoa(sequence), LogAppInfo)
		sequence = sequence + 1

		err = rows.Scan(&questionTypeName, &state, &requestText, &requestErrorText, &response_type, &questionId)
		checkErr(err)

		// questionsArrMap[chat_id].QuestionTypeName = questionTypeName
		// questionsArrMap[chat_id].State = state
		// questionsArrMap[chat_id].RequestText = requestText
		// questionsArrMap[chat_id].RequestErrorText = requestErrorText
		// questionsArrMap[chat_id].ResponseValidationType = responseValidationType

	}
	logger(123, "ok2", LogAppInfo)
	logger(123, strconv.Itoa(sequence), LogAppInfo)
	if sequence != 0 {
		CurrentState = cs
		logger(123, "response_type_"+strconv.Itoa(response_type), LogAppInfo)
		switch response_type {
		case 1:
			msg := tgbotapi.NewMessage(chat_id, requestText)
			msg.ReplyMarkup = tgbotapi.NewHideKeyboard(true)
			bot.Send(msg)
		case 2:
			logger(123, "question_id"+strconv.Itoa(questionId), LogAppInfo)
			rows, err = db.Query(`SELECT count(*) as cnt  from public.question_list ql where ql.question_id=$1;`, questionId)
			checkErr(err)
			for rows.Next() {
				err = rows.Scan(&response_type_list_count)
			}
			logger(123, "response_type_list_count_"+strconv.Itoa(response_type_list_count), LogAppInfo)
			defer rows.Close()
			rows, err = db.Query(`SELECT ql.value  from public.question_list ql where ql.question_id=$1;`, questionId)
			checkErr(err)
			defer rows.Close()
			InlineButtons := make([][]tgbotapi.InlineKeyboardButton, response_type_list_count)
			index := 0
			value := ""
			for rows.Next() {
				err = rows.Scan(&value)
				InlineButtons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(value, value))
				index++
			}
			logger(123, "lenInlineButtons_"+strconv.Itoa(len(InlineButtons)), LogAppInfo)
			msg := tgbotapi.NewMessage(chat_id, requestText)
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineButtons...)
			bot.Send(msg)
		case 3:
			logger(123, "question_id"+strconv.Itoa(questionId), LogAppInfo)
			rows, err = db.Query(`SELECT count(*) as cnt  from public.question_list ql where ql.question_id=$1;`, questionId)
			checkErr(err)
			for rows.Next() {
				err = rows.Scan(&response_type_list_count)
			}
			logger(123, "response_type_list_count_"+strconv.Itoa(response_type_list_count), LogAppInfo)
			defer rows.Close()
			rows, err = db.Query(`SELECT ql.value  from public.question_list ql where ql.question_id=$1;`, questionId)
			checkErr(err)
			defer rows.Close()
			InlineButtons := make([][]tgbotapi.InlineKeyboardButton, response_type_list_count)
			index := 0
			value := ""
			for rows.Next() {
				err = rows.Scan(&value)
				InlineButtons[index] = tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("☒ "+value, value))
				index++
			}
			logger(123, "lenInlineButtons_"+strconv.Itoa(len(InlineButtons)), LogAppInfo)
			msg := tgbotapi.NewMessage(chat_id, requestText)
			msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(InlineButtons...)
			bot.Send(msg)

			var nextButton = tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("▶️ Növbəti", "nextButton"),
				))
			msgNextButton := tgbotapi.NewMessage(chat_id, "Seçim(lər)i edib, növbəti düyməsinə sıxın.")
			msgNextButton.ReplyMarkup = nextButton
			bot.Send(msgNextButton)
		default:
		}
	} else {
		//rand.Seed(time.Now().UTC().UnixNano())
		//reqNumber = rand.Intn(10000000)
		if reqNumber > 0 {
			msg := tgbotapi.NewMessage(chat_id, "Müraciətiniz qəbul olundu. Müraciət nömrəsi: "+strconv.Itoa(reqNumber))
			msg.ReplyMarkup = mainMenu
			bot.Send(msg)
			_, err = db.Exec(`update  public.requests set status=1 where reqnumber=$1;`, reqNumber)
			checkErr(err)
			reqNumber = 0
			//_,err = db.Exec(`insert into public.requests(reqnumber,reqfrom,reqtype,reqtext) values($1,$2,$3,$4);`, reqNumber, chat_id, cmdLine, reqText)
			//checkErr(err)
		}
	}

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func logger(chatid int64, text string, logType string) {
	_, err := db.Exec("insert into public.logs(chat_id,text,type) values($1,$2,$3)", chatid, text, logType)
	checkErr(err)

}

func createNewRequest(request_number int, chat_id int64, request_type string) {
	_, err := db.Exec("insert into public.requests(reqnumber,reqfrom,reqtype) values($1,$2,$3)", request_number, chat_id, request_type)
	checkErr(err)
}

func getNewRequestNumber() int {
	rows, err := db.Query(`select coalesce (max(reqnumber)+1,1) as request_number from public.requests r ;`)
	checkErr(err)
	defer rows.Close()
	var maxRequestNumber int
	for rows.Next() {
		err = rows.Scan(&maxRequestNumber)
	}
	return maxRequestNumber
}

func getRequestNumberId(request_number int) int {
	rows, err := db.Query(`select  id from public.requests r  where reqnumber=$1;`, request_number)
	checkErr(err)
	defer rows.Close()
	var requestNumberId int
	for rows.Next() {
		err = rows.Scan(&requestNumberId)
	}
	return requestNumberId
}

func setNewStatusToRequest(request_number int, chat_id string, status string) {
	_, err := db.Exec("insert into public.request_statuses(request_id,status,insertedby) values($1,$2,$3)", getRequestNumberId(request_number), status, chat_id)
	checkErr(err)
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
	//db, err = sql.Open("postgres", "postgres://bbuitmkqevrfzf:ebcea06a881aee891ebaa6176c11aaa534fc7021091792389315020cf67ec954@ec2-3-209-65-193.compute-1.amazonaws.com:5432/dave9oelfnd2v0")

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
	router.HandleFunc("/messages", messagesGetHandler).Methods("GET")
	router.HandleFunc("/messages/{id}", messagesIdGetHandler).Methods("GET")
	router.HandleFunc("/messagesCount/{id}", messagesCountGetHandler).Methods("GET")
	router.HandleFunc("/users", usersGetHandler).Methods("GET")
	router.HandleFunc("/messageTo/{id}", messageToGetHandler).Methods("GET")
	router.HandleFunc("/requestTypes", requestTypesGetHandler).Methods("GET")
	router.HandleFunc("/servicesRequests/{reqtypeid}", servicesRequestsGetHandler).Methods("GET")
	router.HandleFunc("/servicesrequeststoclient", servicesRequestsToClientGetHandler).Methods("GET")
	router.HandleFunc("/userRequests/{reqnumber}", userRequestsGetHandler).Methods("GET")
	router.HandleFunc("/userRequestSave", userRequestSaveGetHandler).Methods("POST")
	router.HandleFunc("/servicecRequestsRegs", serviceRequestsReqsGetHandler).Methods("GET")
	router.HandleFunc("/Requests", requestsGetHandler).Methods("GET")
	router.HandleFunc("/Requests/{reqnumber}", requestsIdGetHandler).Methods("GET")
	router.HandleFunc("RequestsDone/{reqnumber}", requestsDoneGetHandler).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+port, router))

}

func requestsDoneGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqnumber := params["reqnumber"]
	var telegramid int64
	err := db.QueryRow(`update requests set status=2
 where r.reqnumber=$1 and status=1 returning reqfrom `, reqnumber).Scan(&telegramid)
	if err != nil {
		panic(err)
	}

	msg := tgbotapi.NewMessage(telegramid, `Hörmətli Vətəndaş, `+reqnumber+` saylı müraciətiniz sonlandırıldı. Müraciət etdiyiniz üçün Sizə təşəkkür edirik!`)
	//msg.ReplyMarkup = mainMenu
	bot.Send(msg)

	if err != nil {
		panic(err)
	}

	templates.ExecuteTemplate(w, "requestdone.html", nil)
}

func requestsIdGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqnumber := params["reqnumber"]

	var data RepoComponent
	datas := []RepoComponent{}
	rows, err := db.Query(`select src.component_description,
       coalesce(srcd.data_value,'Məlumat yoxdur')
       from servicerequestscomponents src
     join servicerequestscomponentsdatas srcd on src.id=srcd.servicerequestscomponents_id
     join requests r on r.id=srcd.requests_id
 where r.reqnumber=$1`, reqnumber)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data.Description,
			&data.Value)
		if err != nil {
			panic(err)
		}
		datas = append(datas, data)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	templates.ExecuteTemplate(w, "requestdetails.html", datas)
}

func requestsGetHandler(w http.ResponseWriter, r *http.Request) {

	var data RepoRequest
	datas := []RepoRequest{}
	rows, err := db.Query(`
select  rt.name as req_type_name,
        s.id,
        s.service_name as req_subtype_name,
        coalesce(cast(r.reqnumber as  varchar),''),
        r.datetime as reqdate,
        case r.status when 1 then 'Açıq' when 2 then 'Bağlı' end as status

        from request_type rt
   join servicesrequests s on rt.id = s.request_type_id
   join requests r on r.servicesrequestsid=s.id
   where r.status>0`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(&data.ReqTypeName,
			&data.ReqSubTypeId,
			&data.ReqSubTypeName,
			&data.ReqNumber,
			&data.ReqDate,
			&data.Status)
		if err != nil {
			panic(err)
		}
		datas = append(datas, data)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}

	templates.ExecuteTemplate(w, "requests.html", datas)
}

func serviceRequestsReqsGetHandler(w http.ResponseWriter, r *http.Request) {
	reqfrom, _ := strconv.ParseInt(r.URL.Query().Get("reqfrom"), 10, 64)
	servreqid, _ := strconv.ParseInt(r.URL.Query().Get("servicereqid"), 10, 64)
	fmt.Println(reqfrom)
	fmt.Println(servreqid)

	var id int64
	var reqnumber string
	err := db.QueryRow("insert into requests(reqfrom,servicesrequestsid,status) values($1,$2,0) returning id", reqfrom, servreqid).Scan(&id)
	fmt.Println(id)
	if err != nil {
		panic(err)
	}
	_, err = db.Exec(`update requests set reqnumber=luhn_generate($1)::numeric where id=$2 returning reqnumber::numeric`, id, id)
	if err != nil {
		panic(err)
	}
	rows, err := db.Query(`select reqnumber from requests where id=$1`, id)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&reqnumber)
	}
	fmt.Println(reqnumber)
	txt := `Hörmətli Vətəndaş, Müraciətiniz üzrə sorğunu tamamlamaq üçün xahiş edirik, ilkin tələb olunan məlumatları "Linkə keçid" vasitəsilə keçid edərək, əlavə ediniz.`
	snt := time.Now()
	_, err = db.Exec(`insert into messages(text,sent,sentby,tel_chat_id,message_type) 
                           values($1,$2,$3,$4,$5)`, txt, snt, 1, reqfrom, 1)

	msg := tgbotapi.NewMessage(reqfrom, txt)
	//msg.ReplyMarkup = mainMenu
	//bot.Send(msg)

	var numericKeyboard = tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonURL("Linkə Keçid", "https://sosialbot.herokuapp.com/userRequests/"+reqnumber)))

	//markup := tgbotapi.NewInlineKeyboardMarkup(InlineButtons...)
	msg.ReplyMarkup = &numericKeyboard
	bot.Send(msg)
	//repo := repositoryServiceRequestReqArr{}
	//err := queryServiceRequestReq(&repo, reqfrom, servreqid)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//out, err := json.Marshal(repo)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	//fmt.Fprintf(w, string(out))
}

func messageToGetHandler(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	telegramId, _ := strconv.ParseInt(params["id"], 10, 64)
	msg := tgbotapi.NewMessage(telegramId, r.URL.Query().Get("message"))
	//msg.ReplyMarkup = mainMenu
	bot.Send(msg)

	_, err = db.Exec(`insert into messages(text,sent,sentby,tel_chat_id,message_type,viewedby,viewedat) values($1,$2,$3,$4,$5,$6,$7)`, r.URL.Query().Get("message"), time.Now(), 1, params["id"], 1, 1, time.Now())
	checkErr(err)
	fmt.Fprintf(w, "")
}

func servicesRequestsGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	reqtypeid, _ := strconv.ParseInt(params["reqtypeid"], 10, 64)

	repos := repositoryServiceRequestArr{}
	err := queryServiceRequests(&repos, reqtypeid)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func servicesRequestsToClientGetHandler(w http.ResponseWriter, r *http.Request) {
	repos := repositoryServiceRequestToClientArr{}
	err := queryServiceRequestToClient(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func requestTypesGetHandler(w http.ResponseWriter, r *http.Request) {
	repos := repositoryRequestTypeArr{}
	err := queryRequestTypes(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func usersGetHandler(w http.ResponseWriter, r *http.Request) {
	repos := repositoryUsers{}
	err := queryUserRepos(&repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))

}

func messagesIdGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	repos := repositoryMesages{}
	err := queryMessageReposById(&repos, params["id"])

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	_, err = db.Exec(`update messages set viewedBy=1,viewedAt=$1 where tel_chat_id=$2`, time.Now(), params["id"])
	checkErr(err)

	fmt.Fprintf(w, string(out))

}

func loginGetHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "login.html", nil)
}

func messagesGetHandler(w http.ResponseWriter, r *http.Request) {
	repos := repositoryMesages{}
	err := queryMessageRepos(&repos)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func messagesCountGetHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	repos := repositoryMessagesCountArr{}
	err := queryMessagesCountReposById(&repos, params["id"])

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	out, err := json.Marshal(repos)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Fprintf(w, string(out))
}

func userRequestSaveGetHandler(w http.ResponseWriter, r *http.Request) {

	reqnumber := r.FormValue("reqnumber")
	var requestId int64
	var reqfrom int64
	err = db.QueryRow(`select id,reqfrom from requests where reqnumber=$1`, reqnumber).Scan(&requestId, &reqfrom)

	if err != nil {
		panic(err)
	}
	var rows, err = db.Query(`select coalesce(s.service_name,'') as service_name,
												   coalesce(cast(s2.order_num as varchar),'') as order_num,
												   coalesce(s2.component_description,'') as component_description,
												   coalesce(s2.component_type,'') as component_type,
												   coalesce(cast(s2.data_driven as varchar),'') as data_driven,
		       									   s2.id,
												   coalesce(s3.component_id,'') as component_id,
												   coalesce(s3.component_name,'') as component_name,
												   coalesce(s3.component_value,'') as component_value,
												   coalesce(s3.component_label,'') as component_label,
												   coalesce(s3.component_requiredsize,'') as component_requiredsize,
												   coalesce(s3.component_placeholder,'') as component_placeholder,
												   coalesce(cast(s3.component_minlength as varchar),'') as component_minlegth,
												   coalesce(cast(s3.component_maxlength as varchar),'') as component_maxlength,
												   coalesce(s3.component_title,'') as componet_title,
												   coalesce(s3.component_mindate,'') as component_mindate,
												   coalesce(s3.component_maxdate,'') as component_madate

									 from  servicesrequests  s
												left join servicerequestscomponents s2  on s.id=s2.services_requests_id
												left join servicerequestscomponentsdetails s3  on s2.id=s3.servicerequestscomponents_id

									where s2.data_driven=1
									 order by s2.order_num`)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		repo := RepoData{}
		err = rows.Scan(&repo.service_name,
			&repo.order_num,
			&repo.component_description,
			&repo.component_type,
			&repo.data_driven,
			&repo.componentId,
			&repo.component_id,
			&repo.component_name,
			&repo.component_value,
			&repo.component_label,
			&repo.component_requiredsize,
			&repo.component_placeholder,
			&repo.component_minlength,
			&repo.component_maxlength,
			&repo.component_title,
			&repo.component_mindate,
			&repo.component_maxdate)
		if err != nil {
			panic(err)
		}
		_, err = db.Exec(`insert into servicerequestscomponentsdatas(servicerequestscomponents_id,data_value,requests_id) values($1,$2,$3)`, repo.componentId, r.FormValue(repo.component_name), requestId)
		if err != nil {
			panic(err)
		}

	}

	_, err = db.Exec(`update requests set status=1 where id=$1 and status=0`, requestId)
	if err != nil {
		panic(err)
	}

	tmpl, _ := template.ParseFiles("templates/done.html")
	tmpl.Execute(w, "")

	txt := `Hörmətli Vətəndaş, ` + reqnumber + ` nömrəli müraciətiniz qəbul olundu. Müraciətlə bağlı Operatorla əlaqə saxladıqda, bu qeydiyyat nömrəsini təqdim etməyinizi xahiş edirik. Qısa zamanda Operator Sizinlə əlaqə saxlayacaqdır.`

	msg := tgbotapi.NewMessage(reqfrom, txt)
	//msg.ReplyToMessageID = update.Message.MessageID
	//msg.ReplyMarkup = mainMenu
	bot.Send(msg)
}

func userRequestsGetHandler(w http.ResponseWriter, r *http.Request) {
	var form inputForm
	var field inputField
	var fields []string
	var fieldSelect selectField
	var fieldCheckbox checkboxField

	var h string
	//data := make(map[string]interface{})

	params := mux.Vars(r)
	reqnumber, _ := strconv.ParseInt(params["reqnumber"], 10, 64)
	fmt.Println(reqNumber)
	var rows, err = db.Query(`select coalesce(s.service_name,'') as service_name,
										   coalesce(cast(s2.order_num as varchar),'') as order_num,
										   coalesce(s2.component_description,'') as component_description,
										   coalesce(s2.component_type,'') as component_type,
										   coalesce(cast(s2.data_driven as varchar),'') as data_driven,
										   coalesce(s3.component_id,'') as component_id,
										   coalesce(s3.component_name,'') as component_name,
										   coalesce(s3.component_value,'') as component_value,
										   coalesce(s3.component_label,'') as component_value,
										   coalesce(s3.component_requiredsize,'') as component_requiredsize,
										   coalesce(s3.component_placeholder,'') as component_placeholder,
										   coalesce(cast(s3.component_minlength as varchar),'') as component_minlegth,
										   coalesce(cast(s3.component_maxlength as varchar),'') as component_maxlength,
										   coalesce(s3.component_title,'') as componet_title,
										   coalesce(s3.component_mindate,'') as component_mindate,
										   coalesce(s3.component_maxdate,'') as component_madate
											
							 from  servicesrequests  s
										left join servicerequestscomponents s2  on s.id=s2.services_requests_id
										left join servicerequestscomponentsdetails s3  on s2.id=s3.servicerequestscomponents_id
										join requests r on  r.servicesrequestsid=s.id
        						  where r.reqnumber=$1
							 order by s2.order_num`, reqnumber)
	if err != nil {
		panic(err)
	}

	defer rows.Close()
	for rows.Next() {
		repo := RepoData{}
		err = rows.Scan(
			&repo.service_name,
			&repo.order_num,
			&repo.component_description,
			&repo.component_type,
			&repo.data_driven,
			&repo.component_id,
			&repo.component_name,
			&repo.component_value,
			&repo.component_label,
			&repo.component_requiredsize,
			&repo.component_placeholder,
			&repo.component_minlength,
			&repo.component_maxlength,
			&repo.component_title,
			&repo.component_mindate,
			&repo.component_maxdate,
		)
		if err != nil {
			panic(err)
		}
		if repo.component_type == "input_text" {
			field.template = tplInputTemplate
			field.Id = repo.component_id
			field.Label = repo.component_label
			field.Name = repo.component_name

			field.Placeholder = repo.component_placeholder
			field.Pattern = ""
			field.ReqSize = repo.component_requiredsize
			field.MaxLength = repo.component_maxlength
			field.MinLength = repo.component_minlength
			field.ErrMsg = "Error number 5"

			h = field.appendText()
			fields = append(fields, h)
		}
		if repo.component_type == "select" {
			fieldSelect.Template = tplSelectTemplate
			fieldSelect.Id = repo.component_id
			fieldSelect.Label = repo.component_label
			fieldSelect.Name = repo.component_name

			split := strings.Split(repo.component_value, ",")
			var str = ""
			for _, line := range split {
				str = str + fmt.Sprintf("<option value=\"%s\">%s</option>", line, line)
			}
			fieldSelect.Values = str
			//fmt.Printf(fieldSelect.Values)
			h = fieldSelect.appendSelect()
			fields = append(fields, h)
		}
		if repo.component_type == "inputCheckbox" {
			fieldCheckbox.Template = tplCheckboxTemplate
			fieldCheckbox.Id = repo.component_id
			fieldCheckbox.Label = repo.component_label
			fieldCheckbox.Name = repo.component_name
			fieldCheckbox.Value = repo.component_value

			h = fieldCheckbox.appendCheckbox()
			fields = append(fields, h)
		}

	}

	/*
		for i := 1; i < 10; i++ {
			if i == 1 {
				field.template = tplDateTemplate

				field.Id = "customer" + strconv.Itoa(i)
				field.Label = field.Id
				field.Name = field.Id

				h = field.appendDate()
				fields = append(fields, h)
			} else {
				field.template = tplInputTemplate

				field.Id = "customer" + strconv.Itoa(i)
				field.Label = field.Id
				field.Name = field.Id

				field.Placeholder = field.Id
				field.Pattern = ""
				field.ReqSize = "100"
				field.MaxLength = "8"
				field.MinLength = "1"
				field.ErrMsg = "Error number 5"

				h = field.appendText()
				fields = append(fields, h)
			}

			//fmt.Println(h)
		}
	*/

	//fmt.Println(fields)
	form.Fields = fields
	data := map[string]interface{}{
		"Title":   template.HTML(form.fieldsToString()),
		"Message": "",
		"reqnum":  reqnumber,
	}
	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, data)
}

func queryUserRepos(repos *repositoryUsers) error {
	rows, err := db.Query(`select a.tel_chat_id,a.cnt from (select tel_chat_id,coalesce((select count(*) from messages where tel_chat_id=m.tel_chat_id and viewedat is null group by m.tel_chat_id),0) as cnt from messages m where tel_chat_id is not null group by tel_chat_id) a order by a.cnt desc`)

	if err != nil {
		return err
	}
	defer rows.Close()
	var i = 0
	repo := repositoryUser{}
	for rows.Next() {
		i = 1

		err = rows.Scan(&repo.User, &repo.UnviewedCnt)
		if err != nil {
			return err
		}
		repos.Repos = append(repos.Repos, repo)
	}
	if i == 0 {
		repo.User = "0"
		repo.UnviewedCnt = 0
		repos.Repos = append(repos.Repos, repo)
	}

	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryMessageRepos(repos *repositoryMesages) error {
	//rows, err := db.Query(`select m.id,m.text,m.sent,m.sentby,m.tel_chat_id,m.tel_message_id,m.message_type,coalesce(cast(m.viewedBy as varchar),''),coalesce(cast(m.viewedAt as varchar),''),coalesce(cast(m.replyto as varchar),''),encode(v.voice::bytea,'hex') as hex_voice from messages m left join voices v on m.id=v.messages_id`)
	rows, err := db.Query(`select m.id,coalesce(m.text,'') as text,m.sent,m.sentby,coalesce(cast(m.tel_chat_id as varchar),'') as  tel_chat_id,coalesce(cast(m.tel_message_id as varchar),'') as tel_message_id,coalesce(cast(m.message_type as varchar),'') as message_type,coalesce(cast(m.viewedBy as varchar),'') as  viewedBy,coalesce(cast(m.viewedAt as varchar),'') as viewedAt,coalesce(cast(m.replyto as varchar),'') as replyTo,coalesce(cast(v.duration as varchar),'') as duration,coalesce(encode(v.voice::bytea,'hex'),'') as hex_voice from messages m left join voices v on m.id=v.messages_id where tel_chat_id is not null`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositoryMessage{}
		err = rows.Scan(&repo.Id,
			&repo.Text,
			&repo.Sent,
			&repo.SentBy,
			&repo.TelChatId,
			&repo.TelMessageId,
			&repo.MessageType,
			&repo.ViewedBy,
			&repo.ViewedAt,
			&repo.RepltyTo,
			&repo.Duration,
			&repo.VoiceText)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryMessageReposById(repos *repositoryMesages, id string) error {

	rows, err := db.Query(`select m.id,coalesce(m.text,'') as  text,m.sent,m.sentby,coalesce(cast(m.tel_chat_id as varchar),'') as tel_chat_id,coalesce(cast(m.tel_message_id as varchar),'') as tel_message_id,coalesce(cast(m.message_type as varchar),'') as  message_type,coalesce(cast(m.viewedBy as varchar),'') as  viewedBy,coalesce(cast(m.viewedAt as varchar),'') as viewedAt,coalesce(cast(m.replyto as varchar),'') as replyTo,coalesce(cast(v.duration as varchar),'') as duration,coalesce(encode(v.voice::bytea,'hex'),'') as hex_voice from messages m left join voices v on m.id=v.messages_id where tel_chat_id is not null and tel_chat_id=$1 order by sent asc`, id)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositoryMessage{}
		err = rows.Scan(&repo.Id,
			&repo.Text,
			&repo.Sent,
			&repo.SentBy,
			&repo.TelChatId,
			&repo.TelMessageId,
			&repo.MessageType,
			&repo.ViewedBy,
			&repo.ViewedAt,
			&repo.RepltyTo,
			&repo.Duration,
			&repo.VoiceText)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryMessagesCountReposById(repos *repositoryMessagesCountArr, id string) error {
	rows, err := db.Query(`select count(*) as cnt from messages m left join voices v on m.id=v.messages_id where tel_chat_id is not null and tel_chat_id=$1`, id)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositoryMessagesCount{}
		err = rows.Scan(&repo.Count)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryRequestTypes(repos *repositoryRequestTypeArr) error {

	rows, err := db.Query(`select id,name from request_type`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositoryRequestType{}
		err = rows.Scan(&repo.Id, &repo.Name)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryServiceRequests(repos *repositoryServiceRequestArr, reqtypeid int64) error {

	rows, err := db.Query(`select sr.id,sr.service_name from  servicesrequests sr where request_type_id=$1`, reqtypeid)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositoryServiceRequest{}
		err = rows.Scan(&repo.Id, &repo.Name)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryServiceRequestToClient(repos *repositoryServiceRequestToClientArr) error {

	rows, err := db.Query(`select rt.name as req_type_name,sr.service_name from request_type rt left join servicesrequests sr on rt.id=sr.request_type_id`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := repositoryServiceRequestToClient{}
		err = rows.Scan(&repo.TypeName, &repo.ServiceName)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

func queryServiceRequestReq(repos *repositoryServiceRequestReqArr, reqFrom int64, servicesrequestsid int64) error {

	var id int64
	err := db.QueryRow("insert into requests(reqfrom,servicesrequestsid,status) values($1,$2,0) returning id;", reqFrom, servicesrequestsid).Scan(&id)

	if err != nil {
		return err
	}
	err = db.QueryRow(`update requests set reqnumber=luhn_generate($1)::numeric where id=$2 returning reqnumber`, id, id).Scan(&id)
	if err != nil {
		return err
	}
	repo := repositoryServiceRequestReq{}
	repo.Id = id
	repos.Repos = append(repos.Repos, repo)

	return nil
}

func queryRepoRequests(repos *RepoRequestArr) error {

	rows, err := db.Query(`
select  rt.name as req_type_name,
        s.id,
        s.service_name as req_subtype_name,
        coalesce(cast(r.reqnumber as  varchar),''),
        r.datetime as reqdate,
        r.status

        from request_type rt
   join servicesrequests s on rt.id = s.request_type_id
   join requests r on r.servicesrequestsid=s.id
   where r.status>0`)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		repo := RepoRequest{}
		err = rows.Scan(&repo.ReqSubTypeName,
			&repo.ReqSubTypeId,
			&repo.ReqSubTypeName,
			&repo.ReqNumber,
			&repo.ReqDate,
			&repo.Status)
		if err != nil {
			return err
		}

		repos.Repos = append(repos.Repos, repo)
	}
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

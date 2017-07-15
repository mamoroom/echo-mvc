package script

import (
	"github.com/patrickmn/go-cache"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"

	"github.com/mamoroom/echo-mvc/app/lib/logger"
	"github.com/mamoroom/echo-mvc/app/lib/mail_data"
	"github.com/mamoroom/echo-mvc/app/lib/nazo_data"
	"github.com/mamoroom/echo-mvc/app/lib/util"
	"github.com/mamoroom/echo-mvc/app/models"
	"github.com/mamoroom/echo-mvc/app/models/entity"
	"github.com/mamoroom/echo-mvc/app/service/user"

	"bytes"
	_ "crypto/tls"
	"encoding/base64"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"
)

type ResChecNazoMaster struct {
	IsAvailable bool
	NazoId      int64
	OpenAt      time.Time
}

var PRE_EXEC_INTERVAL = 3600
var SEND_MAIL_INTERVAL_PER_WORKER = 3
var MAX_RCPT_TO_NUM = 100
var data_cache *cache.Cache

func init() {
	data_cache = cache.New(time.Duration(PRE_EXEC_INTERVAL)*time.Second, time.Duration(PRE_EXEC_INTERVAL+100)*time.Second)
}

func CheckNazoMaster(now time.Time) (*ResObject, *ErrorParam) {
	// nazo master取得する
	nazo_master := models.NewNazoMasterR()
	nazo_master_entites, err := nazo_master.FindsAndGetEntities()
	if err != nil {
		return nil, &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "DatabaseAccessError",
			Msg:       "Could not get nazo master",
			Param: map[string]interface{}{
				"detail": err.Error(),
			},
		}
	}

	for _, nazo_master_entity := range *nazo_master_entites {
		if IsIgnoreTarget(nazo_master_entity.Id) {
			continue
		}
		if IsNextExecTarget(now, nazo_master_entity) {
			return resSucceded(&ResChecNazoMaster{
				IsAvailable: true,
				NazoId:      nazo_master_entity.Id,
				OpenAt:      nazo_master_entity.OpenAt,
			}), nil
		}
	}

	// no target nazo master
	return resSucceded(&ResChecNazoMaster{
		IsAvailable: false,
	}), nil

}

func IsIgnoreTarget(nazo_id int64) bool {
	for _, id := range conf.Email.IgnoreNazoId {
		if nazo_id == id {
			return true
		}
	}
	return false
}

func IsNextExecTarget(now time.Time, nazo_master_entity entity.NazoMaster) bool {
	duration := nazo_master_entity.OpenAt.Sub(now)
	//fmt.Println(nazo_master_entity.Id, ":", duration)
	return duration.Seconds() > 0 && duration.Seconds() < float64(PRE_EXEC_INTERVAL)
}

func GetAndSaveHtmlMailTmpl(nazo_id int64, lang string, is_app_released bool, path string) error {

	// make request
	url := conf.Server.Domain + "/mailmagazine/"
	//url := "http://localhost/for_localhost_demo/html_mail.php"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	// make query param
	q := req.URL.Query()
	q.Add("n", util.CastInt64ToStr(nazo_id))
	q.Add("l", lang)
	if is_app_released {
		q.Add("r", "t")
	}

	//q.Add("e", "dev") (仮
	env := os.Getenv("CONFIGOR_ENV")
	switch env {
	case "local":
		q.Add("e", "dev")
	case "dev":
		q.Add("e", env)
	case "stg":
		q.Add("e", env)
	}

	req.URL.RawQuery = q.Encode()

	// request
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// read body
	resp_body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// save io
	err = ioutil.WriteFile(path, resp_body, 0755)
	if err != nil {
		return err
	}

	return nil
}

func GetEmails(nazo_id int64) (map[string][][]string, *ErrorParam) {
	var cache_key_prefix string = util.CastInt64ToStr(nazo_id)
	var cache_key_suffix string = "email"

	_target, found := data_cache.Get(cache_key(cache_key_prefix, cache_key_suffix))
	if !found {
		fmt.Println("get from db")
		all_data, err_param := GetEmailsFromDB()
		if err_param != nil {
			return nil, err_param
		}
		data_cache.Set(cache_key(cache_key_prefix, cache_key_suffix), &all_data, cache.DefaultExpiration)
		_target = &all_data
	}

	if _target == nil {
		return nil, &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "CacheError",
			Msg:       "Could not get emails from cache",
		}
	}

	t, ok := _target.(*map[string][][]string)
	if !ok {
		return nil, &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "CacheDataFormatError",
			Msg:       "Could not cast string slice",
		}
	}
	return *t, nil
}

func cache_key(key_prefix string, key_suffix string) string {
	return key_prefix + "__" + key_suffix
}

func GetEmailsFromDB() (map[string][][]string, *ErrorParam) {

	// user emails
	user := models.NewUserR()
	user_entities, err := user.GetsNotificationUsers()
	if err != nil {
		return nil, &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "DatabaseAccessError",
			Msg:       "Could not get nazo master",
			Param: map[string]interface{}{
				"detail": err.Error(),
			},
		}
	}

	// bounce mails
	mail_bounce_master := models.NewMailBounceMasterR()
	cnt_mail_bounce_master_of_map, err := mail_bounce_master.FindsAndCntlMap()

	var all_emails = make(map[string][][]string)
	var send_worker_emails_of = make(map[string][]string)
	var send_worker_emails_lang_index = map[string]int{
		"ja":   0,
		"en":   0,
		"es":   0,
		"fr":   0,
		"it":   0,
		"de":   0,
		"ko":   0,
		"nl":   0,
		"text": 0, //キャリアメールはjaとみなす
	}

	for _, user_entity := range *user_entities {

		// bounce mailは無視する
		if (*cnt_mail_bounce_master_of_map)[user_entity.AuthEmail] >= conf.Email.InvalidMailbounceCnt {
			continue
		}

		// langの決定. lang or text(=キャリアメール)
		var lang string
		if util.CheckCarrierDomainValidation(user_entity.AuthEmail) {
			lang = "text"
		} else {
			lang = user_entity.Lang
		}

		if send_worker_emails_lang_index[lang] > MAX_RCPT_TO_NUM-1 {
			all_emails[lang] = append(all_emails[lang], send_worker_emails_of[lang])
			send_worker_emails_lang_index[lang] = 0
			send_worker_emails_of[lang] = make([]string, MAX_RCPT_TO_NUM)
		}

		_, ok := send_worker_emails_of[lang]
		if !ok {
			send_worker_emails_of[lang] = make([]string, MAX_RCPT_TO_NUM)
		}

		send_worker_emails_of[lang][send_worker_emails_lang_index[lang]] = user_entity.AuthEmail
		send_worker_emails_lang_index[lang]++
	}

	for lang, emails := range send_worker_emails_of {
		all_emails[lang] = append(all_emails[lang], emails[:send_worker_emails_lang_index[lang]])
	}
	return all_emails, nil

}

type SendParam struct {
	Lang   string
	Emails []string
}

func SendMail(nazo_id int64, target_tmpl_path string, all_emails map[string][][]string, now time.Time, max_client_num int) *ErrorParam {
	// to get nazo label
	nazo_master := models.NewNazoMasterR()
	nazo_master_entity, err := nazo_master.GetById(nazo_id)
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "DatabaseAccessError",
			Msg:       "Could not get nazo master",
			Param: map[string]interface{}{
				"detail": err.Error(),
			},
		}
	}

	// model
	user_model := models.NewUserR()
	dummy_user_entity := user_model.GetEmptyEntity()
	// [todo]: nowはscript実行時間なので40minほど評価時間とは差があることに注意!
	service_user := user.New(&dummy_user_entity, now)
	is_app_released := service_user.IsAppReleased()
	if conf.Email.IsDebug {
		is_app_released = conf.Email.IsDebugAppReleased
	}

	var wg sync.WaitGroup

	var client_num int
	for _, lang_emails := range all_emails {
		client_num += len(lang_emails)
	}

	if client_num > max_client_num {
		client_num = max_client_num
	}

	q := make(chan *SendParam, client_num)
	e := make(chan *ErrorParam, client_num)

	domains := conf.Email.Domains
	for i := 0; i < client_num; i++ {
		wg.Add(1)
		domain := domains[i%len(domains)]
		go SendMailClientWorker(&nazo_master_entity, is_app_released, target_tmpl_path, domain, &wg, q, e)
	}

	// [todo]: エラー処理だけ適切にできるかわからないので要検証...
	var err_param *ErrorParam
	for lang, lang_emails := range all_emails {
		for _, emails := range lang_emails {
			send_param := &SendParam{
				Lang:   lang,
				Emails: emails,
			}
			select {
			case err_param = <-e:
				break
			case q <- send_param:
			}
		}
	}
	close(q)
	wg.Wait()

	return err_param
}

func SendMailClientWorker(_nazo_master_entity *entity.NazoMaster, is_app_released bool, target_tmpl_path string, domain string, wg *sync.WaitGroup, q chan *SendParam, e chan<- *ErrorParam) {
	defer wg.Done()

	// make connection
	for {
		// block //
		_send_param, ok := <-q
		if !ok {
			return
		}
		err_param := SendMailByClient(_nazo_master_entity, is_app_released, target_tmpl_path, _send_param, domain)
		time.Sleep(time.Duration(SEND_MAIL_INTERVAL_PER_WORKER) * time.Second)
		if err_param != nil {
			e <- err_param
			return
		}
	}

}

func SendMailByClient(_nazo_master_entity *entity.NazoMaster, is_app_released bool, target_tmpl_path string, _send_param *SendParam, domain string) *ErrorParam {
	// title
	nazo_data_json, err := nazo_data.GetNazoData(_nazo_master_entity, "nazo_master")
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "NazoJsonDataIoError",
			Msg:       "Could not get nazo json data",
			Param: map[string]interface{}{
				"detail": err.Error(),
				"users":  strings.Join(_send_param.Emails, ","),
			},
		}
	}

	switch _send_param.Lang {
	// textメールはjaのみ想定
	case "text":

		// subject
		mail_json_data, err := mail_data.GetMailData("mail_val")
		if err != nil {
			return &ErrorParam{
				Logger:    logger.SendMailBatchLogger,
				LogFunc:   "Error",
				ErrorType: "MailJsonDataIoError",
				Msg:       "Could not get mail json data",
				Param: map[string]interface{}{
					"detail": err.Error(),
					"users":  strings.Join(_send_param.Emails, ","),
				},
			}
		}
		mail_json_subject := mail_json_data.Data["subject"].(map[string]interface{})["ja"].(string)
		nazo_title := _nazo_master_entity.GetLabel() + ". " + nazo_data_json.NazoListData["title_trans8"].(map[string]interface{})["ja"].(string)
		subject := strings.Replace(mail_json_subject, "$nazo_title", nazo_title, -1)
		encoded_subject, err := EncodingToISO2022JP(subject)
		if err != nil {
			return &ErrorParam{
				Logger:    logger.SendMailBatchLogger,
				LogFunc:   "Error",
				ErrorType: "SubjectEncodeError",
				Msg:       "Could not encode subject data",
				Param: map[string]interface{}{
					"detail": err.Error(),
					"users":  strings.Join(_send_param.Emails, ","),
				},
			}
		}

		//body
		mail_static_body := conf.Email.TextMail.BodyBeforeRelease
		if is_app_released {
			mail_static_body = conf.Email.TextMail.BodyAfterRelease
		}
		__body := strings.Replace(mail_static_body, "$nazo_title", nazo_title, -1)
		uri_suffix := util.CastInt64ToStr(_nazo_master_entity.Id) + "?utm_source=mail&utm_medium=email&utm_campaign=lmj"
		_body := strings.Replace(__body, "$nazo_id", uri_suffix, -1)
		body := strings.Replace(_body, "$domain", conf.Server.Domain, -1)
		return SendTextMailByClient(encoded_subject, body, _send_param.Emails, domain)

	default:
		mail_json_data, err := mail_data.GetMailData("mail_val")
		if err != nil {
			return &ErrorParam{
				Logger:    logger.SendMailBatchLogger,
				LogFunc:   "Error",
				ErrorType: "MailJsonDataIoError",
				Msg:       "Could not get mail json data",
				Param: map[string]interface{}{
					"detail": err.Error(),
					"users":  strings.Join(_send_param.Emails, ","),
				},
			}
		}
		mail_json_from := mail_json_data.Data["from"].(map[string]interface{})[_send_param.Lang].(string)
		mail_json_subject := mail_json_data.Data["subject"].(map[string]interface{})[_send_param.Lang].(string)
		subject := strings.Replace(mail_json_subject, "$nazo_title", _nazo_master_entity.GetLabel()+". "+nazo_data_json.NazoListData["title_trans8"].(map[string]interface{})[_send_param.Lang].(string), -1)
		html_file_path := target_tmpl_path + "/" + _send_param.Lang + ".html"
		return SendHtmlMailByClient(mail_json_from, subject, html_file_path, _send_param.Emails, domain)

	}
}

func SendTextMailByClient(encoded_subject string, body string, send_users_email []string, domain string) *ErrorParam {
	// Send the email body.
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(domain + ":" + conf.Email.Port)
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "SmtpDialError",
			Msg:       "Could not create smtp client",
			Param: map[string]interface{}{
				"detail": err.Error(),
			},
		}
	}
	defer c.Quit()

	// Set the sender and recipient.
	c.Mail(conf.Email.From) // メールの送り主を指定
	for _, email := range send_users_email {
		c.Rcpt(email) // 受信者を指定
	}

	// request
	wc, err := c.Data()
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "SmtpConnectionError",
			Msg:       "Could not connect smtp",
			Param: map[string]interface{}{
				"detail": err.Error(),
				"users":  strings.Join(send_users_email, ","),
			},
		}
	}
	defer wc.Close()

	buf := bytes.NewBufferString("To:recieve@layton.world")
	buf.WriteString("\n")
	buf.WriteString("Subject:" + conf.Email.TextMail.SubjectPrefix + encoded_subject + conf.Email.TextMail.SubjectSuffix) //件名
	buf.WriteString("\n")
	buf.WriteString("MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";")
	buf.WriteString("\r\n")
	buf.WriteString(body)
	if _, err = buf.WriteTo(wc); err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "HtmlBufferWriteError",
			Msg:       "Could not convert html string to buffer",
			Param: map[string]interface{}{
				"detail": err.Error(),
				"users":  strings.Join(send_users_email, ","),
			},
		}
	}
	return nil
}

func EncodingToISO2022JP(str string) (string, error) {
	reader := strings.NewReader(str)
	transformer := japanese.ISO2022JP.NewEncoder()

	iso_encoded_byte, err := ioutil.ReadAll(transform.NewReader(reader, transformer))
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(iso_encoded_byte), nil
}

func SendHtmlMailByClient(from string, subject string, html_file_path string, send_users_email []string, domain string) *ErrorParam {
	// Send the email body.
	// Connect to the remote SMTP server.
	c, err := smtp.Dial(domain + ":" + conf.Email.Port)
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "SmtpDialError",
			Msg:       "Could not create smtp client",
			Param: map[string]interface{}{
				"detail": err.Error(),
			},
		}
	}
	defer c.Quit()

	// Set the sender and recipient.
	c.Mail(conf.Email.From) // メールの送り主を指定
	for _, email := range send_users_email {
		c.Rcpt(email) // 受信者を指定
	}

	// request
	wc, err := c.Data()
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "SmtpConnectionError",
			Msg:       "Could not connect smtp",
			Param: map[string]interface{}{
				"detail": err.Error(),
				"users":  strings.Join(send_users_email, ","),
			},
		}
	}
	defer wc.Close()

	buf := bytes.NewBufferString("To:recieve@layton.world")
	buf.WriteString("\n")
	buf.WriteString("Subject:" + subject) //件名
	buf.WriteString("\n")
	buf.WriteString("MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";")
	buf.WriteString("\n")
	buf.WriteString("From:" + from + "<" + conf.Email.From + ">") //署名
	buf.WriteString("\r\n")

	html_string, err := ParseTemplate(html_file_path, struct{}{})
	if err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "HtmlParseError",
			Msg:       "Could not parse html",
			Param: map[string]interface{}{
				"detail": err.Error(),
				"users":  strings.Join(send_users_email, ","),
			},
		}
	}
	buf.WriteString(html_string)
	if _, err = buf.WriteTo(wc); err != nil {
		return &ErrorParam{
			Logger:    logger.SendMailBatchLogger,
			LogFunc:   "Error",
			ErrorType: "HtmlBufferWriteError",
			Msg:       "Could not convert html string to buffer",
			Param: map[string]interface{}{
				"detail": err.Error(),
				"users":  strings.Join(send_users_email, ","),
			},
		}
	}
	return nil
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

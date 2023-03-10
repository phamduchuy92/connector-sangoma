package main

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/CyCoreSystems/agi"
	"github.com/valyala/fasthttp"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

type NotificationRequest struct {
	Mobile   string `json:"mobile"`
	Status   string `json:"status"`
	Ext      string `json:"ext"`
	Datetime string `json:"datetime"`
}

var authorization = "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL2dyYXRpc29mdC50ZWNoLyIsInVwbiI6InNhbmdvbWEiLCJncm91cHMiOlsiMTE5NWJmZTAtZWRmMi00YzBlLThiNDktZTU5ZDI2MDcyYmM1Il0sImV4cCI6MzE1NTY4ODk4NjQ0MDMxOTksImlhdCI6MTY1Mzg5NTEyNywianRpIjoiMGY0NDA2OWUtMmM2NS00NmM2LWJjNGEtNjg1NDVmMmQ1MGE0In0.Ol4elsuUI6gFh1PIXbalRD3G_6G6D5783cPW-7FTVRLFHu5l-14NH0qBBU4r39xOLcUfD-UxA_tUUoUUS47IMtC17tVCP6ITFjhUvH31e3RkNywg_iHSf70zHfgMm_v1l9VPmhuBeT1ReJ57im2-GZAQA0L2KaWqQtPK4ZlrzwbMmn-zkHDdXa9Y0YalBdg_VTGqBsfvEc_bLef8Jq7rZ8TlYz6K7NZKqBJWc8K1vDyreCWyx9F6wT6RN0pTVyRQFIrIY7pqKfCdlx2WsSIDsidfBfnBwnj_5TqJofATNRWPkhEee8dKqgyBHSHBVLG1JEYM1tpFukcx7nbqwDXpFA"
var username = "sangoma"
var password = "hXKv_ngMEZK$vE8*@55aNYP#"
var webClient = &fasthttp.Client{
	TLSConfig: &tls.Config{InsecureSkipVerify: true},
}

func main() {
	viper.SetConfigName("application") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv() // read value ENV variable

	setDefaultConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Unable to read configuration file `application.yaml`. Using defaults")
	}

	log.Printf("Application started on [%s]", viper.GetString("tcpAddress"))
	agi.Listen(viper.GetString("tcpAddress"), handler)
}

// setDefaultConfig set viper configuration when nothing inside
func setDefaultConfig() {
	viper.SetDefault("tcpAddress", "localhost:4573")
	viper.SetDefault("notify.apiEndpoint", "https://api.sisvietnam.vn/gis/restful/callcenter/webhook/updateCallInfo")
	viper.SetDefault("login.apiEndpoint", "https://api.sisvietnam.vn/gis/restful/callcenter/login")
	viper.SetDefault("login.username", username)
	viper.SetDefault("login.password", password)
}

func handler(a *agi.AGI) {
	defer a.Close()
	var status = "RINGING"
	login(viper.GetString("login.apiEndpoint"), viper.GetString("login.username"), viper.GetString("login.password"))
	log.Printf("Login with %s[%s]", viper.GetString("login.username"), viper.GetString("login.password"))

	callerid, err := a.Get("CALLERID(num)")
	if err != nil {
		log.Printf("Cannot detect calling number")
	}
	log.Printf("callerid %s", callerid)
	exten, err := a.Get("EXTEN")
	if err != nil {
		log.Printf("Cannot detect exten number")
	}
	log.Printf("exten %s", exten)
	channel, err := a.Get("CHANNEL")
	if err != nil {
		log.Printf("Cannot detect channel")
	}
	log.Printf("channel %s", channel)
	dialstatus, err := a.Get("DIALSTATUS")
	if err != nil {
		log.Printf("Cannot detect status")
	}
	log.Printf("dialstatus %s", dialstatus)
	r, _ := regexp.Compile("^Local\\/(.*?)\\@{1}.*")
	if len(r.FindStringSubmatch(channel)) > 0 && exten == "" {
		exten = r.FindStringSubmatch(channel)[1]
		status = "IN_PROGRESS"
	}

	notify(viper.GetString("notify.apiEndpoint"), callerid, exten, status)

	a.Close()
}

func notify(url string, calling string, called string, status string) error {
	notification := NotificationRequest{
		Mobile:   calling,
		Status:   status,
		Ext:      called,
		Datetime: time.Now().Format("2006-01-02 15:04:05"),
	}
	log.Printf("%s | %s -> %s", notification.Datetime, calling, called)
	json_data, err := json.Marshal(notification)
	if err != nil {
		log.Fatal(err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("%s", authorization))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(json_data)
	log.Printf("Body %v", string(json_data))
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Perform the request
	err = webClient.Do(req, resp)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode() == 200 {
		log.Println("Successfully send to CRM")
	} else if resp.StatusCode() == 401 {
		log.Println("Unauthorized")
	} else {
		log.Printf("Internal error with status %d", resp.StatusCode())
	}

	return nil
}

// login perform a login request to server
func login(url string, username string, password string) error {
	var err error
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(url)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", basicAuth(username, password)))
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Perform the request
	err = webClient.Do(req, resp)
	if err != nil {
		log.Println(err)
	}
	if resp.StatusCode() < 300 {
		authorization = fmt.Sprintf("Bearer %s", string(resp.Body()))
		return nil
	}

	log.Printf("Failed to login request with response: %d", resp.StatusCode())
	return nil
}

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

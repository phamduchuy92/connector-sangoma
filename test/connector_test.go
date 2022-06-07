package connector

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

type NotificationRequest struct {
	SDT      string
	status   string
	ext      string
	datetime string
}

func TestLogin(t *testing.T) {
	var webClient = &fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	var username = "sangoma"
	var password = "hXKv_ngMEZK$vE8*@55aNYP#"
	var err error
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://payment.sisvietnam.vn/gis/restful/callcenter/login")
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

	log.Printf("Resp: %+v", string(resp.Body()))
}

func TestNotify(t *testing.T) {
	var err error
	var webClient = &fasthttp.Client{
		TLSConfig: &tls.Config{InsecureSkipVerify: true},
	}
	notification := NotificationRequest{
		SDT:      "123456",
		status:   "ANSWER",
		ext:      "500",
		datetime: time.Now().Format("200601021504"),
	}
	json_data, err := json.Marshal(notification)
	if err != nil {
		log.Fatal(err)
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI("http://payment.sisvietnam.vn/gis/restful/callcenter/webhook/updateCallInfo")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJpc3MiOiJodHRwczovL2dyYXRpc29mdC50ZWNoLyIsInVwbiI6InNhbmdvbWEiLCJncm91cHMiOlsiMTE5NWJmZTAtZWRmMi00YzBlLThiNDktZTU5ZDI2MDcyYmM1Il0sImV4cCI6MzE1NTY4ODk4NjQ0MDMxOTksImlhdCI6MTY1Mzg5NTEyNywianRpIjoiMGY0NDA2OWUtMmM2NS00NmM2LWJjNGEtNjg1NDVmMmQ1MGE0In0.Ol4elsuUI6gFh1PIXbalRD3G_6G6D5783cPW-7FTVRLFHu5l-14NH0qBBU4r39xOLcUfD-UxA_tUUoUUS47IMtC17tVCP6ITFjhUvH31e3RkNywg_iHSf70zHfgMm_v1l9VPmhuBeT1ReJ57im2-GZAQA0L2KaWqQtPK4ZlrzwbMmn-zkHDdXa9Y0YalBdg_VTGqBsfvEc_bLef8Jq7rZ8TlYz6K7NZKqBJWc8K1vDyreCWyx9F6wT6RN0pTVyRQFIrIY7pqKfCdlx2WsSIDsidfBfnBwnj_5TqJofATNRWPkhEee8dKqgyBHSHBVLG1JEYM1tpFukcx7nbqwDXpFA"))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetBody(json_data)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	// Perform the request
	err = webClient.Do(req, resp)
	if err != nil {
		log.Println(err)
	}

	log.Printf("Resp: %+v", string(resp.Body()))
}

func basicAuth(username string, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func TestRegex(t *testing.T) {
	match, _ := regexp.MatchString("^Local\\/(.*?)\\@{1}.*?", "Local/500@from-queue-00000273;1")
	fmt.Println(match)
	r, _ := regexp.Compile("^Local\\/(.*?)\\@{1}.*")
	fmt.Println(r.FindStringSubmatch("Local/500@from-queue-00000273;1")[1])
}

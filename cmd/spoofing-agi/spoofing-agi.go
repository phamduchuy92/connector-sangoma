package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"net/http"
	"net/url"

	"github.com/CyCoreSystems/agi"

	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"

	uuid "github.com/google/uuid"
)

var (
	sipCauseCode = make(map[int]int)
	// config from consul
	runtimeViper = viper.New()
)

func main() {
	viper.SetConfigName("application") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv() // read value ENV variable

	setDefaultConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Println("Unable to read configuration file `application.yaml`. Using defaults")
	}

	reloadCausecodeMapping()

	log.Printf("Application started on [%s]", viper.GetString("tcpAddress"))
	agi.Listen(viper.GetString("tcpAddress"), handler)
}

// setDefaultConfig set viper configuration when nothing inside
func setDefaultConfig() {
	viper.SetDefault("tcpAddress", ":4567")
	viper.SetDefault("timeout", 1000)
	viper.SetDefault("channel", "SIP")
	viper.SetDefault("dryRun", false)
	viper.SetDefault("defaultCauseCode", 1)
}

// reloadCausecodeMapping reload application causecode from viper
func reloadCausecodeMapping() {
	causeCodeMap := viper.GetStringMapString("causeCodeMapping")
	for key, element := range causeCodeMap {
		httpCode, err := strconv.Atoi(key)
		if err != nil {
			log.Printf("Unable to parse HTTP code [%s]", key)
			continue
		}
		sipCode, err := strconv.Atoi(element)
		if err != nil {
			log.Printf("Unable to parse SIP code [%s]", element)
			continue
		}
		sipCauseCode[httpCode] = sipCode
	}
	log.Printf("Causecode Mapping: %+v", sipCauseCode)
}

func handler(a *agi.AGI) {
	defer a.Close()

	uniqID, err := a.Get("UNIQUEID")
	if err != nil {
		log.Printf("Cannot get call UNIQUEID")
	}

	calling, err := a.Get("CALLERID(num)")
	if err != nil {
		log.Printf("Cannot detect calling number")
	}
	called, err := a.Get("EXTEN")
	if err != nil {
		log.Printf("Cannot detect called number")
	}

	peerIP, err := a.Get("PEER_IP")
	if err != nil {
		peerIP = ""
	}
	s := strings.Split(peerIP, ":")
	peerIP = s[0]

	UUID := uuid.New()

	u, err := url.Parse(viper.GetString("apiEndpoint"))
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	q.Set("channel", viper.GetString("channel"))
	q.Set("calling", calling)
	q.Set("called", called)
	q.Set("peer", peerIP)
	q.Set("callref", UUID.String())
	u.RawQuery = q.Encode()

	causeCode := viper.GetInt("defaultCauseCode")
	if !viper.GetBool("dryRun") {
		log.Printf("[%s] | [%s] -> [%s] | [%s] >> [%s]", uniqID, calling, called, peerIP, u)
		causeCode = callback(u, uniqID)
	} else {
		log.Printf("[%s] DRYRUN | [%s] -> [%s] | [%s] >> [%s]", uniqID, calling, called, peerIP, u)
		go callback(u, uniqID)
	}
	log.Printf("[%s] | [%s] -> [%s] | causeCode [%d]", uniqID, calling, called, causeCode)
	a.Set("SPOOFING_CODE", fmt.Sprintf("%d", causeCode))
	a.Close()
	// a.Command(fmt.Sprintf("GOSUB %s code-%d 1", agiContext, causeCode)).Err()
	// a.Command(fmt.Sprintf("HANGUP %d", causeCode)).Err()
}

func callback(u *url.URL, uniqID string) int {
	httpClient := http.Client{
		Timeout: time.Duration(viper.GetInt("timeout")) * time.Millisecond,
	}
	causeCode := viper.GetInt("causeCodeDefault")
	resp, err := httpClient.Get(u.String())

	if err != nil {
		log.Printf("[%s] | ERR: [%s] | %v", uniqID, u, err)
		return causeCode
	}
	// switch {
	// case resp.StatusCode < 400:
	// 	log.Printf("[%s] | Success response", uniqID)
	// 	causeCode = 1
	// 	break
	// case resp.StatusCode < 500:
	// 	log.Printf("[%s] | Client errors", uniqID)
	// 	causeCode = 42
	// 	break
	// default:
	// 	log.Printf("[%s] | Server errors", uniqID)
	// 	causeCode = 1
	// 	break
	// }
	if sipCode, ok := sipCauseCode[resp.StatusCode]; ok {
		causeCode = sipCode
	} else {
		causeCode = resp.StatusCode
	}
	return causeCode
}

// ReadConfigConsul get config from consul
func ReadConfigConsul() {
	runtimeViper.AddRemoteProvider("consul", "localhost:8500", "config/test")
	runtimeViper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	err := runtimeViper.ReadRemoteConfig()
	if err != nil {
		log.Println("cannot read config from consul")
	}
	listenerCh := make(chan bool)

	// open a goroutine to watch remote changes forever
	go func() {
		for {
			err := runtimeViper.WatchRemoteConfigOnChannel()
			if err != nil {
				log.Printf("unable to read remote config: %v", err)
				continue
			}
			for {
				time.Sleep(time.Second * 5) // delay after each request
				listenerCh <- true
			}
		}
	}()
	for {
		select {
		case <-listenerCh:
			runtimeViper.ReadRemoteConfig()
			fmt.Printf("test %+v", runtimeViper)
		}
	}
}

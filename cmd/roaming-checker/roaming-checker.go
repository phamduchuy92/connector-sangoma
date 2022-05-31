package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"net/http"

	"github.com/gofrs/uuid"
	"github.com/spf13/viper"
	"gorm.io/gorm"

	"context"

	redis "github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	_ "github.com/spf13/viper/remote"
	"schneider.vip/problem"
)

// RoamingRequest hold the structure of request
type RoamingRequest struct {
	ClientID    string `json:"clientId,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	RequestID   string `json:"requestId,omitempty"`
	RequestDate string `json:"requestDate,omitempty"`
	SecureCode  string `json:"secureCode,omitempty"`
}

// RoamingResponse hold the structure of response
type RoamingResponse struct {
	ClientID    string `json:"clientId,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	RequestID   string `json:"requestId,omitempty"`
	RequestDate string `json:"requestDate,omitempty"`
	ResponseID  string `json:"responseId,omitempty"`
	IsRoaming   uint   `json:"isRoaming,omitempty"`
}

type RoamingCheckerLog struct {
	gorm.Model
	Prefix      string `json:"prefix,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Provider    string `json:"provider,omitempty"`
	IsRoaming   uint   `json:"isRoaming,omitempty"`
}

type MnpRequestVNPT struct {
	Msisdn string `json:"msisdn,omitempty"`
}

type MnpResponseVNPT struct {
	ErrorCode          string `json:"error_code,omitempty"`
	Message            string `json:"message,omitempty"`
	MnpStatus          string `json:"mnp_status,omitempty"`
	OriginProviderCode string `json:"origin_provider_code,omitempty"`
	ProviderCode       string `json:"provider_code,omitempty"`
}

type requestPayloadStruct struct {
	Prefix      string `json:"prefix,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
}

// config from consul
var runtime_viper = viper.New()
var ctx = context.Background()
var rdb *redis.Client

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	viper.SetConfigName("roaming-checker") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv()        // read value ENV variable
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		log.Println("Unable to read configuration file `roaming-checker.yaml`. Using defaults")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
	})

	// start server
	app := fiber.New()
	// app.Use(cors.New())
	app.Get("/search", ReverseProxy)
	app.Listen(":" + viper.GetString("roaming-checker.port"))
}

func ReverseProxy(c *fiber.Ctx) error {
	prefix := c.Query("prefix")
	phoneNumber := c.Query("phoneNumber")
	if prefix == "" || phoneNumber == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	resp_cache, err := rdb.Get(ctx, prefix+"_"+phoneNumber).Result()
	if err == redis.Nil || err != nil {
		t := time.Now()
		gen := uuid.NewGen()

		reqID, err := gen.NewV4()
		if err != nil {
			log.Fatal("Cannot generate UUID", err)
		}
		clientID := viper.GetString(prefix + ".clientID")
		if len(clientID) == 0 {
			log.Printf("Invalid endpoint %s.clientID", clientID)
			return c.SendStatus(http.StatusNotImplemented)
		}
		message := &RoamingRequest{
			ClientID:    clientID,
			PhoneNumber: phoneNumber,
			RequestID:   reqID.String(),
			RequestDate: t.Format("20060102150405"),
		}
		hash := fmt.Sprintf("%s%s%s%s%s", message.ClientID, message.PhoneNumber, message.RequestID, message.RequestDate, viper.GetString(prefix+".privateKey"))

		message.SecureCode = fmt.Sprintf("%x", sha256.Sum256([]byte(hash)))

		bytesRepresentation, err := json.Marshal(message)
		if err != nil {
			log.Println("Failed to marshall", err)
			return c.SendStatus(http.StatusBadGateway)
		}
		fmt.Printf("Send request: %s\n", bytesRepresentation)

		client := &http.Client{}
		apiEndpoint := viper.GetString(prefix + ".apiEndpoint")
		if len(apiEndpoint) == 0 {
			log.Printf("Invalid endpoint %s.apiEndpoint", prefix)
			return c.SendStatus(http.StatusBadGateway)
		}
		request, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(bytesRepresentation))
		if err != nil {
			log.Println("Failed to check status for " + message.PhoneNumber)
			return c.SendStatus(http.StatusBadGateway)
		}
		request.Header.Add("Content-Type", "application/json; charset=utf-8")
		request.Header.Add("Accept", "application/json")
		request.Header.Add("apikey", viper.GetString(prefix+".apiKey"))

		response, err := client.Do(request)
		if err != nil {
			log.Println("Failed to perform POST request", err)
			return c.SendStatus(http.StatusBadGateway)
		}

		responseBody, err := ioutil.ReadAll(response.Body)
		response.Body.Close()
		if err != nil {
			log.Println("Failed to read body", err)
			return c.SendStatus(http.StatusBadGateway)
		}

		fmt.Printf("Got response body: %T %s", responseBody, responseBody)
		resp := &RoamingResponse{}
		err = json.Unmarshal([]byte(responseBody), &resp)

		if err != nil {
			log.Println("Failed to unmarshall ", responseBody, err)
			return c.SendStatus(http.StatusBadGateway)
		}
		fmt.Printf("Got response: %+v\n", resp)
		rdb.Set(ctx, prefix+"_"+phoneNumber, resp, time.Duration(viper.GetInt("redis.timeout"))*time.Millisecond)

		log := &RoamingCheckerLog{
			Prefix:      prefix,
			PhoneNumber: phoneNumber,
			Provider:    checkProvider(prefix),
			IsRoaming:   resp.IsRoaming,
		}
		if resp.IsRoaming == 0 {
			return c.Status(viper.GetInt("isRoaming.0")).JSON(log)
		} else if resp.IsRoaming == 1 {
			return c.Status(viper.GetInt("isRoaming.1")).JSON(log)
		} else if resp.IsRoaming == 2 {
			return c.Status(viper.GetInt("isRoaming.2")).JSON(log)
		}

	} else {
		resp := &RoamingResponse{}
		json.Unmarshal([]byte(resp_cache), &resp)
		log := &RoamingCheckerLog{
			Prefix:      prefix,
			PhoneNumber: phoneNumber,
			Provider:    checkProvider(prefix),
			IsRoaming:   resp.IsRoaming,
		}

		if resp.IsRoaming == 0 {
			return c.Status(viper.GetInt("isRoaming.0")).JSON(log)
		} else if resp.IsRoaming == 1 {
			return c.Status(viper.GetInt("isRoaming.1")).JSON(log)
		} else if resp.IsRoaming == 2 {
			return c.Status(viper.GetInt("isRoaming.2")).JSON(log)
		}
	}
	return c.SendStatus(http.StatusOK)
}

// ReadConfigConsul get config from consul
func ReadConfigConsul() {
	runtime_viper.AddRemoteProvider("consul", "localhost:8500", "config/test")
	runtime_viper.SetConfigType("yaml") // because there is no file extension in a stream of bytes, supported extensions are "json", "toml", "yaml", "yml", "properties", "props", "prop", "env", "dotenv"

	// read from remote config the first time.
	err := runtime_viper.ReadRemoteConfig()
	if err != nil {
		log.Println("cannot read config from consul")
		return
	}
	listenerCh := make(chan bool)

	// open a goroutine to watch remote changes forever
	go func() {
		for {
			err := runtime_viper.WatchRemoteConfigOnChannel()
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
			runtime_viper.ReadRemoteConfig()
			fmt.Printf("RemoteConfig was updated %+v", runtime_viper)
		}
	}
}

func checkProvider(prefix string) string {
	if prefix == "84005" {
		return "VIETNAMMOBILE"
	} else if prefix == "84004" || prefix == "84006" || prefix == "84008" {
		return "VIETTEL"
	} else if prefix == "84001" {
		return "MOBIFONE"
	}

	return "VINAPHONE"
}

// ProblemJSONErrorHandle send error handle
func ProblemJSONErrorHandle(ctx *fiber.Ctx, err error) error {
	// Statuscode defaults to 500
	code := fiber.StatusInternalServerError

	// Retreive the custom statuscode if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send custom error page
	err = ctx.Status(code).JSON(problem.Of(code))
	if err != nil {
		// In case the SendFile fails
		return ctx.Status(500).JSON(problem.Of(500))
	}

	// Return from handler
	return nil
}

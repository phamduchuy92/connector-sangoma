package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"

	"bytes"
	"io/ioutil"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"context"

	redis "github.com/go-redis/redis/v8"
	_ "github.com/spf13/viper/remote"
	"gorm.io/gorm"
)

type MnpCheckerRequest struct {
	Msisdn string `json:"msisdn,omitempty"`
}

type MnpCheckerResponse struct {
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

type MnpCheckerLog struct {
	gorm.Model
	Prefix      string `json:"prefix,omitempty"`
	PhoneNumber string `json:"phoneNumber,omitempty"`
	Provider    string `json:"provider,omitempty"`
	IsRoaming   uint   `json:"isRoaming,omitempty"`
}

// config from consul
var runtime_viper = viper.New()
var ctx = context.Background()
var rdb *redis.Client

func main() {
	viper.SetConfigName("mnp-checker") // config file name without extension
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")
	viper.AutomaticEnv()        // read value ENV variable
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		// setDefaultConfig()
		log.Println("Unable to read configuration file `mnp-checker.yaml`. Using defaults")
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis.address"),
		Password: viper.GetString("redis.password"), // no password set
		DB:       viper.GetInt("redis.db"),          // use default DB
	})

	// start server
	app := fiber.New()
	// app.Use(cors.New())
	app.Get("/search", CheckMnp)
	app.Listen(":" + viper.GetString("mnp-checker.port"))
}

func CheckMnp(c *fiber.Ctx) error {
	msisdn := c.Query("msisdn")
	if msisdn == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	message := &MnpCheckerRequest{
		Msisdn: msisdn,
	}

	bytesRepresentation, err := json.Marshal(message)
	if err != nil {
		log.Println("Failed to marshall", err)
		return c.SendStatus(http.StatusBadGateway)
	}
	fmt.Printf("Send request: %s\n", bytesRepresentation)

	client := &http.Client{}
	apiEndpoint := viper.GetString("endpoint.apiEndpoint")

	request, err := http.NewRequest("POST", apiEndpoint, bytes.NewBuffer(bytesRepresentation))
	if err != nil {
		log.Println("Failed to check status for " + message.Msisdn)
		return c.SendStatus(http.StatusBadGateway)
	}
	request.Header.Add("Content-Type", "application/json; charset=utf-8")
	request.Header.Add("Accept", "application/json")

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
	resp := &MnpCheckerResponse{}
	err = json.Unmarshal([]byte(responseBody), &resp)

	if err != nil {
		log.Println("Failed to unmarshall ", responseBody, err)
		return c.SendStatus(http.StatusBadGateway)
	}
	fmt.Printf("Got response: %+v\n", resp)

	return c.Status(http.StatusOK).JSON(resp)
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

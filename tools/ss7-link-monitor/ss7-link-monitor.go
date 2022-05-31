package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/hpcloud/tail"
)

var endpointURL string

func args2config() (tail.Config, int64) {
	config := tail.Config{Follow: true}
	n := int64(0)
	maxlinesize := int(0)
	flag.Int64Var(&n, "n", 0, "tail from the last Nth location")
	flag.IntVar(&maxlinesize, "max", 0, "max line size")
	flag.BoolVar(&config.Follow, "f", false, "wait for additional data to be appended to the file")
	flag.BoolVar(&config.ReOpen, "F", false, "follow, and track file rename/rotation")
	flag.BoolVar(&config.Poll, "p", false, "use polling, instead of inotify")
	flag.StringVar(&endpointURL, "u", "", "URL to trigger callback to")
	flag.Parse()
	if config.ReOpen {
		config.Follow = true
	}
	config.MaxLineSize = maxlinesize
	return config, n
}

func main() {
	config, n := args2config()
	if flag.NFlag() < 1 {
		fmt.Println("need one or more files as arguments")
		os.Exit(1)
	}

	if n != 0 {
		config.Location = &tail.SeekInfo{-n, io.SeekEnd}
	}

	done := make(chan bool)
	for _, filename := range flag.Args() {
		go tailFile(filename, config, done)
	}

	for _, _ = range flag.Args() {
		<-done
	}
}

func tailFile(filename string, config tail.Config, done chan bool) {
	defer func() { done <- true }()
	t, err := tail.TailFile(filename, config)
	if err != nil {
		fmt.Println(err)
		return
	}
	for line := range t.Lines {
		if strings.HasSuffix(line.Text, "Transition to=timeout") {
			log.Println(line.Text)
			if endpointURL != "" {
				shutDownIsupLinks()
			}
		}
	}
	err = t.Wait()
	if err != nil {
		fmt.Println(err)
	}
}

func shutDownIsupLinks() {
	agent := fiber.AcquireClient().Post(endpointURL)
	if err := agent.Parse(); err != nil {
		log.Fatal(err)
	}
	code, resp, err := agent.String()
	if err != nil {
		log.Fatal(err)
	}
	if code >= 400 {
		log.Fatalf("return error with code: %d", code)
	}
	log.Printf("OK : << %d | %s", code, resp)
}

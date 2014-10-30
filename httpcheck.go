package main

import (
	"bufio"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/artyom/autoflags"
)

type Config struct {
	List        string        `flag:"urls,file with urls, one url per line"`
	Delay       time.Duration `flag:"delay,delay between each check cycle"`
	NoKeppAlive bool          `flag:"nokeepalive,disable keep alive"`
	Concurrency int           `flag:"n,number of concurrent checks to make"`
}

func main() {
	config := Config{
		Delay:       30 * time.Second,
		Concurrency: 10,
	}
	if err := autoflags.Define(&config); err != nil {
		log.Fatal(err)
	}
	flag.Parse()
	if len(config.List) == 0 || config.Delay <= 0 || config.Concurrency <= 0 {
		flag.Usage()
		os.Exit(1)
	}
	f, err := os.Open(config.List)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	urls := make([]string, 0)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}
		urls = append(urls, line)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	f.Close()
	if len(urls) == 0 {
		log.Fatal("no urls to check")
	}
	client := http.DefaultClient
	client.Transport = http.DefaultTransport
	client.Transport.(*http.Transport).DisableKeepAlives = config.NoKeppAlive
	gate := make(chan struct{}, config.Concurrency)
	var wg sync.WaitGroup
	for {
		for _, url := range urls {
			gate <- struct{}{}
			wg.Add(1)
			go func(url string) {
				defer wg.Done()
				defer func() { <-gate }()
				begin := time.Now()
				resp, err := client.Head(url)
				if err != nil {
					log.Printf("ERROR:\t%s\t%s", url, err)
					return
				}
				log.Printf("%s\t%s\t%d", url, time.Now().Sub(begin), resp.StatusCode)
			}(url)
		}
		wg.Wait()
		time.Sleep(config.Delay)
	}
}

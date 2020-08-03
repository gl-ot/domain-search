package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"sync"
)

const (
	routineNum = 4
)

func main() {
	domain := flag.String("domain", "", "your desired domain")
	enrichingWordsPath := flag.String("words", "words.txt", "path to file with your enriching words")
	flag.Parse()
	if *domain == "" {
		log.Fatalf("Please specify your domain")
	}

	availableDomains := make(chan string)
	wg := sync.WaitGroup{}
	semaphore := make(chan struct{}, routineNum)
	for _, w := range readWords(*enrichingWordsPath) {
		wg.Add(1)
		cleanedWord := clean(w)
		go func() {
			defer wg.Done()
			semaphore <- struct{}{}
			for _, d := range computeDomains(*domain, cleanedWord) {
				body := getBody(d)
				if isAvailable(body) {
					availableDomains <- d
				}
			}
			<-semaphore
		}()
	}

	go func() {
		wg.Wait()
		close(availableDomains)
	}()

	for d := range availableDomains {
		fmt.Println(d)
	}
}

func clean(word string) string {
	word = strings.ReplaceAll(word, " ", "")
	word = strings.ToLower(word)
	return word
}

func computeDomains(domain, word string) []string {
	if word == "" {
		return []string{}
	}
	return []string{word + domain, domain + word}
}

func readWords(path string) []string {
	buf, err := ioutil.ReadFile(path)
	check(err)
	wordsStr := string(buf)
	return strings.Split(wordsStr, "\n")
}

func getBody(word string) []byte {
	url := fmt.Sprintf("https://uk.godaddy.com/domainfind/v1/search/exact?q=%s", word)
	res, err := http.Get(url)
	check(err)
	defer res.Body.Close()
	if res.StatusCode == http.StatusForbidden {
		log.Fatal("Godaddy limit request")
	}
	body, err := ioutil.ReadAll(res.Body)
	check(err)
	return body
}

func isAvailable(body []byte) bool {
	res := gjson.Get(string(body), "ExactMatchDomain.IsAvailable")
	return res.Bool()
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

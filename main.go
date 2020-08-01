package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func main() {
	domain := flag.String("domain", "", "your desired domain")
	enrichingWordsPath := flag.String("words", "words.txt", "path to file with your enriching words")
	resultPath := flag.String("result", "available.txt", "path to result file with available domains")
	flag.Parse()
	if *domain == "" {
		log.Fatalf("Please specify your domain")
	}

	_ = os.Remove(*resultPath)
	resultFile, err := os.OpenFile(*resultPath, os.O_WRONLY|os.O_CREATE, 0755)
	check(err)
	write := func(word string) {
		resultFile.Write([]byte(word))
		resultFile.WriteString("\n")
	}

	for _, rWord := range readWords(*enrichingWordsPath) {
		cleaned := clean(rWord)
		for _, domain := range domainNames(*domain, cleaned) {
			body := getBody(domain)
			if isAvailable(body) {
				write(domain)
			}
		}
	}
}

func clean(word string) string {
	word = strings.ReplaceAll(word, " ", "")
	word = strings.ToLower(word)
	return word
}

func domainNames(domain, word string) []string {
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
	url := fmt.Sprintf("https://uk.godaddy.com/domainfind/v1/search/exact?q=%s&key=ad_dlp_com_ru&pc=&ptl=&solution_set_ids=dpp-us-solution-tier1,dpp-intl-solution-tier4,dpp-intl-solution-tier6&itc=dlp_cheapdomain_com_ru&isc=rudomrub1&req_id=1596292396792", word)
	res, err := http.Get(url)
	check(err)
	defer res.Body.Close()
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

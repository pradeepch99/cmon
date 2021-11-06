package main

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/iangcarroll/cookiemonster/pkg/monster"
)

var (
	monsterWordlist string

	wl = monster.NewWordlist()
)

func init() {
	if err := wl.LoadFromString(monsterWordlist); err != nil {
		panic(err)
	}
}
func MonsterRun(cookie string) (success bool, err error) {
	c := monster.NewCookie(cookie)

	if !c.Decode() {
		return false, errors.New("could not decode")
	}

	if _, success := c.Unsign(wl, 100); !success {
		return false, errors.New("could not unsign")
	}

	return true, nil
}

func webIsReachable(url string) (resp *http.Response) {
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	response, err := http.Get(url)
	if err != nil {
		//fmt.Println(err)
		return
	}
	response.Body.Close()
	return response
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		if (!strings.HasPrefix(line, "http://")) && (!strings.HasPrefix(line, "https://")) {
			line = "https://" + line
		}
		//fmt.Println("[*]Checking domain", line)
		resp := webIsReachable(line)
		if resp != nil {
			fmt.Println("[*]Checking domain", line)
			for _, c := range resp.Cookies() {
				res, err := MonsterRun(c.Value)
				if res {
					fmt.Println("\t", c.Name, "=", c.Value)
					fmt.Println("\t", c.Name, "=>", res)
				} else {
					fmt.Println("\t", c.Name, "=>", err)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

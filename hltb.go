package main

import (
	"fmt"
	"github.com/moovweb/gokogiri"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type gameInfo struct {
	label string
	value string
}

const postUrl = "http://www.howlongtobeat.com/search_main.php?t=games&page=1&sorthead=&sortd=Normal&plat=&detail=0"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Should pass on the game name.")
		return
	}

	fmt.Println("hltb!")
	fmt.Println("======================")

	games := Scrap(strings.Join(os.Args[1:], " "))
	for _, game := range games {
		for _, info := range game {
			fmt.Println(info.label + ": " + info.value)
		}
		fmt.Println("======================")
	}
}

func Scrap(queryString string) [][]gameInfo {
	resp, _ := http.PostForm(postUrl, url.Values{"queryString": {queryString}})
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	doc, _ := gokogiri.ParseHtml(body)
	defer doc.Free()

	rootNode := doc.Root()
	gameNameNodes, _ := rootNode.Search("//h3/a")
	mainStoryNodes, _ := rootNode.Search("//div[text()=\"Main Story\"]")
	mainPlusExtraNodes, _ := rootNode.Search("//div[text()=\"Main + Extra\"]")
	completionistNodes, _ := rootNode.Search("//div[text()=\"Completionist\"]")
	combinedNodes, _ := rootNode.Search("//div[text()=\"Combined\"]")

	games := make([][]gameInfo, len(gameNameNodes))
	for i := 0; i < len(gameNameNodes); i++ {
		games[i] = []gameInfo{
			gameInfo{"Game Name", gameNameNodes[i].Content()},
			gameInfo{"Main Story", mainStoryNodes[i].NextSibling().Content()},
			gameInfo{"Main + Extra", mainPlusExtraNodes[i].NextSibling().Content()},
			gameInfo{"Completionist", completionistNodes[i].NextSibling().Content()},
			gameInfo{"Combined", combinedNodes[i].NextSibling().Content()},
		}
	}
	return games
}

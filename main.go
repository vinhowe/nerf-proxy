package main

import (
	"flag"
	"fmt"
	"github.com/elazarl/goproxy"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

const BlacklistFileName = "blacklist.txt"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	_, err := os.Stat(BlacklistFileName)
	if os.IsNotExist(err) {
		fmt.Printf("%s does not exist, creating it", BlacklistFileName)
		_, err = os.Create(BlacklistFileName)
		check(err)
	}

	data, err := ioutil.ReadFile(BlacklistFileName)
	check(err)

	blockList := strings.Split(string(data), "\n")
	blockList = blockList[0:len(blockList)-1]
	blockListFmt := strings.ReplaceAll(strings.Join(blockList, "|"), ".", "\\.")
	blockListRegex := regexp.MustCompile(fmt.Sprintf(".*%s/*", blockListFmt))
	fmt.Printf("Loaded blacklist with %d items", len(blockList))

	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = false

	proxy.OnRequest(goproxy.ReqHostMatches(blockListRegex)).
		HandleConnect(goproxy.AlwaysReject)

	addr := flag.String("addr", ":8080", "proxy listen address")
	flag.Parse()

	log.Fatal(http.ListenAndServe(*addr, proxy))
}

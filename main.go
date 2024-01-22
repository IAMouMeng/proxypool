package main

import (
	"flag"
	"fmt"
	"proxy/data"
	"proxy/fofa"
	"proxy/log"
	"proxy/proxy"
	"proxy/web"
)

var (
	email     string
	key       string
	ruleFile  string
	page      int
	size      int
	thread    int
	timeout   int
	checkTime int
	getTime   int
	testUrl   string
	apiHost   string
	apiPort   int
)

func init() {
	flag.StringVar(&email, "e", "", "Search API Email")
	flag.StringVar(&ruleFile, "r", "rule.txt", "Search rule file")
	flag.StringVar(&key, "k", "", "Search API Key")
	flag.IntVar(&page, "g", 1, "Search data pages")
	flag.IntVar(&size, "z", 10000, "Search data size")
	flag.IntVar(&thread, "d", 1000, "Check proxy thread")
	flag.StringVar(&testUrl, "u", "http://google.com", "Test url")
	flag.IntVar(&timeout, "t", 10, "Proxy timeout (seconds)")
	flag.IntVar(&checkTime, "ct", 7200, "Check database proxy time (seconds)")
	flag.IntVar(&getTime, "gt", 1800, "Get proxy time (seconds)")
	flag.StringVar(&apiHost, "host", "127.0.0.1", "Web API host")
	flag.IntVar(&apiPort, "port", 3511, "Web API port")

	flag.Usage = usage
}

func init() {
	data.InitDb("./database.db")
}

func main() {

	flag.Parse()

	if email != "" && key != "" {
		client := fofa.NewClient(email, key)

		info, err := client.GetAccountInfo()

		if err != nil {
			log.LogError("Error getting account info:", err)
			return
		}

		logAccountInfo(info)

		proxy.CreateGetTask(getTime, client, ruleFile, page, size, thread, testUrl, timeout)
		proxy.CreateCheckTask(checkTime, testUrl, timeout)
		web.InitApi(apiHost, apiPort)
	} else {
		flag.Usage()
	}
}

func logAccountInfo(info map[string]interface{}) {
	log.LogInfo(fmt.Sprintf("Login Success!! Userinfo: Username: %v VIP: %t VIP Level: %.0f Remaining API Queries: %.0f Remaining API Data: %.0f",
		info["username"],
		info["isvip"],
		info["vip_level"],
		info["remain_api_query"],
		info["remain_api_data"],
	))
}

func usage() {
	fmt.Println(`Welcome to ProxyPool system
Usage: ./proxy [-e email] [-k key] [...]
Options:`)
	flag.PrintDefaults()
}

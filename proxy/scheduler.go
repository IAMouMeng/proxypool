package proxy

import (
	"fmt"
	"os"
	"proxy/data"
	"proxy/fofa"
	"proxy/log"
	"strings"
	"time"
)

type Task struct {
	Name string
	Func func()
}

type Scheduler struct {
	tasks []Task
}

func NewScheduler() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) AddTask(task Task) {
	s.tasks = append(s.tasks, task)
}

func (s *Scheduler) Start(interval time.Duration) {
	s.runTasks()

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.runTasks()
		}
	}
}

func (s *Scheduler) runTasks() {
	for _, task := range s.tasks {
		fmt.Printf("Running task %s\n", task.Name)
		task.Func()
	}
}

func CreateGetTask(second int, client *fofa.Client, ruleFile string, page int, size int, thread int, testUrl string, timeout int) {
	scheduler := NewScheduler()

	scheduler.AddTask(Task{
		Name: "GetProxyTask",
		Func: func() {
			engine := NewEngine(thread)
			engine.Run()

			fileContent, err := os.ReadFile(ruleFile)

			if err != nil {
				log.LogError("无法读取文件:", err)
				return
			}

			lines := strings.Split(string(fileContent), "\n")

			for _, rule := range lines {
				response, err := client.SearchData(rule, page, size)
				if err != nil {
					log.LogError("Error searching data:", err)
					return
				}

				for _, dataResult := range response["results"].([]interface{}) {
					engine.SubmitTask(ProxyData{TestUrl: testUrl, TimeOut: timeout, Data: dataResult.([]interface{})})
				}
				log.LogInfo(fmt.Sprintf("Search Success!! Query:%s Size:%d", response["query"], len(response["results"].([]interface{}))))
			}
		},
	})

	go scheduler.Start(time.Duration(second) * time.Second)
}

func CreateCheckTask(second int, testUrl string, timeout int) {
	scheduler := NewScheduler()

	scheduler.AddTask(Task{
		Name: "CheckProxyTask",
		Func: func() {
			proxyList, _ := data.GlobalEngine.GetProxyList("", "")

			for _, proxyData := range proxyList {
				res := CheckProxyValidity(proxyData.Proxy, time.Duration(timeout), testUrl) // stop use worker

				if !res {
					err := data.GlobalEngine.DeleteProxy(proxyData.Proxy)
					if err != nil {
						log.LogError("database error", err)
					}

					log.LogInfo(fmt.Sprintf("From Tasker: Database proxy %s has been deleted", proxyData.Proxy))
				}
			}
		},
	})

	go scheduler.Start(time.Duration(second) * time.Second)
}

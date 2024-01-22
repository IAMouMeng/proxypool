package proxy

import (
	"fmt"
	"proxy/data"
	"proxy/log"
	"sync"
	"time"
)

type Engine struct {
	Wg      *sync.WaitGroup
	Tasks   chan ProxyData
	Result  chan string
	Threads int
	Jobs    int
}

func NewEngine(threads int) *Engine {
	return &Engine{
		Wg:      &sync.WaitGroup{},
		Tasks:   make(chan ProxyData, 100000),
		Threads: threads,
		Jobs:    0,
	}
}

func (e *Engine) Run() {
	go e.Scheduler()
}

func (e *Engine) Scheduler() {
	for i := 0; i < e.Threads; i++ {
		e.worker(e.Tasks)
	}
}

func (e *Engine) SubmitTask(data ProxyData) {
	go func() {
		e.Tasks <- data
	}()
}

func (e *Engine) worker(tasks chan ProxyData) {
	go func() {
		for dataProxy := range tasks {
			proxyAddress := fmt.Sprintf("%s:%s", dataProxy.Data.([]interface{})[0].(string), dataProxy.Data.([]interface{})[1].(string))

			res := CheckProxyValidity(proxyAddress, time.Duration(dataProxy.TimeOut), dataProxy.TestUrl)
			if res {
				log.LogInfo(fmt.Sprintf("Found available proxy: %s Country: %s", proxyAddress, dataProxy.Data.([]interface{})[2].(string)))

				err := data.GlobalEngine.InsertProxy(data.ProxyData{Proxy: proxyAddress, Country: dataProxy.Data.([]interface{})[2].(string), Type: "http"})

				if err != nil {
					log.LogError("database error:", err)
				}
			}
		}
	}()
}

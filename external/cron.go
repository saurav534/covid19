package external

import (
	"log"
	"time"
)

type Cron struct {
	done chan bool
}

func NewCron() *Cron {
	return &Cron{
		done: make(chan bool),
	}
}

func (c *Cron) StartCron() {
	go func() {
		tickers := make([]*time.Ticker, 0)
		t1 := time.NewTicker(309 * time.Second)
		t2 := time.NewTicker(327 * time.Second)
		tickers = append(tickers, t1, t2)
		for {
			select {
			case <-c.done:
				for _, ticker := range tickers {
					ticker.Stop()
				}
				log.Println("All ticker stopped")
				return
			case <-t1.C:
				UpdateCrowdData()
			case <-t2.C:
				UpdateMohData()
			}

		}
	}()
}

func (c *Cron) StopCron() {
	c.done <- true
}

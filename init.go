package loadbalancer

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func UpdateHealthCron(pool *ServerPool) {
	c := cron.New()
	c.AddFunc("@every 5s", func() {
		fmt.Println("running cron after 5 sec")
		ctx := context.Background()
		for _ , server := pool.Servers{
			err := GetHealth(ctx,server.URL)
			if err != nil {
				fmt.Printf("server is down %s\n",server.URL.Host)
				server.SetHealth(true)
			} else{
				server.SetHealth(false)
			}
		}
	})

	c.Start()
}

func GetHealth(ctx context.Context, url *url.URL) error {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", fmt.Sprintf("http://"+url.Host+"/health"), nil)
	resp, err := client.Do(request)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("Server is down")
	}
	return nil
}

func RoundRobinScheduler(pool *ServerPool) *Server {
	var counts int
	for i:= pool.GetNextAvailableServer();;i++{
		if counts> int(pool.ServerCount){
			return nil
		}
		idx := (i)%pool.ServerCount
		fmt.Println("next available server is :",idx)
		if !pool.Servers[idx].GetHealth(){
			pool.UpdateLastServerUsed(idx)
			return pool.Servers[idx]
		}
		counts++
	}
}

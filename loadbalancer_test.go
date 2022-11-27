package loadbalancer

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var serverPool *ServerPool

func initServerPool() {
	serverPool = GetNewServerPool()
	server1 := GetNewServer("https://localhost:9000")
	server2 := GetNewServer("http://localhost:8000")
	serverPool.RegisterServer(server1)
	serverPool.RegisterServer(server2)
}

func PingStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hey I am Alive",
	})
}

func LoadBalancer(c *gin.Context) {
	server := RoundRobinScheduler(serverPool)
	fmt.Println("server is: ", server)
	server.ReverseProxy.ServeHTTP(c.Writer, c.Request)
}

func initRoutes(router *gin.Engine) {
	router.GET("/ping", PingStatus)
	router.Any("/api/*pattern", LoadBalancer)
}

func main() {
	router := gin.Default()
	initRoutes(router)
	initServerPool()
	UpdateHealthCron(serverPool)
	router.Run(":8800")
}

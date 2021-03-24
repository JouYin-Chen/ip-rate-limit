package middleware

import (
	"ip-rate-limit/helper"
	"net/http"
	"fmt"
	"net"
	"strconv"
	"github.com/gin-gonic/gin"
)
const LimitationTimes int = 40
const TimeOutSeconds int = 60

func RateLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientIP, err := getClientIPByRequest(c.Request)
		if err != nil {
			fmt.Println("Getting req.RemoteAddr", err)
			RespondWithError(400, "Getting req.RemoteAddr error", c)
		} else {
			rateLimitingByIp(c, clientIP)
			IPAddressTracking(clientIP)
			c.Next()
		}
	}
}

func RespondWithError(code int, message string, c *gin.Context) {
	resp := map[string]string{"ErrorMsg": message}

	c.JSON(code, resp)
	c.Abort()
}

// getClientIPByRequest tries to get directly from the Request.
func getClientIPByRequest(req *http.Request) (ip string, err error) {
	IPAddress := req.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		fmt.Println("get ip from header: X-Forwarded-For")
		IPAddress = req.Header.Get("X-Forwarded-For")
	}

	if IPAddress == "" {
	// Try via request
		remoteAddr, port, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
				return "", err
		} else {
			IPAddress = remoteAddr
			fmt.Println("With req.RemoteAddr found IP:%v; Port: %v", remoteAddr, port)
		}
	}

	if IPAddress == "" {
		netIP := net.ParseIP(ip)
		fmt.Println("get ip from net", string(netIP))
		if string(netIP) != "" {
			IPAddress = string(netIP)
		}
	}

	if IPAddress == "" {
		message := fmt.Sprintf("Parsing IP from Request.RemoteAddr got nothing.")
		// fmt.Println(message)
		return "", fmt.Errorf(message)
	}
	return IPAddress, nil
}
func IPAddressTracking(IP string) {
	var ipKey = fmt.Sprintf("IP:%v", IP)
	helper.SetValueExpByKey(ipKey, "0")
	helper.SetExpireKey(ipKey, TimeOutSeconds)
	helper.IncrNumberByKey(ipKey)
}

func rateLimitingByIp(c *gin.Context, IP string) {
	var ipKey = fmt.Sprintf("IP:%v", IP)

	value, _ := helper.FindValueByKey(ipKey)
	requests, _ := strconv.Atoi(string(value))
	if requests + 1 > LimitationTimes {
		errMsg := fmt.Sprintf("Too Many Requests, IP address %v have %v requests", IP, requests + 1)

		RespondWithError(429, errMsg, c)
	}

}
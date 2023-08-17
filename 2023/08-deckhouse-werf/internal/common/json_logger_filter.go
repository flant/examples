package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func JsonLogger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(
		func(params gin.LogFormatterParams) string {
			log := make(map[string]interface{})

			log["status_code"] = params.StatusCode
			log["path"] = params.Path
			log["method"] = params.Method
			log["start_time"] = params.TimeStamp.Format("2023/01/02 - 13:04:05")
			log["remote_addr"] = params.ClientIP
			log["response_time"] = params.Latency.String()

			str, _ := json.Marshal(log)
			return string(str) + "\n"
		},
	)
}

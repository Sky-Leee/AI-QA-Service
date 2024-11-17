package coze

import (
	"ai-qa-service/internal/logger"
	"time"

	"github.com/spf13/viper"
)

var (
	api   string
	host  string
	token string

	genBotID   string
	judgeBotID string

	user string

	timeout time.Duration
)

func InitCoze() {
	api = viper.GetString("coze.api")
	host = viper.GetString("coze.host")
	token = viper.GetString("coze.token")

	genBotID = viper.GetString("coze.bot_id.gen")
	judgeBotID = viper.GetString("coze.bot_id.judge")

	user = viper.GetString("coze.user")

	timeout = time.Duration(viper.GetInt("coze.timeout")) * time.Second

	logger.Infof("Coze initialized: api=%s, host=%s, token=%s, genBotID=%s, judgeBotID=%s, user=%s, timeout=%s", api, host, token, genBotID, judgeBotID, user, timeout)
}

func getArgs(botID, query string) *RequestBody {
	return &RequestBody{
		Method: "POST",
		URL:    api,
		Host:   host,
		Token:  token,
		Body: CozeRequestBody{
			BotID:  botID,
			User:   user,
			Query:  query,
			Stream: false,
		},
	}
}

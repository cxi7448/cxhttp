package xtelegram

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var mTgRobot *tgbotapi.BotAPI

var sChatId = int64(0)
var sToken string

// 创建
func InitRobot(chatId int64, token string) {
	sChatId = chatId
	sToken = token
	TgRobot, tgErr := tgbotapi.NewBotAPI(sToken)
	if tgErr != nil {
		clLog.Error("创建Telegram机器人失败! 错误: %v", tgErr)
		return
	}
	mTgRobot = TgRobot
}

// 发送消息到telegram
func SendMessage(_message string, args ...interface{}) error {
	if mTgRobot == nil {
		clLog.Error("telegram 机器人初始化失败")
		return fmt.Errorf("telegram 机器人初始化失败")
	}
	var sendMsg = _message
	if args != nil && len(args) > 0 {
		sendMsg = fmt.Sprintf(sendMsg, args...)
	}

	msg := tgbotapi.NewMessage(sChatId, sendMsg)
	_, err := mTgRobot.Send(msg)
	if err != nil {
		clLog.Error("发送Telegram消息失败![%v]-[%v]", sChatId, sToken)
		return err
	}
	return nil
}

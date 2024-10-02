package xtelegram

import (
	"fmt"
	"github.com/cxi7448/cxhttp/clUtil/clLog"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

var mTgRobot *XTRobot

// 创建
func InitRobot(chatId int64, token string) {
	mTgRobot = New(chatId, token)
}

// 发送消息到telegram
func SendMessage(_message string, args ...interface{}) error {
	return mTgRobot.SendMessage(_message, args...)
}

func New(chatid int64, token string) *XTRobot {
	return &XTRobot{
		ChatID: chatid,
		Token:  token,
	}
}

func (this *XTRobot) init() {
	if this.botApi != nil {
		return
	}
	robot, err := tgbotapi.NewBotAPI(this.Token)
	if err != nil {
		clLog.Error("创建Telegram机器人失败! 错误: %v", err)
		return
	}
	this.botApi = robot
}

func (this *XTRobot) SendMessage(_message string, args ...interface{}) error {
	this.init()
	if this.botApi == nil {
		clLog.Error("telegram 机器人初始化失败")
		return fmt.Errorf("telegram 机器人初始化失败")
	}
	var sendMsg = _message
	if args != nil && len(args) > 0 {
		sendMsg = fmt.Sprintf(sendMsg, args...)
	}

	msg := tgbotapi.NewMessage(this.ChatID, sendMsg)
	_, err := this.botApi.Send(msg)
	if err != nil {
		clLog.Error("发送Telegram消息失败![%v]-[%v]", this.ChatID, this.Token)
		return err
	}
	return nil
}

type XTRobot struct {
	ChatID int64
	Token  string
	botApi *tgbotapi.BotAPI
}

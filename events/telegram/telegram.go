package telegram

import "Telegram_Bot/claims/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	// storage
}

func New(client *telegram.Client, offset int) *Processor {}

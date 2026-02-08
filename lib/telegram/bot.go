package telegram

import (
	"fmt"
	"io"
	"net/http"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot    *tgbotapi.BotAPI
	userId int64
}

func NewTelegramBot(apiKey string, userId int64) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		return nil, fmt.Errorf("Error while creating new Telegram bot : %v", err)
	}

	bot.Debug = false

	return &TelegramBot{
		bot:    bot,
		userId: userId,
	}, nil
}

// blocks the caller, so call in  a separate goroutine
func (t *TelegramBot) Listen() <-chan tgbotapi.Update {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 30
	return t.bot.GetUpdatesChan(updateConfig)
}

// save the image sent to the bot to a local file
func (t *TelegramBot) SaveTelegramImage(update tgbotapi.Update, path string) error {
	photo := update.Message.Photo[len(update.Message.Photo)-1]
	file, err := t.bot.GetFile(tgbotapi.FileConfig{
		FileID: photo.FileID,
	})
	if err != nil {
		return fmt.Errorf("SaveTelegramImage GetFile failed: %v", err)
	}
	fileURL := file.Link(t.bot.Token)

	resp, err := http.Get(fileURL)
	if err != nil {
		return fmt.Errorf("SaveTelegramImage download failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("SaveTelegramImage bad http status: %s", resp.Status)
	}

	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("SaveTelegramImage create file failed: %v", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("SaveTelegramImage save file failed: %v", err)
	}

	return nil
}

func (t *TelegramBot) StopListening() {
	t.bot.StopReceivingUpdates()
}

func (t *TelegramBot) SendMessage(message string) error {
	msg := tgbotapi.NewMessage(t.userId, message)

	if _, err := t.bot.Send(msg); err != nil {
		return err
	}
	return nil
}

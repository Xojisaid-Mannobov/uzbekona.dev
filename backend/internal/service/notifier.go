package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"uzbekona.dev/backend/internal/config"
	"uzbekona.dev/backend/internal/model"
)

// Notifier — yangi so'rov kelganda jamoani xabardor qiladi.
type Notifier interface {
	NewContact(c model.Contact)
	NewApplication(a model.Application)
}

type noopNotifier struct{}

func (noopNotifier) NewContact(model.Contact)         {}
func (noopNotifier) NewApplication(model.Application) {}

// telegramNotifier — TELEGRAM_BOT_TOKEN va TELEGRAM_CHAT_ID berilsa, so'rovni Telegram'ga yuboradi.
type telegramNotifier struct {
	token  string
	chatID string
	client *http.Client
}

func NewTelegramNotifier(cfg *config.Config) Notifier {
	if cfg.TelegramBotToken == "" || cfg.TelegramChatID == "" {
		return noopNotifier{}
	}
	return &telegramNotifier{
		token:  cfg.TelegramBotToken,
		chatID: cfg.TelegramChatID,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (t *telegramNotifier) NewContact(c model.Contact) {
	var b strings.Builder
	fmt.Fprintf(&b, "<b>Yangi so‘rov #%d</b>\n\n", c.ID)
	line := func(label, value string) {
		if value != "" {
			fmt.Fprintf(&b, "<b>%s:</b> %s\n", label, html.EscapeString(value))
		}
	}
	line("Ism", c.Name)
	line("Aloqa", c.Contact)
	line("Email", c.Email)
	line("Loyiha turi", c.ProjectType)
	line("Budjet", c.Budget)
	fmt.Fprintf(&b, "\n%s", html.EscapeString(truncate(c.Message, 3000)))
	t.send(b.String())
}

// NewApplication — "Jamoaga qo'shilish" arizasi haqida xabar.
func (t *telegramNotifier) NewApplication(a model.Application) {
	var b strings.Builder
	fmt.Fprintf(&b, "<b>Yangi nomzod #%d — %s</b>\n\n", a.ID, html.EscapeString(a.Position))
	line := func(label, value string) {
		if value != "" {
			fmt.Fprintf(&b, "<b>%s:</b> %s\n", label, html.EscapeString(value))
		}
	}
	line("Ism", a.FullName)
	line("Email", a.Email)
	line("Telefon", a.Phone)
	line("Telegram", a.Telegram)
	line("Tajriba", a.Experience)
	line("Portfolio", a.PortfolioURL)
	line("Rezyume", a.ResumeURL)
	fmt.Fprintf(&b, "\n%s", html.EscapeString(truncate(a.About, 3000)))
	t.send(b.String())
}

func (t *telegramNotifier) send(text string) {
	body, _ := json.Marshal(map[string]any{
		"chat_id":    t.chatID,
		"text":       text,
		"parse_mode": "HTML",
	})
	resp, err := t.client.Post("https://api.telegram.org/bot"+t.token+"/sendMessage", "application/json", bytes.NewReader(body))
	if err != nil {
		slog.Warn("telegram bildirishnoma yuborilmadi", "error", err)
		return
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		slog.Warn("telegram bildirishnoma rad etildi", "status", resp.StatusCode)
	}
}

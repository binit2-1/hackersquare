package bot

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/binit2-1/hackersquare/apps/api/internal/repository/pg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InitTelegramBot() (*tgbotapi.BotAPI, error) {
	token := os.Getenv("TELEGRAM_BOT_TOKEN")

	fmt.Printf("🚨 DEBUG: Token string length is: %d\n", len(token))
    if len(token) > 5 {
        fmt.Printf("🚨 DEBUG: Token starts with: %s\n", token[:5])
        fmt.Printf("🚨 DEBUG: Token ends with: %s\n", token[len(token)-5:])
    }
	
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN environment variable is empty")
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Telegram: %w", err)
	}

	bot.Debug = true
	log.Printf("Successfully authorized on Telegram as @%s", bot.Self.UserName)

	return bot, nil
}

// RunTelegramListener listens for commands and saves user preferences
func RunTelegramListener(bot *tgbotapi.BotAPI, repo *pg.PostgresEventRepo) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.GetUpdatesChan(u)

	fmt.Println("🤖 Telegram Listener started... Waiting for commands.")

	for update := range updates {
		// Only look at messages that are bot commands (start with "/")
		if update.Message == nil || !update.Message.IsCommand() {
			continue
		}

		chatID := fmt.Sprintf("%d", update.Message.Chat.ID)
		command := update.Message.Command()
		args := update.Message.CommandArguments()

		// Handle the /start welcome message
		if command == "start" {
			welcomeMsg := "Welcome to HackerSquare Alerts! 🚀\n\nTo get personalized hackathon pings, tell me your location and tech stack using the /setup command.\n\n*Format:* `/setup Country | Tag1, Tag2`\n*Example:* `/setup India | React, Go`"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcomeMsg)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
			continue
		}

		// Handle the /setup preference saving
		if command == "setup" {
			// Split the arguments by the pipe "|" character
			parts := strings.Split(args, "|")
			if len(parts) != 2 {
				errorMsg := "⚠️ Invalid format. Please use a pipe `|` to separate country and tags.\n\n*Example:* `/setup India | Next.js, Supabase`"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorMsg)
				msg.ParseMode = "Markdown"
				bot.Send(msg)
				continue
			}

			// Clean up the country string
			country := strings.TrimSpace(parts[0])

			// Clean up the tags array
			rawTags := strings.Split(parts[1], ",")
			var tags []string
			for _, t := range rawTags {
				cleanTag := strings.TrimSpace(t)
				if cleanTag != "" {
					tags = append(tags, cleanTag)
				}
			}

			// Save it to the database!
			err := repo.UpsertSubscription(context.Background(), "telegram", chatID, tags, country)
			if err != nil {
				log.Printf("DB Error saving subscription for %s: %v", chatID, err)
				bot.Send(tgbotapi.NewMessage(update.Message.Chat.ID, "❌ Server error: Failed to save preferences. Please try again later."))
				continue
			}

			// Confirm success with the user
			reply := fmt.Sprintf("✅ Preferences saved successfully!\n\nI will now ping this chat whenever a new hackathon in *%s* matching `%v` is added.", country, tags)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			msg.ParseMode = "Markdown"
			bot.Send(msg)
		}
	}
}

func StartHackathonNotifier(bot *tgbotapi.BotAPI, repo *pg.PostgresEventRepo) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	lastCheck := time.Now()

	fmt.Println("⏰ Telegram notifier ticker has started...")

	for range ticker.C {
		newHackathons, err := repo.GetNewHackathonsSince(context.Background(), lastCheck)
		if err != nil {
			log.Printf("Notifier DB Error: %v\n", err)
			continue
		}

		if len(newHackathons) > 0 {
			lastCheck = time.Now()
		}

		// Loop through the new hackathons and route them. 'h' is defined here!
		for _, h := range newHackathons {
			subscribers, err := repo.GetMatchingChats(context.Background(), h.Location, []string{})
			if err != nil || len(subscribers) == 0 {
				continue // Nobody matched this one, skip to the next hackathon
			}

			// Note: If PrizeUSD is a pointer in your struct, use *h.PrizeUSD here
			msgText := fmt.Sprintf(
				"🚀 *New Hackathon Alert!*\n\n* %s *\n📍 %s\n💰 Prize Pool: $%.2f\n📅 Starts: %s\n🔗 [Apply Here](%s)",
				h.Title,
				h.Location,
				*h.PrizeUSD, // Assuming it's a pointer based on your previous code
				h.StartDate.Format("Jan 02, 2006"),
				h.ApplyURL,
			)

			for _, chatIDStr := range subscribers {
				chatID, parseErr := strconv.ParseInt(chatIDStr, 10, 64)
				if parseErr != nil {
					continue
				}

				msg := tgbotapi.NewMessage(chatID, msgText)
				msg.ParseMode = "Markdown"
				msg.DisableWebPagePreview = true

				if _, err := bot.Send(msg); err != nil {
					log.Printf("Failed to send to %s: %v\n", chatIDStr, err)
				} else {
					fmt.Printf("✅ Sent targeted alert for '%s' to ChatID: %s\n", h.Title, chatIDStr)
				}
			}
		}
	}
}

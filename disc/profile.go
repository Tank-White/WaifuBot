package disc

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Karitham/WaifuBot/database"
	"github.com/Karitham/WaifuBot/query"
	"github.com/diamondburned/arikawa/v2/bot/extras/arguments"
	"github.com/diamondburned/arikawa/v2/discord"
	"github.com/diamondburned/arikawa/v2/gateway"
	"go.mongodb.org/mongo-driver/mongo"
)

// Name represent the name of a character
type Name = string

// Profile displays user profile
func (b *Bot) Profile(m *gateway.MessageCreateEvent, _ ...*arguments.UserMention) (*discord.Embed, error) {
	user := parseUser(m)

	uData, err := database.ViewUserData(user.ID)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	return &discord.Embed{
		Title: fmt.Sprintf("%s's profile", user.Username),
		Description: fmt.Sprintf(
			"%s\n%s last rolled %s ago.\nThey have rolled %d waifus and claimed %d.\nFavorite waifu is %s",
			uData.Quote, user.Username, time.Since(uData.Date).Truncate(time.Minute), len(uData.Waifus), uData.ClaimedWaifus, uData.Favorite.Name,
		),
		Thumbnail: &discord.EmbedThumbnail{URL: uData.Favorite.Image},
	}, nil
}

// Favorite sets a waifu as favorite
func (b *Bot) Favorite(m *gateway.MessageCreateEvent, name ...Name) (string, error) {
	if name == nil {
		return "", errors.New("no character name entered")
	}
	n := strings.Join(name, " ")

	id := parseArgs(n)
	searchArgs := query.CharSearchInput{
		ID:   id,
		Name: n,
	}

	char, err := query.CharSearch(searchArgs)
	if err != nil {
		return "", err
	}

	err = database.Favorite(char).Set(m.Author.ID)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("New waifu set, check your profile\n<%s>", char.Character.SiteURL), nil
}

// Quote sets a quote on the user profile
func (b *Bot) Quote(m *gateway.MessageCreateEvent, quote ...string) (string, error) {
	if quote == nil {
		return "", errors.New("no quote entered")
	}

	q := strings.Join(quote, " ")
	return fmt.Sprintf("New quote set :\n%s", q), database.Quote(q).Set(m.Author.ID)
}

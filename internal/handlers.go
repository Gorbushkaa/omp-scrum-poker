package handlers

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strconv"
	"strings"
)

func MessageHandler(s *discordgo.Session, m *discordgo.MessageCreate, taskUrl *string, storyPoints map[string]string) error {

	prefix := fmt.Sprintf("<@%s> ", s.State.User.ID)

	if strings.HasPrefix(m.Content, prefix) {
		*taskUrl = "https://jira.omprussia.ru/" + strings.Replace(m.Content, prefix, "", 1)
		storyPoints = make(map[string]string)
		btn := discordgo.Button{
			Label:    "Начать Scrum Poker",
			Style:    discordgo.PrimaryButton,
			CustomID: "start",
		}

		actions := discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{btn},
		}

		embed := &discordgo.MessageEmbed{
			Title:       "OMP Scrum Poker",
			Description: *taskUrl,
		}

		data := &discordgo.MessageSend{
			Components: []discordgo.MessageComponent{actions},
			Embed:      embed,
		}

		_, err := s.ChannelMessageSendComplex(m.ChannelID, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func InteractionHandler(s *discordgo.Session, i *discordgo.InteractionCreate, taskUrl *string, storyPoints map[string]string) error {

	var fibNum = [...]int{1, 2, 3, 5, 8, 13, 21, 34, 55, 89}
	var fibMenuOptions []discordgo.SelectMenuOption
	for _, num := range fibNum {
		option := discordgo.SelectMenuOption{
			Label:   strconv.Itoa(num),
			Value:   strconv.Itoa(num),
			Default: false,
		}
		fibMenuOptions = append(fibMenuOptions, option)
	}

	fullSelect := discordgo.SelectMenu{
		CustomID:    "storypoints",
		Placeholder: "Сколько StoryPoint'ов поставишь?",
		MaxValues:   1,
		Options:     fibMenuOptions,
		Disabled:    false,
	}

	btn := discordgo.Button{
		Label:    "Открыть карты",
		Style:    discordgo.DangerButton,
		CustomID: "open",
	}

	selectActions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{fullSelect},
	}

	buttonActions := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{btn},
	}

	embed := discordgo.MessageEmbed{
		Title:       "OMP Scrum Poker",
		Description: *taskUrl,
	}

	embeds := []*discordgo.MessageEmbed{
		&embed,
	}

	messageEdit := discordgo.MessageEdit{
		Components: []discordgo.MessageComponent{selectActions, buttonActions},
		Embeds:     embeds,
		ID:         i.Message.ID,
		Channel:    i.ChannelID,
	}

	switch i.MessageComponentData().CustomID {

	case "start":
		_, err := s.ChannelMessageEditComplex(&messageEdit)
		if err != nil {
			return err
		}

	case "storypoints":
		storyPoints[i.Member.User.Mention()] = i.MessageComponentData().Values[0]

		messageEdit.SetContent(GetVotedUsers(storyPoints, true))
		_, err := s.ChannelMessageEditComplex(&messageEdit)
		if err != nil {
			return err
		}

	case "open":
		messageEdit := discordgo.MessageEdit{
			Components: []discordgo.MessageComponent{},
			Embeds:     embeds,
			ID:         i.Message.ID,
			Channel:    i.ChannelID,
		}

		messageEdit.SetContent(GetVotedUsers(storyPoints, false))
		_, err := s.ChannelMessageEditComplex(&messageEdit)
		if err != nil {
			return err
		}
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
	})
	if err != nil {
		return err
	}
	return nil
}

func GetVotedUsers(usersList map[string]string, closed bool) string {
	var votedUsers string

	for userName, storyPoint := range usersList {
		switch closed {
		case true:
			votedUsers += fmt.Sprintf("%s 🦦, ", userName)
		case false:
			votedUsers += fmt.Sprintf("%s %s, ", userName, storyPoint)
		}
	}

	return votedUsers
}

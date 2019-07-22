package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type DiscordService struct {
	client        *ChatPlugClient
	discordClient *discordgo.Session
}

type DiscordServiceConfiguration struct {
	BotToken string `json:"botToken"`
}

func (ds *DiscordService) Startup(args []string) {
	ds.client = &ChatPlugClient{}
	fmt.Println("serviceID: " + args[1])
	ds.client.Connect(args[1], "http://localhost:2137/query", "ws://localhost:2137/query")

	if !ds.IsConfigured() {
		config := ds.client.AwaitConfiguration(ds.GetConfigurationSchema())
		ds.SaveConfiguration(config.FieldValues)
	}

	serviceConfiguration, err := ds.GetConfiguration()

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("doopa22: " + serviceConfiguration.BotToken)
	ds.discordClient, err = discordgo.New("Bot " + serviceConfiguration.BotToken)
	ds.discordClient.AddHandler(ds.discordMessageCreate)

	ds.discordClient.Open()
	msgChan := ds.client.SubscribeToNewMessages()
	defer ds.client.Close()

	for msg := range msgChan {
		webhooks, _ := ds.discordClient.ChannelWebhooks(msg.TargetThreadID)

		hasWebhook := false
		var webhook *discordgo.Webhook

		for _, hook := range webhooks {
			if strings.HasPrefix(hook.Name, "ChatPlug ") {
				hasWebhook = true
				webhook = hook
			}
		}

		if !hasWebhook {
			channel, _ := ds.discordClient.Channel(msg.TargetThreadID)
			webhook, _ = ds.discordClient.WebhookCreate(msg.TargetThreadID, "ChatPlug "+channel.Name, "https://i.imgur.com/l2QP9Go.png")
		}

		data := &discordgo.WebhookParams{
			Content:   msg.Message.Body,
			Username:  msg.Message.Author.Username,
			AvatarURL: msg.Message.Author.AvatarURL,
		}

		ds.discordClient.WebhookExecute(webhook.ID, webhook.Token, true, data)

		for _, attachment := range msg.Message.Attachments {
			fmt.Printf("doopsko haha yes %s\n", attachment.SourceURL)

			data := &discordgo.WebhookParams{
				Username:  msg.Message.Author.Username,
				AvatarURL: msg.Message.Author.AvatarURL,
				// File:      "https://i.imgur.com/ZGPxFN2.jpg",
			}

			// url := fmt.Sprintf("https://discordapp.com/webhooks/%s/%s", webhook.ID, webhook.Token)

			err := ds.discordClient.WebhookExecute(webhook.ID, webhook.Token, true, data)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}

func (ds *DiscordService) discordMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	webhook, err := ds.discordClient.Webhook(m.WebhookID)
	if err == nil && webhook != nil {
		if strings.HasPrefix(webhook.Name, "ChatPlug ") {
			return
		}
	}

	attachments := make([]*AttachmentInput, 0)

	for _, discordAttachment := range m.Attachments {
		attachment := AttachmentInput{
			Type:      "IMAGE",
			OriginID:  discordAttachment.ID,
			SourceURL: discordAttachment.URL,
		}

		fmt.Println("ATTACHMENT! " + attachment.OriginID + " hoho " + attachment.SourceURL)
		attachments = append(attachments, &attachment)
	}
	fmt.Printf("New message from %s in channel %s: %s\n", m.Author.Username, m.ChannelID, m.Content)
	fmt.Printf("mid: %s, uid: %s, avatar: %s\n", m.ID, m.Author.ID, m.Author.AvatarURL("medium"))
	fmt.Println("---")

	ds.client.SendMessage(
		m.Content,
		m.ID,
		m.ChannelID,
		m.Author.Username,
		m.Author.ID,
		m.Author.AvatarURL("medium"),
		attachments,
	)
}

func (ds *DiscordService) GetConfigurationSchema() []ConfigurationField {
	conf := make([]ConfigurationField, 0)
	ques1 := ConfigurationField{
		Type:         "STRING",
		Hint:         "Your Discord bot token",
		DefaultValue: "",
		Optional:     false,
		Mask:         true,
	}
	conf = append(conf, ques1)
	return conf
}

func (ds *DiscordService) GetConfiguration() (*DiscordServiceConfiguration, error) {
	file, err := ioutil.ReadFile("config.json")

	if err != nil {
		return nil, err
	}

	data := DiscordServiceConfiguration{}

	err = json.Unmarshal([]byte(file), &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func (ds *DiscordService) SaveConfiguration(conf []string) {
	confStruct := DiscordServiceConfiguration{
		BotToken: conf[0],
	}

	file, _ := json.MarshalIndent(&confStruct, "", " ")

	_ = ioutil.WriteFile("config.json", file, 0644)
}

func (ds *DiscordService) IsConfigured() bool {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return false
	}
	return true
}

func main() {
	client := DiscordService{}
	fmt.Println("Starting the service...")
	client.Startup(os.Args)
}

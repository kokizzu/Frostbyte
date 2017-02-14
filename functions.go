package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
)

// GetPageContents - Get page content based on URL.
// url: Valid url of image.
func GetPageContents(url string) ([]byte, error) {
	req, err := http.NewRequest("GET", url, nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body, err
}

// IsManager Check to see if a user has ManageServer Permissions.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
func IsManager(s *discordgo.Session, GuildID string, AuthorID string) bool {
	// Check the user permissions of the guild.
	perms, err := s.State.UserChannelPermissions(AuthorID, GuildID)
	if err == nil {
		if (perms & discordgo.PermissionManageServer) > 0 {
			return true
		}
	} else {
		return false
	}
	return false
}

// Save - Saves Database to config.json
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: The Message Object sent back from Discord.
func (bot *Object) Save() {
	for {
		<-time.After(5 * time.Minute)
		js, err := json.MarshalIndent(bot, "", "  ")
		if err == nil {
			ioutil.WriteFile("config.json", js, 0777)
		}
	}
}

/* Not functional yet
func (bot *Object) PruneMessages() {
	for {
		<-time.After(1 * time.Hour)
		for _, m := range bot.System.Messages {
			if m.Timestamp < time.Now()-(3600*24*7) {

			}
		}
	}
}
*/

// GetRoleID - Grabs the Discord Role ID
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// role: The Discord role
func (bot *Object) GetRoleID(s *discordgo.Session, role string) string {
	var id string
	r, err := s.State.Guild(bot.Guild)
	if err == nil {
		for _, v := range r.Roles {
			if v.Name == role {
				id = v.ID
			}
		}
	}
	return id
}

// MemberHasRole - Checks to see if the user has a role.
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// role: The Discord role
func (bot *Object) MemberHasRole(s *discordgo.Session, AuthorID string, role string) bool {
	therole := bot.GetRoleID(s, role)
	z, err := s.State.Member(bot.Guild, AuthorID)
	if err != nil {
		z, err = s.GuildMember(bot.Guild, AuthorID)
		if err != nil {
			fmt.Println("Error ->", err)
			return false
		}
	}
	for r := range z.Roles {
		if therole == z.Roles[r] {
			return true
		}
	}
	return false
}

// Register - Register new object.
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// m: Message Object sent back from Discord.
func (bot *Object) Register(s *discordgo.Session, m *discordgo.MessageCreate) {
	// check and make sure the server already exists in my collection.
	if bot.System != nil {
		return
	}
	c, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println(err)
		return
	}

	bot.Guild = c.GuildID
	chn := &Channels{
		Autorole: "",
		Greeting: "",
		ByeMsg:   "",
	}

	// Create a new Info pointer.
	info := &System{
		Prefix:   ".",
		Autorole: "",
		Greeting: "",
		ByeMsg:   "",
		Channels: chn,
	}
	// Add our Info object to the bot map.
	bot.System = info
}

// Task - Store new messages to object.
// bot: Main Object with all your settings.
// s: The Current Session between the bot and discord
// role: The Discord role
func (bot *Object) Task(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Don't track the bots messages.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if bot.System == nil {
		return
	}
	// Create a new pointer.
	msg := &Messages{
		ID:        m.ID,
		Author:    m.Author.ID,
		Channel:   m.ChannelID,
		Timestamp: time.Now().Unix(),
	}
	// Add this Message to our Info object.
	bot.System.Messages = append(bot.System.Messages, msg)
}

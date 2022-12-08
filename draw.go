/*
   Copyright (C) 2022 Tianyu Zhu eric@ericz.me

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
)

func draw(s *discordgo.Session, m *discordgo.MessageCreate, msg string) {
	// Token
	token := k.String("ai.draw.token")
	var bearer = "Bearer " + token

	// send API req
	type req struct {
		Prompt string `json:"prompt"`
		n      int    `json:"n"`
		size   string `json:"size"`
		user   string
	}
	request := req{
		Prompt: msg,
		n:      1,
		size:   "1024x1024",
		user:   m.Author.ID,
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "Error", m.MessageReference)
		return
	}
	httpReq, err := http.NewRequest("POST", "https://api.openai.com/v1/images/generations", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "Error", m.MessageReference)
		return
	}
	httpReq.Header.Set("Authorization", bearer)
	httpReq.Header.Set("Content-Type", "application/json")
	httpClient := &http.Client{}
	resp, err := httpClient.Do(httpReq)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
		s.ChannelMessageSendReply(m.ChannelID, "Error http", m.Reference())
		return
	}

	// Decode response
	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	user := m.Author.Username + "#" + m.Author.Discriminator

	// Handle respons
	if result == nil {
		log.Error().Str("user", user).Str("prompt", msg).Msg("Chat: result is nil")
		s.ChannelMessageSendReply(m.ChannelID, "Draw: result is nil", m.Reference())
		return
	}
	if result["data"] == nil {
		resultStr := fmt.Sprintf("%#v", result)
		log.Error().Str("user", user).Str("prompt", msg).Str("resp", resultStr).Msg("Draw: data is nil")
		s.ChannelMessageSendReply(m.ChannelID, "Draw: data is nil (likely OpenAI rejecting request due to inappropriate prompt)", m.Reference())
		return
	}

	// Send
	aiRespStr := result["data"].([]interface{})[0].(map[string]interface{})["url"].(string)
	log.Info().Str("user", user).Str("prompt", msg).Str("resp", aiRespStr).Msg("Draw: success")
	s.ChannelMessageSendReply(m.ChannelID, aiRespStr, m.Reference())
}
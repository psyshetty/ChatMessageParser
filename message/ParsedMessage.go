package message

import "github.com/psyshetty/ChatMessageParser/message/link"

type ParsedMessage struct {
  Mentions []string `json:"mentions"`
  Emoticons []string `json:"emoticons"`
  Links []link.Link `json:"links"`
}


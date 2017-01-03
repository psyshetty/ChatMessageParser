package main

import (
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "encoding/json"
  "fmt"
  "github.com/psyshetty/ChatMessageParser/message"
  "github.com/psyshetty/ChatMessageParser/message/link"
  "regexp"
  "strings"
  "html"
  "net/url"
)

func setMentions(message string, pm *message.ParsedMessage) {
  var mentionsExp = regexp.MustCompile(`@\w+`)
  matches  := mentionsExp.FindAllString(message, -1)
  for _, match := range matches {
    (*pm).Mentions = append((*pm).Mentions, strings.TrimSpace(match[1:]))
  }
  fmt.Println("Mentions: ", (*pm).Mentions)
}

func setEmoticons(message string, pm *message.ParsedMessage) {
  var emoticonsExp = regexp.MustCompile(`\(\w+\)`)
  matches  := emoticonsExp.FindAllString(message, -1)
  for _, match := range matches {
    (*pm).Emoticons = append((*pm).Emoticons, match[1:len(match)-1])
  }
  fmt.Println("Emoticons: ", (*pm).Emoticons)
}

func setLinks(message string, pm *message.ParsedMessage) {
  var linksExp = regexp.MustCompile(`http(s?)://\S+`)
  matches  := linksExp.FindAllString(message, -1)
  for _, match := range matches {
    // test if URL is valid
    _, err := url.ParseRequestURI(match)
    if err == nil {
      titles := html.EscapeString(link.GetTitle(match))
      fmt.Println("Title from func : ", titles)
      newlink := link.Link{Url:match, Titles:titles}
      (*pm).Links = append((*pm).Links, newlink)
    }
  }
  fmt.Println("Links: ", (*pm).Links)
}

//post chat message API endpoint
func parseChatMessage(w http.ResponseWriter, req *http.Request) {
  var chatMessage message.ChatMessage
  _ = json.NewDecoder(req.Body).Decode(&chatMessage)
  fmt.Println("Chat Message : ", chatMessage.Message)
  mentionsSlice := []string{}
  emoticonsSlice := []string{}
  linksSlice := []link.Link{}
  pm := message.ParsedMessage{Mentions:mentionsSlice,Emoticons:emoticonsSlice,Links:linksSlice}
  setMentions(chatMessage.Message, &pm)
  setEmoticons(chatMessage.Message, &pm)
  setLinks(chatMessage.Message, &pm)
  enc := json.NewEncoder(w)
  enc.SetEscapeHTML(false)
  enc.Encode(pm)
}

func main() {
  router := mux.NewRouter()
  router.HandleFunc("/chat/message/parse", parseChatMessage).Methods("POST")
  log.Fatal(http.ListenAndServe(":12345", router))
}

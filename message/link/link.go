package link

import (
  "golang.org/x/net/html"
  "io"
  "net/http"
)

type Link struct {
  Url string  `json:"url"`
  Titles string `json:"title"`
}

func GetTitle(url string) string {
  resp, err := http.Get(url)
  if err != nil {
    panic(err)
  }
  defer resp.Body.Close()
  titles, ok := getHtmlTitle(resp.Body);
  if ok {
    println(titles)
  } else {
    println("Fail to get HTML title")
  }
  return titles
}

func getHtmlTitle(r io.Reader) (string, bool) {
  doc, err := html.Parse(r)
  if err != nil {
    panic("Fail to parse html")
  }

  return traverse(doc)
}

func isTitleElement(n *html.Node) bool {
  return n.Type == html.ElementNode && n.Data == "title"
}


func traverse(n *html.Node) (string, bool) {
  if isTitleElement(n) {
    return n.FirstChild.Data, true
  }

  for c := n.FirstChild; c != nil; c = c.NextSibling {
    result, ok := traverse(c)
    if ok {
      return result, ok
    }
  }

  return "", false
}



package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"qtbooru/pkg/api/post"
	"strings"

	"github.com/joho/godotenv"
)

type Booru int
const (
	E621 Booru = iota
	E926

	E621_URL string = "https://e621.net/posts.json"
	E926_URL string = "https://e926.net/posts.json"
)
func URL(u string) string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	if url := os.Getenv("OVERRIDE_URL"); url == "" {
		return u
	} else {
		fmt.Println("redirecting to " + url)
		return url
	}
}
var BooruURL = map[Booru]string{
	E621: URL(E621_URL),
	E926: URL(E926_URL),
}
func (b Booru) String() string {
	return BooruURL[b]
}

type RequestBuilder struct {
	Site Booru
	Params *[]string
	Tags *[]string
	User string
	Key string
	Agent string
}

func (r *RequestBuilder) Process(client *http.Client) (*[]*post.Post, error) {
	tags := strings.Join(*r.Tags, " ")
	if tags != "" { tags = "tags=" + tags }
	params := strings.Join(append(*r.Params, tags), "&")
	if params != "" { params = "?" + params }

	req, err := http.NewRequest("GET", r.Site.String() + params, nil)
	if err != nil { return nil, nil }

	req.Header.Set("User-Agent", r.Agent)
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(r.User+":"+r.Key)))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	p := Posts{}
	err = json.Unmarshal(body, &p)
	if err != nil {
		log.Fatal(err)
	}
	if !p.Success && len(p.Message) > 0 {
		return nil, fmt.Errorf("%s", p.Message)
	}
	return p.Posts, nil
}

type Posts struct {
	Posts   *[]*post.Post `json:"posts"`
	Success bool          `json:"success"`
	Message string        `json:"message"`
}


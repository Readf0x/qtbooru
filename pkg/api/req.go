package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"qtbooru/pkg/api/post"
	"strings"

	"github.com/joho/godotenv"
)

const Agent = "QtBooru/indev_v0 (created by readf0x)"

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
}

func (r *RequestBuilder) Build() (req *http.Request, err error) {
	params := strings.Join(*r.Params, "&")
	params += strings.Join(*r.Tags, " ")
	if params != "" { params = "?" + params }

	req, err = http.NewRequest("GET", r.Site.String() + params, nil)
	if err != nil { return }

	req.Header.Set("User-Agent", Agent)
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(r.User+":"+r.Key)))
	return
}

func NewRequest(site Booru, params *[]string, tags *[]string, user string, key string) (req *http.Request, err error) {
	req, err = (&RequestBuilder{
		Site: site,
		Params: params,
		Tags: tags,
		User: user,
		Key: key,
	}).Build()
	return
}

func Process(req *http.Request) *[]*post.Post {
	client := &http.Client{}

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
	return p.Posts
}

type Posts struct {
	Posts *[]*post.Post `json:"posts"`
}


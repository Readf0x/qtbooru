package api
//go:generate go tool oapi-codegen -config cfg.yaml ../../api.yaml
import (
	"encoding/base64"
	"net/http"
	"strings"
)

type Booru int
const (
	E621 Booru = iota
	E926
)
var BooruUrl = map[Booru]string{
	E621: "https://e621.net/posts.json",
	E926: "https://e926.net/posts.json",
}
func (b Booru) String() string {
	return BooruUrl[b]
}

type RequestBuilder struct {
	Site Booru
	Arguments *[]string
}

func (r *RequestBuilder) Build(user string, key string, agent string) (req *http.Request, err error) {
	args := strings.Join(*r.Arguments, "&")
	if args != "" { args = "?" + args }

	req, err = http.NewRequest("GET", r.Site.String() + args, nil)
	if err != nil { return }

	req.Header.Set("User-Agent", agent)
	req.Header.Set("Authorization", "Basic " + base64.StdEncoding.EncodeToString([]byte(user+":"+key)))
	return
}

func NewRequest(site Booru, args *[]string, user string, key string, agent string) (req *http.Request, err error) {
	req, err = (&RequestBuilder{
		Site: site,
		Arguments: args,
	}).Build(user, key, agent)
	return
}


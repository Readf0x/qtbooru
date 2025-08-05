package post

import (
	"encoding/json"
	"time"
)

type PostID int
type UserID int

type Post struct {
	ID            PostID     `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	File          FileFull   `json:"file"`
	Preview       File       `json:"preview"`
	Sample        FileSample `json:"sample"`
	Score         Score      `json:"score"`
	Tags          TagList    `json:"tags"`
	LockedTags    *[]string  `json:"locked_tags"`
	ChangeSeq     int        `json:"change_seq"`
	Flags         Flags      `json:"flags"`
	Rating        Rating     `json:"rating"`
	FavCount      int        `json:"fav_count"`
	Sources       []string   `json:"sources"`
	Pools         []PostID   `json:"pools"`
	Relationships Relations  `json:"relationships"`
	ApproverID    *UserID    `json:"approver_id"`
	UploaderID    UserID     `json:"uploader_id"`
	Description   string     `json:"description"`
	CommentCount  int        `json:"comment_count"`
	IsFavorited   bool       `json:"is_favorited"`
	HasNotes      bool       `json:"has_notes"`
	Duration      *float32   `json:"duration"`
}

type Extension uint8
const (
	JPEG Extension = iota
	PNG
	GIF
	WEBM
)
var ExtensionString = map[Extension]string{
	JPEG: "jpg",
	PNG:  "png",
	GIF:  "gif",
	WEBM: "webm",
}
var StringExtension = map[string]Extension{
	"jpg":  JPEG,
	"png":  PNG,
	"gif":  GIF,
	"webm": WEBM,
}
func (e Extension) String() string {
	return ExtensionString[e]
}
func (e *Extension) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*e = StringExtension[s]
	return nil
}

type File struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url"`
}

type FileFull struct {
	File
	Type Extension `json:"ext"`
	Size int64     `json:"size"`
	Md5  Md5Hash   `json:"md5"`
}

type Md5Hash [16]byte
func (m *Md5Hash) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	copy(m[:], s)
	return nil
}

type FileSample struct {
	File
	Alternates *Alternates `json:"alternates"`
}

type Alternates struct {
	Has      bool                 `json:"has"`
	Original *PostSampleAlternate `json:"original"`
	Variants *AlternateVariants   `json:"variants"`
	Samples  *AlternateSamples    `json:"samples"`
}

type AlternateVariants struct {
	Webm *PostSampleAlternate `json:"webm"`
	Mp4  *PostSampleAlternate `json:"mp4"`
}

type AlternateSamples struct {
	P480 *PostSampleAlternate `json:"480p"`
	P720 *PostSampleAlternate `json:"720p"`
}

type PostSampleAlternate struct {
	File
	Fps   float64 `json:"fps"`
	Codec string  `json:"codec"`
	Size  int64   `json:"size"`
}

type Score struct {
	Up   int `json:"up"`
	Down int `json:"down"`
	Sum  int `json:"total"`
}

type TagList struct {
	General     *[]string `json:"general"`
	Artist      *[]string `json:"artist"`
	Contributor *[]string `json:"contributor"`
	Copyright   *[]string `json:"copyright"`
	Character   *[]string `json:"character"`
	Species     *[]string `json:"species"`
	Invalid     *[]string `json:"invalid"`
	Meta        *[]string `json:"meta"`
	Lore        *[]string `json:"lore"`
}

type Flags struct {
	Pending      bool `json:"pending"`
	Flagged      bool `json:"flagged"`
	NoteLocked   bool `json:"note_locked"`
	StatusLocked bool `json:"status_locked"`
	RatingLocked bool `json:"rating_locked"`
	Deleted      bool `json:"deleted"`
}

type Rating uint8
const (
	Safe Rating = iota
	Questionable
	Explicit
)
var Ratings = map[byte]Rating{
	's': Safe,
	'q': Questionable,
	'e': Explicit,
}
func (r *Rating) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*r = Ratings[s[0]]
	return nil
}

type Relations struct {
	ParentID          PostID    `json:"parent_id"`
	HasChildren       bool      `json:"has_children"`
	HasActiveChildren bool      `json:"has_active_children"`
	Children          *[]PostID `json:"children"`
}


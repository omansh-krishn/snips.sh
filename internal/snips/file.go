package snips

import (
	"net/url"
	"strconv"
	"time"

	"github.com/robherley/snips.sh/internal/config"
	"github.com/robherley/snips.sh/internal/signer"
)

const (
	FileTypeBinary   = "binary"
	FileTypeMarkdown = "markdown"
)

type File struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
	Size      uint64
	Content   []byte
	Private   bool
	Type      string
	UserID    string
}

func (f *File) IsBinary() bool {
	return f.Type == FileTypeBinary
}

func (f *File) IsMarkdown() bool {
	return f.Type == FileTypeMarkdown
}

func (f *File) GetSignedURL(cfg *config.Config, ttl time.Duration) (url.URL, time.Time) {
	expires := time.Now().Add(ttl).UTC()

	println("want:", cfg.HTTPAddressForFile(f.ID).Path)

	pathToSign := url.URL{
		Path: cfg.HTTPAddressForFile(f.ID).Path,
		RawQuery: url.Values{
			"exp": []string{strconv.FormatInt(expires.Unix(), 10)},
		}.Encode(),
	}

	signedFileURL := signer.New(cfg.HMACKey).SignURL(pathToSign)
	signedFileURL.Scheme = cfg.HTTP.External.Scheme
	signedFileURL.Host = cfg.HTTP.External.Host

	return signedFileURL, expires
}

func (f *File) Visibility() string {
	if f.Private {
		return "private"
	}

	return "public"
}

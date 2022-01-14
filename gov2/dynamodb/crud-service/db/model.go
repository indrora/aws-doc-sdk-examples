package db

import (
	"net/url"
	"time"

	"github.com/rs/xid"
)

type (
	Link struct {
		Id        string
		Uri       string
		DeleteKey string
		CreatedOn time.Time
		NumHits   uint64
	}
)

func CreateLink(longUrl string) (*Link, error) {

	// Validate the url:

	mUrl, err := url.Parse(longUrl)
	createTime := time.Now()
	linkId := xid.NewWithTime(createTime)
	deleteKey := xid.New().String()

	if err != nil {
		return nil, err
	}

	mLink := &Link{
		Id:        linkId.String(),
		DeleteKey: deleteKey,
		Uri:       (*mUrl).String(),
		CreatedOn: createTime,
		NumHits:   0,
	}

	return mLink, nil

}

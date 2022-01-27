package db

import (
	"errors"
	"net/url"
	"time"

	"github.com/rs/xid"
)

type (
	Link struct {
		Id        string    `dynamodbav:"Id" json:"id"`
		Uri       string    `dynamodbav:"Uri" json:"uri"`
		Email     string    `dynamodbav:"Email" json:"-"`
		DeleteKey string    `dynamodbav:"DeleteKey" json:"-"`
		CreatedOn time.Time `dynamodbav:"CreatedOn" json:"createdOn"`
		NumHits   uint64    `dynamodbav:"Hits" json:"hits"`
	}
)

func CreateLink(longUrl string, email string) (Link, error) {

	// Validate the url:
	mUrl, err := url.Parse(longUrl)

	if err != nil {
		return Link{}, err
	}

	if email == "" {
		return Link{}, errors.New("invalid email")
	}

	createTime := time.Now()
	linkId := xid.NewWithTime(createTime)
	deleteKey := xid.New().String()

	mLink := Link{
		Id:        linkId.String(),
		DeleteKey: deleteKey,
		Email:     email,
		Uri:       (*mUrl).String(),
		CreatedOn: createTime,
		NumHits:   0,
	}

	return mLink, nil

}

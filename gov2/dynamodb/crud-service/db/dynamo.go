package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/rs/zerolog/log"
)

type DynamoConnection struct {
	client    *dynamodb.Client
	TableName string
}

func GetDDBConnection(tablename string, conf aws.Config) *DynamoConnection {

	client := dynamodb.NewFromConfig(conf)

	conn := &DynamoConnection{
		TableName: tablename,
		client:    client,
	}

	return conn
}

// Create a link
func (db *DynamoConnection) Add(link Link) {

	marshalledLink, err := attributevalue.MarshalMap(link)

	if err != nil {
		log.Error().AnErr("err", err).Msg("Failed to marshal link to ddb")
		return
	}

	_, err = db.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item: marshalledLink,
		TableName: &db.TableName,
	})
	if err != nil {
		log.Error().AnErr("err", err).Msg("Failed to insert link into dynamodb")
	}
}

// List recent links
func (db *DynamoConnection) ListRecent() []Link {
	return nil
}

// Get a link by its ID
func (db *DynamoConnection) Get(id string) *Link {
	response, err := db.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		return nil
	}

	mLink := Link{}

	err = attributevalue.Unmarshal(response.Item, mLink)

	if err != nil {
		return nil
	}

	return &mLink
}

// Destroy a link by its ID
func (db *DynamoConnection) Delete(id string) bool {
	// check if the link exists
	link := db.Get(id)
	if link == nil {
		return false // didn't exist, fails to delete.
	}

	_, err := db.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"id": types.AttributeValueMemberS{Value: id}
		},
	})
	
	if err != nil {
		log.Error().AnErr("err", err).Str("id", id).Msg("Failed to destroy link in ddb")
		return false
	}
	return true
}

// increment the view count on a link by its ID
func (db *DynamoConnection) Increment(id string) {
	// no return
}

package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
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
func (db DynamoConnection) Add(link Link) error {

	marshalledLink, err := attributevalue.MarshalMap(link)

	if err != nil {
		log.Error().AnErr("err", err).Msg("Failed to marshal link to ddb")
		return errors.New("internal error")
	}

	_, err = db.client.PutItem(context.Background(), &dynamodb.PutItemInput{
		Item:      marshalledLink,
		TableName: &db.TableName,
	})

	if err != nil {
		log.Error().AnErr("err", err).Msg("Failed to insert link into dynamodb")
		return errors.New("database error")
	}

	return nil
}

// List recent links
func (db DynamoConnection) ListByEmail(email string) []Link {

	filterExpression := expression.Name("Email").Equal(expression.Value(email))

	expr, err := expression.NewBuilder().WithFilter(filterExpression).Build()

	if err != nil {
		log.Error().AnErr("err", err).Msg("Failed to create filter expression")
		return nil
	}

	result, err := db.client.Scan(context.Background(), &dynamodb.ScanInput{
		TableName:                 &db.TableName,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})

	if err != nil {
		log.Error().AnErr("err", err).Msg("Failed to scan table")
		return nil
	}

	links := []Link{}
	err = attributevalue.UnmarshalListOfMaps(result.Items, &links)
	if err != nil {
		log.Err(err).Msg("Failed to unmarshall scanned email results to link")
		return nil
	}

	return links
}

// Get a link by its ID
func (db DynamoConnection) Get(id string) *Link {
	response, err := db.client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		log.Err(err).Str("id", id).Msg("Failed to retrieve a link")
		return nil
	}

	// Check to make sure there was something at that key:

	if response.Item == nil {
		// This means there was no item at that key.
		return nil
	}

	mLink := new(Link)

	err = attributevalue.UnmarshalMap(response.Item, mLink)

	if err != nil {
		log.Err(err).Msg("Failed to unmarshal DynamoDB response to Item")
		return nil
	}

	return mLink
}

// Destroy a link by its ID
func (db DynamoConnection) Delete(id string) bool {
	// check if the link exists
	link := db.Get(id)
	if link == nil {
		return false // didn't exist, fails to delete.
	}

	_, err := db.client.DeleteItem(context.Background(), &dynamodb.DeleteItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{Value: id},
		},
	})

	if err != nil {
		log.Error().AnErr("err", err).Str("id", id).Msg("Failed to destroy link in ddb")
		return false
	}
	return true
}

// increment the view count on a link by its ID
func (db DynamoConnection) Increment(id string) {
	_, err := db.client.UpdateItem(context.Background(), &dynamodb.UpdateItemInput{
		TableName: &db.TableName,
		Key: map[string]types.AttributeValue{
			"Id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
		UpdateExpression: aws.String("ADD Hits :v"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberN{
				Value: "1",
			},
		},
	})
	if err != nil {
		log.Error().AnErr("err", err).Str("id", id).Msg("Failed to update hit count")
	}
}

package database

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/holmes89/magpie/lib"
	v1 "github.com/holmes89/magpie/lib/handlers/rest/v1"
)

var (
	_        v1.SiteService = (*Conn)(nil)
	notfound *types.ResourceNotFoundException
)

func (conn *Conn) Get(ctx context.Context, id string) (lib.Site, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"SK": &types.AttributeValueMemberS{Value: id},
			"ID": &types.AttributeValueMemberS{Value: "site"},
		},
	}

	rs := lib.Site{}
	resp, err := conn.db.GetItem(ctx, params)
	if err != nil {
		if errors.As(err, &notfound) {
			log.Println("no resources found")
			return lib.Site{}, nil
		}
		log.Println("unable to find Site", err)
		return lib.Site{}, errors.New("unable to fetch Site")
	}

	if err := attributevalue.UnmarshalMap(resp.Item, &rs); err != nil {
		log.Println("unable to unmarshal Site", err)
		return lib.Site{}, errors.New("failed to scan Site")
	}

	return rs, nil
}

func (conn *Conn) GetAll(ctx context.Context) ([]lib.Site, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		Limit:                  aws.Int32(10),
		KeyConditionExpression: aws.String("ID = :key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":key": &types.AttributeValueMemberS{Value: "site"},
		},
		ScanIndexForward: aws.Bool(false),
	}

	sites := make([]lib.Site, 0)
	resp, err := conn.db.Query(ctx, params)
	if err != nil {
		if errors.As(err, &notfound) {
			log.Println("no resources found")
			return sites, nil
		}
		log.Println("unable to fetch Sites", err)
		return sites, errors.New("unable to fetch all Sites")
	}
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &sites); err != nil {
		log.Println("unable to unmarshal Sites", err)
		return sites, errors.New("unable to fetch all Sites")
	}
	r := make([]lib.Site, 0)
	for _, s := range sites {
		s.Meta = nil
		r = append(r, s)
	}
	return r, nil
}

func (conn *Conn) Create(ctx context.Context, r lib.Site) error {
	site := r
	site.Type = "site"
	rs, err := attributevalue.MarshalMap(site)
	if err != nil {
		log.Println("unable to marshal Site message", err)
		return errors.New("failed to insert Site")
	}
	params := &dynamodb.PutItemInput{
		Item:      rs,
		TableName: aws.String(tableName),
	}
	if _, err := conn.db.PutItem(ctx, params); err != nil {
		log.Println("unable to put Site message", err)
		return errors.New("failed to insert Site")
	}
	return nil
}

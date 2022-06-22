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

type Site struct {
	ID string
	SK string
	lib.Site
}

func (conn *Conn) Get(ctx context.Context, id string) (lib.Site, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"SK": &types.AttributeValueMemberS{Value: id},
			"ID": &types.AttributeValueMemberS{Value: "site"},
		},
	}

	rs := Site{}
	resp, err := conn.db.GetItem(ctx, params)
	if err != nil {
		if errors.As(err, &notfound) {
			log.Println("no devices found")
			return lib.Site{}, nil
		}
		log.Println("unable to find Site", err)
		return lib.Site{}, errors.New("unable to fetch Site")
	}

	if err := attributevalue.UnmarshalMap(resp.Item, &rs); err != nil {
		log.Println("unable to unmarshal Site", err)
		return lib.Site{}, errors.New("failed to scan Site")
	}

	return rs.Site, nil
}

func (conn *Conn) GetAll(ctx context.Context, r string) ([]lib.Site, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		Limit:                  aws.Int32(10),
		KeyConditionExpression: aws.String("ID = :key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":key": &types.AttributeValueMemberS{Value: "site"},
		},
		ScanIndexForward: aws.Bool(false),
	}

	resp, err := conn.db.Query(ctx, params)
	if err != nil {
		if errors.As(err, &notfound) {
			log.Println("no log found")
			return nil, nil
		}
		log.Println("unable to fetch Sites", err)
		return nil, errors.New("unable to fetch all Sites")
	}
	var SiteEntities []Site
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &SiteEntities); err != nil {
		log.Println("unable to unmarshal Sites", err)
		return nil, errors.New("unable to fetch all Sites")
	}

	Sites := make([]lib.Site, 0)
	for _, Site := range SiteEntities {
		Sites = append(Sites, Site.Site)
	}

	return Sites, nil
}

func convertSite(r lib.Site) Site {
	return Site{
		ID:   "stie",
		SK:   r.ResourceID,
		Site: r,
	}
}

func (conn *Conn) Create(ctx context.Context, r lib.Site) error {
	rs, err := attributevalue.MarshalMap(convertSite(r))
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

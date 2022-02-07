package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/holmes89/magpie/lib"
	v1 "github.com/holmes89/magpie/lib/handlers/rest/v1"
)

var (
	_        v1.ResourceService = (*Conn)(nil)
	notfound *types.ResourceNotFoundException
)

type resource struct {
	ID string
	SK string
	lib.Resource
}

func (conn *Conn) Get(ctx context.Context, r string, id string) (lib.Resource, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]types.AttributeValue{
			"SK": &types.AttributeValueMemberS{Value: id},
			"ID": &types.AttributeValueMemberS{Value: resourceKey(r)},
		},
	}

	rs := resource{}
	resp, err := conn.db.GetItem(ctx, params)
	if err != nil {
		if errors.As(err, &notfound) {
			log.Println("no devices found")
			return lib.Resource{}, nil
		}
		log.Println("unable to find resource", err)
		return lib.Resource{}, errors.New("unable to fetch resource")
	}

	if err := attributevalue.UnmarshalMap(resp.Item, &rs); err != nil {
		log.Println("unable to unmarshal resource", err)
		return lib.Resource{}, errors.New("failed to scan resource")
	}

	return rs.Resource, nil
}

func resourceKey(id string) string {
	return fmt.Sprintf("r#%s", strings.ToLower(id))
}

func (conn *Conn) GetAll(ctx context.Context, r string) ([]lib.Resource, error) {
	params := &dynamodb.QueryInput{
		TableName:              aws.String(tableName),
		Limit:                  aws.Int32(10),
		KeyConditionExpression: aws.String("ID = :key"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":key": &types.AttributeValueMemberS{Value: resourceKey(r)},
		},
		ScanIndexForward: aws.Bool(false),
	}

	resp, err := conn.db.Query(ctx, params)
	if err != nil {
		if errors.As(err, &notfound) {
			log.Println("no log found")
			return nil, nil
		}
		log.Println("unable to fetch resources", err)
		return nil, errors.New("unable to fetch all resources")
	}
	var resourceEntities []resource
	if err := attributevalue.UnmarshalListOfMaps(resp.Items, &resourceEntities); err != nil {
		log.Println("unable to unmarshal resources", err)
		return nil, errors.New("unable to fetch all resources")
	}

	resources := make([]lib.Resource, 0)
	for _, resource := range resourceEntities {
		resources = append(resources, resource.Resource)
	}

	return resources, nil
}

func convertResource(r lib.Resource) resource {
	return resource{
		ID:       resourceKey(r.Type),
		SK:       r.ResourceID,
		Resource: r,
	}
}

func (conn *Conn) Create(ctx context.Context, r lib.Resource) error {
	rs, err := attributevalue.MarshalMap(convertResource(r))
	if err != nil {
		log.Println("unable to marshal resource message", err)
		return errors.New("failed to insert resource")
	}
	params := &dynamodb.PutItemInput{
		Item:      rs,
		TableName: aws.String(tableName),
	}
	if _, err := conn.db.PutItem(ctx, params); err != nil {
		log.Println("unable to put resource message", err)
		return errors.New("failed to insert resource")
	}
	return nil
}

package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// InitSdk sets up AWS SDK
func (d *DynamoDB) InitSdk() error {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		return errors.New("AWS SDK configuration error, " + err.Error())
	}
	d.client = dynamodb.NewFromConfig(cfg)
	return nil
}

// DynamoDB is a repository to be used by this project
type DynamoDB struct {
	TableName      string
	KeyAttribute   string
	ValueAttribute string
	client         *dynamodb.Client
}

func (d *DynamoDB) Lookup(key string) (string, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(d.TableName),
		Key: map[string]*types.AttributeValue{
			d.KeyAttribute: &types.AttributeValue{S: aws.String(key)},
		},
	}
	res, err := d.client.GetItem(context.Background(), params)
	if err != nil {
		return "", err
	}

	value, ok := res.Item[d.ValueAttribute]
	if !ok || value.S == nil {
		return "", errors.New("item or attribute not found")
	}
	return *value.S, nil
}

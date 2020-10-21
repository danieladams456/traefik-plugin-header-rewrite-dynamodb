package traefik_plugin_header_rewrite_dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// temp until we wire up the traefik configuration
const (
	tableName      = "traefik_header_lookups"
	keyAttribute   = "key"
	valueAttribute = "value"
)

func get(key string) (string, error) {
	cfg, err := config.LoadDefaultConfig()
	if err != nil {
		return "", errors.New("AWS SDK configuration error, " + err.Error())
	}
	client := dynamodb.NewFromConfig(cfg)

	params := &dynamodb.GetItemInput{
		TableName: aws.String(tableName),
		Key: map[string]*types.AttributeValue{
			keyAttribute: &types.AttributeValue{S: aws.String(key)},
		},
	}
	res, err := client.GetItem(context.Background(), params)
	if err != nil {
		return "", err
	}

	value, ok := res.Item[valueAttribute]
	if !ok {
		return "", errors.New("value attribute not found in item")
	}
	if value.S == nil {
		return "", errors.New("value attribute is not a string")
	}
	return *value.S, nil
}

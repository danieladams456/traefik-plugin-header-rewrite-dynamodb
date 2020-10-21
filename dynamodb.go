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

// This function is now implemented in headerrewrite.go as part of the HeaderRewrite type
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
	if !ok || value.S == nil {
		return "", errors.New("item or attribute not found")
	}
	return *value.S, nil
}

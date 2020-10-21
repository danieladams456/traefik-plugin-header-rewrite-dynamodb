package traefik_plugin_header_rewrite_dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

func (a *HeaderRewrite) lookup(key string) (string, error) {
	params := &dynamodb.GetItemInput{
		TableName: aws.String(a.tableName),
		Key: map[string]*types.AttributeValue{
			a.keyAttribute: &types.AttributeValue{S: aws.String(key)},
		},
	}
	res, err := a.dynamodbClient.GetItem(context.Background(), params)
	if err != nil {
		return "", err
	}

	value, ok := res.Item[a.valueAttribute]
	if !ok || value.S == nil {
		return "", errors.New("item or attribute not found")
	}
	return *value.S, nil
}

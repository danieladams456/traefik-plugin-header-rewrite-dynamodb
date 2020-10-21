// Package traefik_plugin_header_rewrite_dynamodb rewrites headers based off of K/V stored in DynamoDB
// TODO change package name and probably repo name to align to go convention
package traefik_plugin_header_rewrite_dynamodb

import (
	"context"
	"errors"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsSdkConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Config the plugin configuration.
type Config struct {
	SourceHeader   string
	TargetHeader   string
	TableName      string
	KeyAttribute   string
	ValueAttribute string
}

// CreateConfig creates the default plugin configuration.
// KeyAttribute and ValueAttribute can be defaulted
func CreateConfig() *Config {
	return &Config{
		KeyAttribute:   "key",
		ValueAttribute: "value",
	}
}

// HeaderRewrite a header rewrite plugin.
type HeaderRewrite struct {
	next http.Handler
	name string
	// copy attributes from Config to non-exported
	sourceHeader   string
	targetHeader   string
	tableName      string
	keyAttribute   string
	valueAttribute string
	dynamodbClient *dynamodb.Client
}

// New creates a HeaderRewrite plugin
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	// validate config
	if config.SourceHeader == "" {
		return nil, errors.New("SourceHeader cannot be empty")
	}
	if config.TargetHeader == "" {
		return nil, errors.New("TargetHeader cannot be empty")
	}
	if config.TableName == "" {
		return nil, errors.New("TableName cannot be empty")
	}

	// create dynamodb client, later can refactor out into a method
	// passing plugin config if we need to influence things more than LoadDefaultConfig()
	cfg, err := awsSdkConfig.LoadDefaultConfig()
	if err != nil {
		return nil, errors.New("AWS SDK configuration error, " + err.Error())
	}
	dynamodbClient := dynamodb.NewFromConfig(cfg)

	return &HeaderRewrite{
		next:           next,
		name:           name,
		sourceHeader:   config.SourceHeader,
		targetHeader:   config.TargetHeader,
		tableName:      config.TableName,
		keyAttribute:   config.KeyAttribute,
		valueAttribute: config.ValueAttribute,
		// TBD if we separate this whole thing out into a dynamodbRepository type
		// that encapsulates the construction and lookup
		dynamodbClient: dynamodbClient,
	}, nil
}

func (a *HeaderRewrite) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	// will only get first value of a header, intended behavior
	if key := req.Header.Get(a.sourceHeader); key != "" {

		val, err := a.lookup(key)
		if err != nil {
			req.Header.Set(a.targetHeader, val)
		}
	}
	a.next.ServeHTTP(rw, req)
}

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

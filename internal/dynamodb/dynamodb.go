package dynamodb

import (
	"errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// InitSdk sets up AWS SDK
func (r *Repository) InitSdk() {
	session := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))
	r.client = dynamodb.New(session)
}

// Repository is a repository to be used by this project
type Repository struct {
	TableName      string
	KeyAttribute   string
	ValueAttribute string
	client         *dynamodb.DynamoDB
}

// Lookup retrieves the value associated with the specified key
func (r *Repository) Lookup(key string) (string, error) {
	res, err := r.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]*dynamodb.AttributeValue{
			r.KeyAttribute: {
				S: aws.String(key),
			},
		},
	})
	if err != nil {
		return "", err
	}
	value, ok := res.Item[r.ValueAttribute]
	if !ok || value.S == nil {
		return "", errors.New("item or attribute not found")
	}
	return *value.S, nil
}

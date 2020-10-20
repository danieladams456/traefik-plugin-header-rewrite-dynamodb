// Package traefik_plugin_header_rewrite_dynamodb rewrites headers based off of K/V stored in DynamoDB
// TODO change package name and probably repo name to align to go convention
package traefik_plugin_header_rewrite_dynamodb

import (
	"context"
	"errors"
	"net/http"
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
}

// New created a new Demo plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	if config.SourceHeader == "" {
		return nil, errors.New("SourceHeader cannot be empty")
	}
	if config.TargetHeader == "" {
		return nil, errors.New("TargetHeader cannot be empty")
	}
	if config.TableName == "" {
		return nil, errors.New("TableName cannot be empty")
	}

	return &HeaderRewrite{
		next:           next,
		name:           name,
		sourceHeader:   config.SourceHeader,
		targetHeader:   config.TargetHeader,
		tableName:      config.TableName,
		keyAttribute:   config.KeyAttribute,
		valueAttribute: config.ValueAttribute,
	}, nil
}

// func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
// 	for key, value := range a.headers {
// 		tmpl, err := a.template.Parse(value)
// 		if err != nil {
// 			http.Error(rw, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		writer := &bytes.Buffer{}

// 		err = tmpl.Execute(writer, req)
// 		if err != nil {
// 			http.Error(rw, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		req.Header.Set(key, writer.String())
// 	}

// 	a.next.ServeHTTP(rw, req)
// }

// TODO
func (a *HeaderRewrite) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
}

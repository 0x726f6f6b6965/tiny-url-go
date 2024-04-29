package storage

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	appConfig "github.com/0x726f6f6b6965/tiny-url-go/internal/config"
	"github.com/0x726f6f6b6965/tiny-url-go/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type dynamo struct {
	DynamoClient *dynamodb.Client
}

const (
	pk                 = "pk"
	sk                 = "sk"
	pkNotExists string = "attribute_not_exists(pk)"
	pkExists    string = "attribute_exists(pk)"
)

var (
	ErrInvalidKey = errors.New("invalid key")
	ErrNotSupport = errors.New("not support")
	ErrDynamoDB   = errors.New("dynamodb error")
	ErrNotFound   = errors.New("not found")
)

// Delete implements utils.Storage.
// key format: <table>;<partition key>;<sort key>
func (d *dynamo) Delete(ctx context.Context, key string) error {
	keys := strings.Split(key, ";")
	if len(keys) != 3 {
		return ErrInvalidKey
	}
	tableName := keys[0]
	partitionKey := keys[1]
	sortKey := keys[2]

	result := make(map[string]types.AttributeValue)
	result[pk] = &types.AttributeValueMemberS{Value: partitionKey}
	result[sk] = &types.AttributeValueMemberS{Value: sortKey}
	_, err := d.DynamoClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName:           aws.String(tableName),
		Key:                 result,
		ConditionExpression: aws.String(pkExists),
	})
	return err
}

// Get implements utils.Storage.
// key format: <table>;<partition key>;<sort key>;<action>
func (d *dynamo) Get(ctx context.Context, key string) (interface{}, error) {
	keys := strings.Split(key, ";")
	if len(keys) != 4 {
		return nil, ErrInvalidKey
	}
	tableName := keys[0]
	partitionKey := keys[1]
	sortKey := keys[2]
	action := keys[3]

	switch strings.ToLower(action) {
	case "get":
		result := make(map[string]types.AttributeValue)
		result[pk] = &types.AttributeValueMemberS{Value: partitionKey}
		result[sk] = &types.AttributeValueMemberS{Value: sortKey}
		data, err := d.DynamoClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
			TableName: aws.String(tableName),
			Key:       result,
		})
		if err != nil {
			return nil, errors.Join(ErrDynamoDB, err)
		}
		if data.Item == nil {
			return nil, ErrNotFound
		}
		return data.Item, nil
	case "query":
		strs := strings.Split(sortKey, " ")
		if len(strs) != 2 {
			return nil, ErrInvalidKey
		}
		switch strings.ToLower(strs[0]) {
		case "beginwith":
			keyEx := expression.KeyAnd(
				expression.Key(pk).Equal(expression.Value(partitionKey)),
				expression.KeyBeginsWith(expression.Key(sk), strs[1]))
			expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
			if err != nil {
				return nil, errors.Join(ErrDynamoDB, err)
			}
			queryPaginator := dynamodb.NewQueryPaginator(d.DynamoClient, &dynamodb.QueryInput{
				TableName:                 aws.String(tableName),
				ExpressionAttributeNames:  expr.Names(),
				ExpressionAttributeValues: expr.Values(),
				KeyConditionExpression:    expr.KeyCondition(),
			})
			if !queryPaginator.HasMorePages() {
				return nil, ErrNotFound
			}
			return queryPaginator, nil
		default:
			return nil, ErrNotSupport
		}

	default:
		return nil, ErrNotSupport
	}
}

// Save implements utils.Storage.
// key format: <table>;<partition key>;<sort key>
func (d *dynamo) Save(ctx context.Context, key string, value interface{}) error {
	keys := strings.Split(key, ";")
	if len(keys) != 3 {
		return ErrInvalidKey
	}
	tableName := keys[0]
	partitionKey := keys[1]
	sortKey := keys[2]
	item, err := attributevalue.MarshalMap(value)
	if err != nil {
		return err
	}
	item[pk] = &types.AttributeValueMemberS{
		Value: partitionKey,
	}
	item[sk] = &types.AttributeValueMemberS{
		Value: sortKey,
	}
	_, err = d.DynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(tableName),
		Item:                item,
		ConditionExpression: aws.String(pkNotExists),
	})
	return err
}

// Update implements utils.Storage.
// key format: <table>;<partition key>;<sort key>
func (d *dynamo) Update(ctx context.Context, key string, value interface{}, updateMask []string) error {
	keys := strings.Split(key, ";")
	if len(keys) != 3 {
		return ErrInvalidKey
	}
	tableName := keys[0]
	partitionKey := keys[1]
	sortKey := keys[2]

	expr, err := getUpdateExpression(value, updateMask)
	if err != nil {
		return errors.Join(ErrDynamoDB, err)
	}
	result := make(map[string]types.AttributeValue)
	result[pk] = &types.AttributeValueMemberS{Value: partitionKey}
	result[sk] = &types.AttributeValueMemberS{Value: sortKey}
	_, err = d.DynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 aws.String(tableName),
		Key:                       result,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression:          expr.Update(),
		ReturnValues:              types.ReturnValueNone,
		ConditionExpression:       aws.String(pkExists),
	})
	if err != nil {
		return errors.Join(ErrDynamoDB, err)
	}
	return nil
}

func NewDynamoDB(ctx context.Context, appCfg *appConfig.StorageConfig) (utils.Storage, error) {
	var (
		err error
		cfg aws.Config
	)
	cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(appCfg.Region))
	if err != nil {
		return nil, err
	}
	return &dynamo{DynamoClient: dynamodb.NewFromConfig(cfg)}, nil
}

func NewDevDynamoDB(ctx context.Context, appCfg *appConfig.StorageConfig) (utils.Storage, error) {
	cfg, _ := config.LoadDefaultConfig(ctx,
		config.WithRegion(appCfg.Region),
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
			func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{URL: fmt.Sprintf("http://%s:%d", appCfg.Host, appCfg.Port)}, nil
			})),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	return &dynamo{DynamoClient: dynamodb.NewFromConfig(cfg)}, nil
}

func getUpdateExpression(in interface{}, updateMask []string) (expression.Expression, error) {
	var (
		vals   = reflect.ValueOf(in)
		start  = true
		update expression.UpdateBuilder
	)
	for _, key := range updateMask {
		sKey := utils.ToCamelCase(key)

		if vals.FieldByName(sKey).IsValid() {
			if start {
				update = expression.Set(expression.Name(key), expression.Value(vals.FieldByName(sKey).Interface()))
				start = false
			} else {
				update.Set(expression.Name(key), expression.Value(vals.FieldByName(sKey).Interface()))
			}
		}
	}
	return expression.NewBuilder().WithUpdate(update).Build()
}

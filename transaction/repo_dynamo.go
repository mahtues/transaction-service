package transaction

import (
	"github.com/pkg/errors"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type DynamoRepository struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoRepository(client *dynamodb.DynamoDB, tableName string) *DynamoRepository {
	d := &DynamoRepository{}
	d.Init(client, tableName)
	return d
}

func (d *DynamoRepository) Init(client *dynamodb.DynamoDB, tableName string) {
	*d = DynamoRepository{
		client:    client,
		tableName: tableName,
	}
}

func (d *DynamoRepository) SaveTransaction(transaction Transaction) error {
	av, err := dynamodbattribute.MarshalMap(transaction)
	if err != nil {
		return errors.Wrap(err, "transaction marshal")
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(d.tableName),
	}

	_, err = d.client.PutItem(input)
	if err != nil {
		return errors.Wrap(err, "put transaction")
	}

	return nil
}

func (d *DynamoRepository) LoadTransaction(id string) (Transaction, error) {
	output, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {S: aws.String(id)},
		},
	})
	if err != nil {
		return Transaction{}, errors.Wrap(err, "get transaction")
	}

	if output.Item == nil {
		return Transaction{}, errors.New("transaction not found")
	}

	var transaction Transaction

	if err = dynamodbattribute.UnmarshalMap(output.Item, &transaction); err != nil {
		return Transaction{}, errors.Wrap(err, "unmarshal transaction")
	}

	return transaction, nil
}

func CreateTable(client *dynamodb.DynamoDB, name string) error {
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(name),
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("id"),
				KeyType:       aws.String("HASH"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	}

	if _, err := client.CreateTable(input); err != nil {
		return errors.Wrap(err, "error creating table")
	}

	return nil
}

func DeleteTable(client *dynamodb.DynamoDB, name string) error {
	input := &dynamodb.DeleteTableInput{
		TableName: aws.String(name),
	}

	if _, err := client.DeleteTable(input); err != nil {
		return errors.Wrap(err, "error deleting table")
	}

	return nil
}

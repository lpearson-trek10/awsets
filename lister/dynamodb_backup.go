package lister

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/trek10inc/awsets/option"
	"github.com/trek10inc/awsets/resource"
)

type AWSDynamoDBBackup struct {
}

func init() {
	i := AWSDynamoDBBackup{}
	listers = append(listers, i)
}

func (l AWSDynamoDBBackup) Types() []resource.ResourceType {
	return []resource.ResourceType{
		resource.DynamoDbBackup,
	}
}

func (l AWSDynamoDBBackup) List(cfg option.AWSetsConfig) (*resource.Group, error) {
	svc := dynamodb.NewFromConfig(cfg.AWSCfg)

	rg := resource.NewGroup()
	err := Paginator(func(nt *string) (*string, error) {
		res, err := svc.ListBackups(cfg.Context, &dynamodb.ListBackupsInput{
			BackupType:              types.BackupTypeFilterAll,
			Limit:                   aws.Int32(100),
			ExclusiveStartBackupArn: nt,
		})
		if err != nil {
			return nil, err
		}
		for _, bk := range res.BackupSummaries {
			bkr := resource.New(cfg, resource.DynamoDbBackup, bk.BackupName, bk.BackupName, bk)
			bkr.AddRelation(resource.DynamoDbTable, bk.TableName, "")
			rg.AddResource(bkr)
		}
		return res.LastEvaluatedBackupArn, nil
	})
	return rg, err
}

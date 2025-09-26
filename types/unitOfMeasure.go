package types

import (
	"context"
	"encoding/json"

	"rhize-data-collection-import/domain"

	"github.com/hasura/go-graphql-client"
)

func CreateUnitOfMeasure(ctx context.Context, client *graphql.Client, input []domain.AddUnitOfMeasureInput) error {
	var m struct {
		AddUnitOfMeasure struct {
			UnitOfMeasure struct {
				ID       string `graphql:"id"`
				DataType string `graphql:"dataType"`
			} `graphql:"unitOfMeasure"`
		} `graphql:"addUnitOfMeasure(input: $uom)"`
	}

	reqVar := map[string]interface{}{
		"uom": input,
	}
	jsonBytes, err := client.MutateRaw(ctx, &m, reqVar)
	if err != nil {
		return err
	}

	var response struct {
		AddUnitOfMeasure struct {
			UnitOfMeasure []*struct {
				ID       string           `json:"id,omitempty"`
				DataType *domain.DataType `json:"dataType,omitempty"`
			} `json:"unitOfMeasure"`
		} `json:"addUnitOfMeasure"`
	}

	err = json.Unmarshal(jsonBytes, &response)

	return err
}

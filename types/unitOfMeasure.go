package types

import (
	"context"
	"encoding/json"

	"rhize-data-collection-import/domain"

	"github.com/hasura/go-graphql-client"
)

func GetUnitOfMeasure(ctx context.Context, client *graphql.Client, id string) (*domain.UnitOfMeasure, error) {
	var query struct {
		GetUnitOfMeasure struct {
			ID       string           `graphql:"id"`
			DataType *domain.DataType `graphql:"dataType"`
		} `graphql:"getUnitOfMeasure(id: $id)"`
	}

	variables := map[string]interface{}{
		"id": id,
	}

	bytes, err := client.QueryRaw(ctx, &query, variables)
	if err != nil {
		return nil, err
	}

	var response struct {
		GetUnitOfMeasure *domain.UnitOfMeasure `json:"getUnitOfMeasure"`
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, err
	}

	return response.GetUnitOfMeasure, nil
}

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

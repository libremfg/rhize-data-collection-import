package types

import (
	"context"
	"encoding/json"

	"rhize-data-collection-import/domain"

	"github.com/hasura/go-graphql-client"
)

func GetDataSource(ctx context.Context, client *graphql.Client, label string) (*domain.DataSource, error) {

	var response struct {
		GetDataSource *domain.DataSource `json:"getDataSource"`
	}

	var q struct {
		GetDataSource struct {
			ID            string `graphql:"id"`
			Label         string `graphql:"label"`
			IID           string `graphql:"iid"`
			ActiveVersion struct {
				ID     string `graphql:"id"`
				IID    string `graphql:"iid"`
				Topics []struct {
					Iid         string `graphql:"iid"`
					Label       string `graphql:"label"`
					Description string `graphql:"description"`
				} `graphql:"topics(first: 20000)"`
				ConnectionString string `graphql:"connectionString"`
			} `graphql:"activeVersion"`
			Versions []struct {
				IID           string `graphql:"iid"`
				ID            string `graphql:"id"`
				Version       string `graphql:"version"`
				VersionStatus string `graphql:"versionStatus"`
			} `graphql:"versions(filter:{versionStatus:{eq:ACTIVE}},order:{asc:version}, first:1)"`
		} `graphql:"getDataSource(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(label),
	}
	jsonResult, err := client.QueryRaw(context.Background(), &q, variables)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		return nil, err
	}

	return response.GetDataSource, nil
}

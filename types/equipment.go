package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"rhize-data-collection-import/domain"

	"github.com/hasura/go-graphql-client"
)

func GetEquipmentPayload(name string, description string, level domain.EquipmentElementLevel, uiSortIndex int) *domain.AddEquipmentInput {
	id := name

	index := uiSortIndex
	return &domain.AddEquipmentInput{
		ID:          id,
		Label:       name,
		UISortIndex: &index,
		Versions: []*domain.EquipmentVersionRef{
			{
				Description:    &description,
				DisplayName:    &name,
				EquipmentLevel: &level,
				ID:             StringPtr(id),
				TimeZoneName:   StringPtr("America/New_York"),
				Version:        StringPtr("1"),
				VersionStatus:  &versionStateActive,
			},
		},
	}
}

func CreateEquipment(ctx context.Context, client *graphql.Client, equipmentInput *domain.AddEquipmentInput) error {
	var err error
	var jsonResult []byte

	var m struct {
		AddEquipment struct {
			Equipment struct {
				Id  string `graphql:"id"`
				Iid string `graphql:"iid"`
			} `graphql:"equipment"`
		} `graphql:"addEquipment(input: $equipment)"`
	}

	type ids struct {
		ID  string `json:"id"`
		Iid string `json:"iid"`
	}

	var response struct {
		AddEquipment struct {
			Equipment []ids `json:"equipment"`
		} `json:"addEquipment"`
	}

	variables := map[string]interface{}{
		"equipment": []domain.AddEquipmentInput{
			*equipmentInput,
		},
	}

	jsonResult, err = client.NamedMutateRaw(ctx, "AddEquipment", m, variables)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		panic(err)
	}

	if len(response.AddEquipment.Equipment) == 0 {
		msg := fmt.Sprintf("expected equipment in %v", response)
		return errors.New(msg)
	}

	updateEquipment := GetEquipment(ctx, client, equipmentInput)

	if len(updateEquipment.Versions) == 0 {
		panic("expected at least 1 version")
	}

	var updateMutation struct {
		UpdateEquipment struct {
			equipment struct {
				id string `graphql:"id"`
			} `graphql:"equipment"`
		} `graphql:"updateEquipment(input: $input)"`
	}

	var updateResponse struct {
		UpdateEquipment struct {
			Equipment []struct {
				ID string `json:"id"`
			} `json:"equipment"`
		} `json:"updateEquipment"`
	}

	vars := map[string]interface{}{
		"input": domain.UpdateEquipmentInput{
			Filter: &domain.EquipmentFilter{
				ID: &domain.StringExactFilterStringFullTextFilterStringRegExpFilter{
					Eq: &updateEquipment.ID,
				},
			},
			Set: &domain.EquipmentPatch{
				ActiveVersion: &domain.EquipmentVersionRef{
					Iid: &updateEquipment.Versions[0].Iid,
				},
			},
		},
	}

	jsonResult, err = client.NamedMutateRaw(ctx, "UpdateEquipment", updateMutation, vars)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &updateResponse)
	if err != nil {
		panic(err)
	}

	if len(updateResponse.UpdateEquipment.Equipment) == 0 {
		msg := fmt.Sprintf("expected equipment in %v", updateResponse)
		return errors.New(msg)
	}

	return nil
}

func GetEquipment(ctx context.Context, client *graphql.Client, equipment *domain.AddEquipmentInput) *domain.Equipment {

	var response struct {
		GetEquipment *domain.Equipment `json:"getEquipment"`
	}

	var q struct {
		GetEquipment struct {
			ID            string `graphql:"id"`
			Label         string `graphql:"label"`
			IID           string `graphql:"iid"`
			ActiveVersion struct {
				IID string `graphql:"iid"`
			} `graphql:"activeVersion"`
			Versions []struct {
				IID           string `graphql:"iid"`
				ID            string `graphql:"id"`
				Version       string `graphql:"version"`
				VersionStatus string `graphql:"versionStatus"`
			} `graphql:"versions(filter:{versionStatus:{eq:ACTIVE}},order:{asc:version}, first:1)"`
		} `graphql:"getEquipment(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(equipment.ID),
	}
	jsonResult, err := client.QueryRaw(context.Background(), &q, variables)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		panic(err)
	}

	return response.GetEquipment
}

package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"rhize-data-collection-import/domain"

	"github.com/hasura/go-graphql-client"
)

var versionStateActive = domain.VersionStateActive

func GetEquipmentClassPayload(name string, description *string, level *domain.EquipmentElementLevel, uiSortIndex int) *domain.AddEquipmentClassInput {
	id := name
	index := uiSortIndex

	return &domain.AddEquipmentClassInput{
		ID:          id,
		Label:       name,
		UISortIndex: &index,
		Versions: []*domain.EquipmentClassVersionRef{
			{
				Description:    description,
				DisplayName:    &name,
				EquipmentLevel: level,
				ID:             StringPtr(id),
				Version:        StringPtr("1"),
				VersionStatus:  &versionStateActive,
			},
		},
	}
}

func CreateEquipmentClass(ctx context.Context, client *graphql.Client, equipmentClassInput *domain.AddEquipmentClassInput) (string, error) {
	var err error
	var jsonResult []byte

	var m struct {
		AddEquipmentClass struct {
			equipment struct {
				id       string `graphql:"id"`
				Versions []struct {
					IID        string `graphql:"iid"`
					Properties []struct {
						ID           string `graphql:"id"`
						IID          string `graphql:"iid"`
						Label        string `graphql:"label"`
						BindingType  string `graphql:"bindingType"`
						PropertyType string `graphql:"propertyType"`
					} `graphql:"properties"`
				} `graphql:"versions"`
			} `graphql:"equipmentClass"`
		} `graphql:"addEquipmentClass(input: $equipmentClass)"`
	}

	type ids struct {
		ID string `json:"id"`
	}

	var response struct {
		AddEquipmentClass struct {
			EquipmentClass []ids `json:"equipmentClass"`
		} `json:"addEquipmentClass"`
	}

	variables := map[string]interface{}{
		"equipmentClass": []domain.AddEquipmentClassInput{
			*equipmentClassInput,
		},
	}

	jsonResult, err = client.NamedMutateRaw(ctx, "AddEquipmentClass", m, variables)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		panic(err)
	}

	if len(response.AddEquipmentClass.EquipmentClass) == 0 {
		msg := fmt.Sprintf("expected equipment class in %v", response)
		return "", errors.New(msg)
	}

	updateEquipmentClass := GetEquipmentClass(ctx, client, equipmentClassInput)

	if len(updateEquipmentClass.Versions) == 0 {
		panic("expected at least 1 version")
	}

	var updateMutation struct {
		UpdateEquipmentClass struct {
			equipmentClass struct {
				id  string `graphql:"id"`
				iid string `graphql:"iid"`
			} `graphql:"equipmentClass"`
		} `graphql:"updateEquipmentClass(input: $input)"`
	}

	var updateResponse struct {
		UpdateEquipmentClass struct {
			EquipmentClass []struct {
				ID  string `json:"id"`
				Iid string `json:"iid"`
			} `json:"equipmentClass"`
		} `json:"updateEquipmentClass"`
	}

	vars := map[string]interface{}{
		"input": domain.UpdateEquipmentClassInput{
			Filter: &domain.EquipmentClassFilter{
				ID: &domain.StringExactFilterStringFullTextFilterStringRegExpFilter{
					Eq: &updateEquipmentClass.ID,
				},
			},
			Set: &domain.EquipmentClassPatch{
				ActiveVersion: &domain.EquipmentClassVersionRef{
					Iid: &updateEquipmentClass.Versions[0].Iid,
				},
			},
		},
	}

	jsonResult, err = client.NamedMutateRaw(ctx, "UpdateEquipmentClass", updateMutation, vars)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &updateResponse)
	if err != nil {
		panic(err)
	}

	if len(updateResponse.UpdateEquipmentClass.EquipmentClass) == 0 {
		msg := fmt.Sprintf("expected equipment class in %v", updateResponse)
		return "", errors.New(msg)
	}

	for _, property := range equipmentClassInput.Versions[0].Properties {
		if property.BindingType == nil {
			panic("expected binding type")
		}
	}

	return updateResponse.UpdateEquipmentClass.EquipmentClass[0].Iid, nil
}

func GetEquipmentClass(ctx context.Context, client *graphql.Client, equipmentClass *domain.AddEquipmentClassInput) *domain.EquipmentClass {

	var response struct {
		GetEquipmentClass *domain.EquipmentClass `json:"getEquipmentClass"`
	}

	var q struct {
		GetEquipmentClass struct {
			ID            string `graphql:"id"`
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
		} `graphql:"getEquipmentClass(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(equipmentClass.ID),
	}

	jsonResult, err := client.QueryRaw(context.Background(), &q, variables)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		panic(err)
	}

	return response.GetEquipmentClass
}
func GetEquipmentClassVersion(ctx context.Context, client *graphql.Client, equipmentClassId string, version string) (*domain.EquipmentClassVersion, error) {

	var response struct {
		GetEquipmentClassVersion *domain.EquipmentClassVersion `json:"getEquipmentClassVersion"`
	}

	var q struct {
		GetEquipmentClassVersion struct {
			IID        string `graphql:"iid"`
			ID         string `graphql:"id"`
			Version    string `graphql:"version"`
			Properties []struct {
				ID                 string `graphql:"id"`
				IID                string `graphql:"iid"`
				Label              string `graphql:"label"`
				Description        string `graphql:"description"`
				Value              string `graphql:"value"`
				BindingType        string `graphql:"bindingType"`
				PropertyType       string `graphql:"propertyType"`
				ValueUnitOfMeasure struct {
					ID       string `graphql:"id"`
					DataType string `graphql:"dataType"`
				} `graphql:"valueUnitOfMeasure"`
				Parent struct {
					Iid string `graphql:"iid"`
					ID  string `graphql:"id"`
				} `graphql:"parent"`
			} `graphql:"properties"`
		} `graphql:"getEquipmentClassVersion(id:$id, version:$version)"`
	}

	variables := map[string]interface{}{
		"id":      graphql.String(equipmentClassId),
		"version": graphql.String(version),
	}

	jsonResult, err := client.QueryRaw(context.Background(), &q, variables)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		return nil, err
	}

	return response.GetEquipmentClassVersion, nil
}
func GetEquipmentClassAllVersions(ctx context.Context, client *graphql.Client, equipmentClass *domain.AddEquipmentClassInput) *domain.EquipmentClass {

	var response struct {
		GetEquipmentClass *domain.EquipmentClass `json:"getEquipmentClass"`
	}

	var q struct {
		GetEquipmentClass struct {
			ID            string `graphql:"id"`
			IID           string `graphql:"iid"`
			ActiveVersion struct {
				IID string `graphql:"iid"`
			} `graphql:"activeVersion"`
			Versions []struct {
				IID           string `graphql:"iid"`
				ID            string `graphql:"id"`
				Version       string `graphql:"version"`
				VersionStatus string `graphql:"versionStatus"`
				Properties    []struct {
					ID                 string `graphql:"id"`
					IID                string `graphql:"iid"`
					Label              string `graphql:"label"`
					Description        string `graphql:"description"`
					Value              string `graphql:"value"`
					BindingType        string `graphql:"bindingType"`
					PropertyType       string `graphql:"propertyType"`
					ValueUnitOfMeasure struct {
						ID       string `graphql:"id"`
						DataType string `graphql:"dataType"`
					} `graphql:"valueUnitOfMeasure"`
					Parent struct {
						Iid string `graphql:"iid"`
						ID  string `graphql:"id"`
					} `graphql:"parent"`
				} `graphql:"properties"`
			} `graphql:"versions(order:{asc:version})"`
		} `graphql:"getEquipmentClass(id:$id)"`
	}

	variables := map[string]interface{}{
		"id": graphql.String(equipmentClass.ID),
	}

	jsonResult, err := client.QueryRaw(context.Background(), &q, variables)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		panic(err)
	}

	return response.GetEquipmentClass
}

func CreateEquipmentClassProperty(ctx context.Context, client *graphql.Client, equipmentClassPropertyInput *domain.AddEquipmentClassPropertyInput) error {
	var err error
	var jsonResult []byte

	var m struct {
		AddEquipmentClassProperty struct {
			NumUids int `graphql:"numUids"`
		} `graphql:"addEquipmentClassProperty(input: $equipmentClassProperty, upsert: true)"`
	}

	var response struct {
		AddEquipmentClassProperty struct {
			NumUids int `json:"numUids"`
		} `json:"addEquipmentClassProperty"`
	}

	variables := map[string]interface{}{
		"equipmentClassProperty": []domain.AddEquipmentClassPropertyInput{
			*equipmentClassPropertyInput,
		},
	}

	jsonResult, err = client.NamedMutateRaw(ctx, "AddEquipmentClassProperty", m, variables)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		return err
	}

	return nil
}

func DeleteEquipmentClassProperty(ctx context.Context, client *graphql.Client, equipmentClassPropertyIid string) error {
	var m struct {
		DeleteEquipmentClassProperty struct {
			NumUids int `graphql:"numUids"`
		} `graphql:"deleteEquipmentClassProperty(filter: {iid: $iid})"`
	}

	variables := map[string]interface{}{
		"iid": []graphql.ID{
			graphql.ID(equipmentClassPropertyIid),
		},
	}

	bytes, err := client.MutateRaw(ctx, m, variables)
	if err != nil {
		return err
	}

	var response struct {
		DeleteEquipmentClassProperty struct {
			NumUids int `json:"numUids"`
		} `json:"deleteEquipmentClassProperty"`
	}

	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return err
	}

	if response.DeleteEquipmentClassProperty.NumUids != 1 {
		return fmt.Errorf("expected to delete 1 equipment class property, deleted %d", response.DeleteEquipmentClassProperty.NumUids)
	}

	return nil
}

func SaveEquipmentClassVersionAs(ctx context.Context, client *graphql.Client, equipmentClassId string, equipmentClassVersion string, equipmentClassVersionTo string) (string, error) {
	var err error
	var jsonResult []byte

	var m struct {
		SaveEquipmentClassVersionAs struct {
			Iid string `graphql:"iid"`
		} `graphql:"saveEquipmentClassVersionAs(fromID: $fromId, fromVersion: $fromVersion, toID: $toId, toVersion: $toVersion)"`
	}

	var response struct {
		SaveEquipmentClassVersionAs []struct {
			Iid string `json:"iid"`
		} `json:"saveEquipmentClassVersionAs"`
	}

	variables := map[string]interface{}{
		"fromId":      equipmentClassId,
		"fromVersion": equipmentClassVersion,
		"toId":        equipmentClassId,
		"toVersion":   equipmentClassVersionTo,
	}

	jsonResult, err = client.NamedMutateRaw(ctx, "SaveEquipmentClassVersionAs", m, variables)
	if err != nil {
		return "", err
	}

	err = json.Unmarshal(jsonResult, &response)
	if err != nil {
		return "", err
	}

	return response.SaveEquipmentClassVersionAs[0].Iid, err
}

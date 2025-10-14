package importSheet

import (
	"context"
	"log"
	"strconv"

	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
)

func EquipmentModel(ctx context.Context, client *graphql.Client, equipmentImportData ImportData) {
	if equipmentImportData.Datasource == "" {
		log.Printf("\tNo Datasource provided, skipping Equipment bindings")
		return
	}
	setupEquipment(ctx, client, equipmentImportData.EquipmentImportData, equipmentImportData.EquipmentClassImportData.EquipmentClassName, equipmentImportData.Datasource)
}

func setupEquipment(ctx context.Context, client *graphql.Client, equipmentImportData []EquipmentImportData, equipmentClass string, datasource string) {
	log.Printf("\tSetting up Equipment bindings")

	// First check that the Datasource exists, we cannot bind without it
	ds, err := types.GetDataSource(ctx, client, datasource)
	if err != nil {
		log.Printf("\t\tCould not get Datasource \"%s\", skipping Equipment setup. Error: %s\n", datasource, err)
		return
	}
	if ds == nil {
		log.Printf("\t\tDatasource \"%s\" not found, add Datasource and rerun utility.\n", datasource)
		return
	}

	for _, equipment := range equipmentImportData {
		log.Printf("\t\tAdding bindings for Equipment \"%s\"\n", equipment.EquipmentName)

		equipmentVersions := types.GetEquipmentAllVersions(ctx, client, equipment.EquipmentName)
		// Must ensure that equipment exists
		if equipmentVersions == nil {
			// If it does not exist, log a warning
			log.Printf("\t\t\tEquipment with ID \"%s\" does not exist, make this Equipment and run the utility again\n", equipment.EquipmentName)
			continue
		}
		latestVersion := pickLatestEquipmentVersion(equipmentVersions.Versions, false)
		// And that the equipment class exists on that equipment
		found := false
		for _, ec := range latestVersion.EquipmentClasses {
			if ec.ID == equipmentClass {
				found = true
				break
			}
		}
		if !found {
			log.Printf("\t\t\tEquipment with ID \"%s\" and Version \"%s\" does not have Equipment Class with ID \"%s\", add this class to the equipment and run the utility again\n", equipment.EquipmentName, latestVersion.Version, equipmentClass)
			continue
		}
		// And that datasource exists on that equipment
		found = false
		for _, dataSource := range latestVersion.DataSources {
			if dataSource.DataSource != nil && dataSource.DataSource.ID == datasource {
				found = true
				break
			}
		}
		if !found {
			log.Printf("\t\t\tEquipment with ID \"%s\" and Version \"%s\" does not have DataSource with ID \"%s\", add this datasource to the equipment and run the utility again\n", equipment.EquipmentName, latestVersion.Version, datasource)
			continue
		}

		propertyNameAliases := make([]*domain.PropertyNameAliasRef, 0)
		for _, binding := range equipment.EquipmentTagBindings {
			// Make sure that the topic exists in the Datasource
			found := false
			for _, topic := range ds.ActiveVersion.Topics {
				if topic.Label == binding.Tag {
					found = true
					break
				}
			}
			if !found {
				log.Printf("\t\t\tCould not find topic \"%s\" inside of Datasource \"%s\", skipping this binding", binding.Tag, datasource)
				continue
			}

			// Check that property is not already binded, but only if not an active version
			// If binded, recommend removing the binding in the UI
			if latestVersion.VersionStatus != domain.VersionStateActive {
				found = false
				for _, alias := range latestVersion.PropertyNameAliases {
					if alias.PropertyLabel == binding.PropertyID {
						found = true
						break
					}
				}
				if found {
					// To-Do: Grab iid of alias and run a delete on it
					log.Printf("\t\t\tBinding for Property \"%s\" already exists, skipping this binding.", binding.PropertyID)
					continue
				}
			}

			propertyNameAliases = append(propertyNameAliases, &domain.PropertyNameAliasRef{
				DataSource: &domain.DataSourceRef{
					ID: types.StringPtr(datasource),
				},
				DataSourceTopicLabel: types.StringPtr(binding.Tag),
				PropertyLabel:        types.StringPtr(binding.PropertyID),
			})
		}

		// If latest version is an Active Version, then make a new draft version
		if latestVersion.VersionStatus == domain.VersionStateActive {
			// For now just logging a warning
			log.Printf("\t\t\tLatest version of Equipment with ID \"%s\" is an Active Version, please create a new Draft Version and rerun the utility\n", equipment.EquipmentName)
			continue
		} else {
			if len(propertyNameAliases) == 0 {
				log.Println("\t\t\tNo available topics to bind, skipping binding")
				continue
			}

			// Assume Draft
			err := types.SetEquipmentBinds(ctx, client, equipment.EquipmentName, latestVersion.Version, datasource, propertyNameAliases)
			if err != nil {
				log.Panicln(err)
			}
		}
	}
}

func pickLatestEquipmentVersion(versions []*domain.EquipmentVersion, draftOnly bool) *domain.EquipmentVersion {

	var thisVersions []*domain.EquipmentVersion

	if draftOnly {
		for _, version := range versions {
			if version.VersionStatus == domain.VersionStateDraft {
				thisVersions = append(thisVersions, version)
			}
		}
		return nil
	} else {
		thisVersions = make([]*domain.EquipmentVersion, len(versions))
		copy(thisVersions, versions)
	}

	if len(thisVersions) == 0 {
		return nil
	}

	if len(thisVersions) == 1 {
		return thisVersions[0]
	}

	index := 0
	latestVersionNum := thisVersions[index].Version

	for i, version := range thisVersions[1:] {
		versionNum, err := strconv.Atoi(version.Version)
		if err != nil {
			// compare strings if not integers
			if version.Version > latestVersionNum {
				latestVersionNum = version.Version
				index = i
			}
		}
		currentNum, err := strconv.Atoi(latestVersionNum)
		if err != nil {
			// compare strings if not integers
			if version.Version > latestVersionNum {
				latestVersionNum = version.Version
				index = i
			}
		}
		if versionNum > currentNum {
			latestVersionNum = version.Version
			index = i + 1
		}
	}

	return thisVersions[index]
}

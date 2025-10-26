package commands

import (
	"context"
	"log"
	"strconv"

	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
)

var (
	EquipmentCmd = &cobra.Command{
		Use:     "equipment",
		Short:   "Import equipment",
		Aliases: []string{"e", "eq"},
		Run:     importEquipment,
	}
)

func importEquipment(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log.Println("Starting import for Equipment")
	equipment(ctx, Client, ImportData)
	log.Println("Finished import for Equipment")
}

// func importEquipment(ctx context.Context, client *graphql.Client, equipmentImportData []EquipmentImportData, equipmentClass string, datasource string) {
func equipment(ctx context.Context, client *graphql.Client, importData types.ImportData) {
	log.Printf("\tSetting up Equipment bindings")

	// First check that the Datasource exists, we cannot bind without it
	ds, err := types.GetDataSource(ctx, client, importData.Datasource)
	if err != nil {
		log.Printf("\t\tCould not get Datasource \"%s\", skipping Equipment setup. Error: %s\n", importData.Datasource, err)
		return
	}
	if ds == nil {
		log.Printf("\t\tDatasource \"%s\" not found, add Datasource and rerun utility.\n", importData.Datasource)
		return
	}

	for _, equipment := range importData.Equipment {
		log.Printf("\t\tAdding bindings for Equipment \"%s\"\n", equipment.ID)

		equipmentVersions := types.GetEquipmentAllVersions(ctx, client, equipment.ID)
		// Must ensure that equipment exists
		if equipmentVersions == nil {
			// If it does not exist, log a warning
			log.Printf("\t\t\tEquipment with ID \"%s\" does not exist, make this Equipment and run the utility again\n", equipment.ID)
			continue
		}
		latestVersion := pickLatestEquipmentVersion(equipmentVersions.Versions, false)
		// And that the equipment class exists on that equipment
		found := false
		for _, ec := range latestVersion.EquipmentClasses {
			if ec.ID == importData.EquipmentClass.Label {
				found = true
				break
			}
		}
		if !found {
			log.Printf("\t\t\tEquipment with ID \"%s\" and Version \"%s\" does not have Equipment Class with ID \"%s\", add this class to the equipment and run the utility again\n", equipment.ID, latestVersion.Version, equipmentClass)
			continue
		}
		// And that datasource exists on that equipment
		found = false
		for _, dataSource := range latestVersion.DataSources {
			if dataSource.DataSource != nil && dataSource.DataSource.ID == importData.Datasource {
				found = true
				break
			}
		}
		if !found {
			log.Printf("\t\t\tEquipment with ID \"%s\" and Version \"%s\" does not have DataSource with ID \"%s\", add this datasource to the equipment and run the utility again\n", equipment.ID, latestVersion.Version, importData.Datasource)
			continue
		}

		propertyNameAliases := make([]*domain.PropertyNameAliasRef, 0)
		for _, binding := range equipment.TagBindings {
			// Make sure that the topic exists in the Datasource
			found := false
			for _, topic := range ds.ActiveVersion.Topics {
				if topic.Label == binding.Tag {
					found = true
					break
				}
			}
			if !found {
				log.Printf("\t\t\tCould not find topic \"%s\" inside of Datasource \"%s\", skipping this binding", binding.Tag, importData.Datasource)
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

			propertyNameAlias := domain.PropertyNameAliasRef{
				DataSource: &domain.DataSourceRef{
					ID: types.StringPtr(importData.Datasource),
				},
				DataSourceTopicLabel: types.StringPtr(binding.Tag),
				PropertyLabel:        types.StringPtr(binding.PropertyID),
			}
			if binding.Expression != "" {
				propertyNameAlias.Expression = types.StringPtr(binding.Expression)
			}

			propertyNameAliases = append(propertyNameAliases, &propertyNameAlias)
		}

		// If latest version is an Active Version, then make a new draft version
		if latestVersion.VersionStatus == domain.VersionStateActive {
			// For now just logging a warning
			log.Printf("\t\t\tLatest version of Equipment with ID \"%s\" is an Active Version, please create a new Draft Version and rerun the utility\n", equipment.ID)
			continue
		} else {
			if len(propertyNameAliases) == 0 {
				log.Println("\t\t\tNo available topics to bind, skipping binding")
				continue
			}

			// Assume Draft
			err := types.SetEquipmentBinds(ctx, client, equipment.ID, latestVersion.Version, importData.Datasource, propertyNameAliases)
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

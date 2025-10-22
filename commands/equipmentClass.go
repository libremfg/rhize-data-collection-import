package commands

import (
	"context"
	"fmt"
	"log"
	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"
	"strconv"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
)

var (
	EquipmentClassCmd = &cobra.Command{
		Use:     "equipmentClass",
		Short:   "Import Equipment Class from file",
		Aliases: []string{"ec"},
		Run:     importEquipmentClass,
	}
)

func init() {
	// Setup Flags
}

func importEquipmentClass(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log.Println("Starting import for Equipment Class")
	equipmentClass(ctx, Client, ImportData.EquipmentClass)
	log.Println("Finished import for Equipment Class")
}

func equipmentClass(ctx context.Context, client *graphql.Client, equipmentClassData types.ImportEquipmentClass) {
	properties := make([]*domain.EquipmentClassPropertyRef, 0)

	bound := domain.PropertyBindingTypeBound
	static := domain.PropertyBindingTypeStatic
	propertyType := domain.Isa95PropertyTypeDefaultType

	equipmentClassName := equipmentClassData.Label

	for _, property := range equipmentClassData.Properties {
		if !property.Use {
			continue
		}

		propertyPath := property.ID
		if strings.HasPrefix(propertyPath, equipmentClassName) {
			propertyPath = propertyPath[len(equipmentClassName)+1:]
		}

		// Filter out initial starting dot if it exists
		if propertyPath[0] == '.' {
			propertyPath = propertyPath[1:]
		}
		propertiesSplit := strings.Split(propertyPath, ".")

		if len(propertiesSplit) == 0 {
			continue
		}

		// Get Property ID/Label
		propertyId := fmt.Sprintf("%s.1.%s", equipmentClassName, propertyPath)
		propertyName := propertiesSplit[len(propertiesSplit)-1]

		// Get Parent Property ID
		parentPropertyId := ""
		if len(propertiesSplit) > 1 {
			parentPropertyId = fmt.Sprintf("%s.1.%s", equipmentClassName, strings.Join(propertiesSplit[:len(propertiesSplit)-1], "."))
		}

		// Get UoM
		uom := property.UnitOfMeasure.ID

		bindingType := &bound
		if property.UnitOfMeasure.DataType == "" {
			bindingType = &static
		}

		// Create Property Ref
		property := domain.EquipmentClassPropertyRef{
			ID:           types.StringPtr(propertyId),
			Label:        types.StringPtr(propertyName),
			BindingType:  bindingType,
			PropertyType: &propertyType,
		}
		if uom != "" {
			property.ValueUnitOfMeasure = &domain.UnitOfMeasureRef{
				ID: types.StringPtr(uom),
			}
		}
		if len(propertiesSplit) > 1 {
			property.Parent = &domain.EquipmentClassPropertyRef{
				ID: types.StringPtr(parentPropertyId),
			}
		}
		properties = append(properties, &property)
	}

	// Ensure Parent Properties are present
search:
	for _, property := range properties {
		if property.Parent == nil {
			continue
		}
		for _, p := range properties {
			// Check if ID is present, if so just continue
			if *property.Parent.ID == *p.ID {
				continue search
			}
		}
		// If not present, then crash
		log.Fatalf("Could not parse equipment properties, parent property \"%s\" is missing. Check that it is defined in the spreadsheet and in Use (X).", *property.Parent.ID)
	}

	// Add Equipment Class
	log.Println("\tAdding Equipment Class")

	// Check if Equipment Class already exists
	equipmentClassId := equipmentClassName
	var equipmentClassVersion string

	equipmentClass := types.GetEquipmentClassAllVersions(ctx, client, &domain.AddEquipmentClassInput{
		ID: equipmentClassId,
	})

	// If exists then
	if equipmentClass != nil {
		log.Println("\t\tEquipment Class already exists, updating instead")
		// - Get latest version, if a draft version, then update that version
		// - Else if not a draft version, then make a new draft version and update that version
		latestVersion := pickLatestVersion(equipmentClass.Versions, false)

		if latestVersion.VersionStatus == domain.VersionStateDraft {

			prefix := len(equipmentClassName) + 3 // +3 for ".1."
			for i := range properties {
				id := *properties[i].ID
				id = equipmentClassName + "." + latestVersion.Version + "." + id[prefix:]
				properties[i].ID = &id

				if properties[i].Parent != nil {
					parentId := *properties[i].Parent.ID
					parentId = equipmentClassName + "." + latestVersion.Version + "." + parentId[prefix:]
					properties[i].Parent.ID = &parentId
				}
			}

			log.Printf("\t\tLatest version is a draft version (%s), updating that version\n", latestVersion.Version)
			for _, property := range latestVersion.Properties {
				log.Printf("\t\tRemoving property \"%s\"\n", property.ID)
				err := types.DeleteEquipmentClassProperty(ctx, client, property.Iid)
				if err != nil {
					log.Printf("\t\t\tFailed to remove property with ID \"%s\": %s\n", property.ID, err)
				}
			}
			for _, property := range properties {
				log.Printf("\t\tAdding property \"%s\"\n", *property.ID)
				err := types.CreateEquipmentClassProperty(ctx, client, &domain.AddEquipmentClassPropertyInput{
					ID:           *property.ID,
					Label:        *property.Label,
					Parent:       property.Parent,
					BindingType:  property.BindingType,
					PropertyType: *property.PropertyType,
					EquipmentClassVersion: &domain.EquipmentClassVersionRef{
						Iid: &latestVersion.Iid,
					},
					ValueUnitOfMeasure: property.ValueUnitOfMeasure,
				})
				if err != nil {
					log.Printf("\t\t\tFailed to input property with ID \"%s\": %s\n", *property.ID, err)
				}
			}
		} else {
			log.Printf("\t\tLatest version is an active version (%s), creating a new draft version\n", latestVersion.Version)

			newVersion := getNewVersion(ctx, client, equipmentClass)

			// Create a draft version
			newVersionIid, err := types.SaveEquipmentClassVersionAs(
				ctx,
				client,
				equipmentClassId,
				latestVersion.Version,
				newVersion,
			)
			if err != nil {
				panic(err)
			}

			version, err := types.GetEquipmentClassVersion(ctx, client, equipmentClassId, newVersion)
			if err != nil {
				panic(err)
			}

			prefix := len(equipmentClassName) + 3 // +3 for ".1."
			for i := range properties {
				id := *properties[i].ID
				id = equipmentClassName + "." + newVersion + "." + id[prefix:]
				properties[i].ID = &id

				if properties[i].Parent != nil {
					parentId := *properties[i].Parent.ID
					parentId = equipmentClassName + "." + newVersion + "." + parentId[prefix:]
					properties[i].Parent.ID = &parentId
				}
			}

			for _, property := range version.Properties {
				log.Printf("\t\tRemoving property \"%s\"\n", property.ID)
				err := types.DeleteEquipmentClassProperty(ctx, client, property.Iid)
				if err != nil {
					log.Printf("\t\t\tFailed to remove property with ID \"%s\": %s\n", property.ID, err)
				}
			}
			for _, property := range properties {
				log.Printf("\t\tAdding property \"%s\"\n", *property.ID)
				err := types.CreateEquipmentClassProperty(ctx, client, &domain.AddEquipmentClassPropertyInput{
					ID:           *property.ID,
					Label:        *property.Label,
					Parent:       property.Parent,
					BindingType:  property.BindingType,
					PropertyType: *property.PropertyType,
					EquipmentClassVersion: &domain.EquipmentClassVersionRef{
						Iid: &newVersionIid,
					},
					ValueUnitOfMeasure: property.ValueUnitOfMeasure,
				})
				if err != nil {
					log.Printf("\t\t\tFailed to input property with ID \"%s\": %s\n", *property.ID, err)
				}
			}
		}
	} else {
		// Else if doesn't exist, then add a new Equipment Class
		uiSortIndex := 1
		processCell := domain.EquipmentElementLevelProcessCell

		extruder := types.GetEquipmentClassPayload(equipmentClassName, types.StringPtr(equipmentClassData.Description), &processCell, uiSortIndex)

		var err error
		_, err = types.CreateEquipmentClass(ctx, client, extruder)
		if err != nil {
			panic(err)
		}

		equipmentClassVersion = "1"
		// Add Equipment Properties
		log.Println("\tAdding Equipment Properties")

		for _, property := range properties {
			err := types.CreateEquipmentClassProperty(ctx, client, &domain.AddEquipmentClassPropertyInput{
				ID:           *property.ID,
				Label:        *property.Label,
				Parent:       property.Parent,
				BindingType:  property.BindingType,
				PropertyType: *property.PropertyType,
				EquipmentClassVersion: &domain.EquipmentClassVersionRef{
					ID:      &equipmentClassName,
					Version: types.StringPtr(equipmentClassVersion),
				},
				ValueUnitOfMeasure: property.ValueUnitOfMeasure,
			})
			if err != nil {
				log.Printf("\t\tFailed to input property with ID \"%s\": %s\n", *property.ID, err)
			}
		}
	}
}

func pickLatestVersion(versions []*domain.EquipmentClassVersion, draftOnly bool) *domain.EquipmentClassVersion {

	var thisVersions []*domain.EquipmentClassVersion

	if draftOnly {
		for _, version := range versions {
			if version.VersionStatus == domain.VersionStateDraft {
				thisVersions = append(thisVersions, version)
			}
		}
		return nil
	} else {
		thisVersions = make([]*domain.EquipmentClassVersion, len(versions))
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

func getNewVersion(ctx context.Context, client *graphql.Client, equipmentClass *domain.EquipmentClass) string {
	latestVersion := pickLatestVersion(equipmentClass.Versions, false)

	latestVersionNum, err := strconv.Atoi(latestVersion.Version)
	if err != nil {
		ok := true
		latestVersionNum = len(equipmentClass.Versions) + 1
		for ok {
			equipmentClassVersion, err := types.GetEquipmentClassVersion(ctx, client, equipmentClass.ID, strconv.Itoa(latestVersionNum))
			if err != nil {
				panic(err)
			}

			if equipmentClassVersion == nil {
				ok = false
			} else {
				latestVersionNum++
			}
		}
	} else {
		latestVersionNum++
	}

	return strconv.Itoa(latestVersionNum)
}

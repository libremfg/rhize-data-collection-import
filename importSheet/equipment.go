package importSheet

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
)

func EquipmentModel(ctx context.Context, client *graphql.Client, equipmentImportData EquipmentImportData) {
	properties := make([]*domain.EquipmentClassPropertyRef, 0)

	bound := domain.PropertyBindingTypeBound
	static := domain.PropertyBindingTypeStatic
	classType := domain.Isa95PropertyTypeClassType

	equipmentClassName := equipmentImportData.EquipmentClassName

	for _, property := range equipmentImportData.EquipmentClassProperties {
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
			PropertyType: &classType,
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
			if property.Parent.ID == p.ID {
				continue search
			}
		}
		// If not present, then crash
		log.Fatalf("Could not parse equipment properties, parent property \"%s\" is missing.", *property.Parent.ID)
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
		latestVersion := equipmentClass.Versions[len(equipmentClass.Versions)-1]
		latestVersionNum, err := strconv.Atoi(latestVersion.Version)
		if err != nil {
			panic(err)
		}

		switch latestVersion.VersionStatus {
		case domain.VersionStateDraft:
			// Do nothing
			equipmentClassVersion = strconv.Itoa(latestVersionNum)
		case domain.VersionStateActive:
			fallthrough
		default:
			// Create a draft version
			_, err = types.SaveEquipmentClassVersionAs(
				ctx,
				client,
				equipmentClassId,
				strconv.Itoa(latestVersionNum),
				strconv.Itoa(latestVersionNum+1),
			)
			equipmentClassVersion = strconv.Itoa(latestVersionNum + 1)
		}
		if err != nil {
			panic(err)
		}
	} else {
		// Else if doesn't exist, then add a new Equipment Class
		uiSortIndex := 1
		processCell := domain.EquipmentElementLevelProcessCell

		extruder := types.GetEquipmentClassPayload(equipmentClassName, types.StringPtr(equipmentImportData.EquipmentClassDescription), &processCell, uiSortIndex)

		var err error
		_, err = types.CreateEquipmentClass(ctx, client, extruder)
		if err != nil {
			panic(err)
		}

		equipmentClassVersion = "1"
	}

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

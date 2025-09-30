package importSheet

import (
	"context"
	"fmt"
	"log"
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
		// If not present, add details to parent
		parentPropertyId := property.Parent.ID
		parentPropertyList := strings.Split(*parentPropertyId, ".")
		property.Parent = &domain.EquipmentClassPropertyRef{
			ID:           parentPropertyId,
			Label:        types.StringPtr(parentPropertyList[len(parentPropertyList)-1]),
			BindingType:  &bound,
			PropertyType: &classType,
			EquipmentClassVersion: &domain.EquipmentClassVersionRef{
				ID:      &equipmentClassName,
				Version: types.StringPtr("1"),
			},
		}
	}

	// Add Equipment Class
	log.Println("\tAdding Equipment Class")

	uiSortIndex := 1
	processCell := domain.EquipmentElementLevelProcessCell

	extruder := types.GetEquipmentClassPayload(equipmentClassName, types.StringPtr(equipmentImportData.EquipmentClassDescription), &processCell, uiSortIndex)

	err := types.CreateEquipmentClass(ctx, client, extruder)
	if err != nil {
		panic(err)
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
				Version: types.StringPtr("1"),
			},
			ValueUnitOfMeasure: property.ValueUnitOfMeasure,
		})
		if err != nil {
			fmt.Printf("Failed to input property with ID \"%s\": %s\n", *property.ID, err)
		}
	}
}

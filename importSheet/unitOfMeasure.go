package importSheet

import (
	"context"
	"log"

	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
)

func UnitOfMeasure(ctx context.Context, client *graphql.Client, equipmentImportData EquipmentImportData) {
	units := make([]domain.AddUnitOfMeasureInput, 0)

out:
	for _, property := range equipmentImportData.EquipmentClassProperties {
		uom := property.UnitOfMeasure.ID
		if uom == "" {
			continue out
		}

		// Check that UoM was not already added
		for _, unit := range units {
			if unit.ID == uom {
				continue out
			}
		}

		log.Printf("\tAdding UoM for %s", uom)

		var dataType domain.DataType

		switch property.UnitOfMeasure.DataType {
		case "Double":
			dataType = domain.DataTypeFloat
		case "UInt16":
			dataType = domain.DataTypeUINt16
		case "UInt32":
			dataType = domain.DataTypeUINt32
		case "Boolean":
			dataType = domain.DataTypeBool
		case "Byte":
			dataType = domain.DataTypeByte
		case "DateTime":
			dataType = domain.DataTypeDateTime
		case "String":
			fallthrough
		case "LocalizedText":
			fallthrough
		default:
			dataType = domain.DataTypeString
		}

		units = append(units, domain.AddUnitOfMeasureInput{
			ID:       uom,
			DataType: &dataType,
		})
	}

	err := types.CreateUnitOfMeasure(ctx, client, units)
	if err != nil {
		log.Printf("could not add unit of measure: %s", err.Error())
	}
}

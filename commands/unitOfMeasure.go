package commands

import (
	"context"
	"log"

	"rhize-data-collection-import/domain"
	"rhize-data-collection-import/types"

	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
)

var (
	UnitOfMeasureCmd = &cobra.Command{
		Use:     "unitOfMeasure",
		Short:   "Import UoMs from file",
		Aliases: []string{"uom", "units"},
		Run:     importUnitOfMeasure,
	}
)

func importUnitOfMeasure(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	log.Println("Starting import for Units of Measure")
	unitOfMeasure(ctx, Client, ImportData.EquipmentClass.Properties)
	log.Println("Finished import for Units of Measure")
}

func unitOfMeasure(ctx context.Context, client *graphql.Client, properties []types.ImportEquipmentClassProperty) {
	units := make([]domain.AddUnitOfMeasureInput, 0)

out:
	for _, property := range properties {
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

		unit := domain.AddUnitOfMeasureInput{
			ID:       uom,
			DataType: &dataType,
		}
		units = append(units, unit)

		existingUoM, err := types.GetUnitOfMeasure(ctx, client, uom)
		if err != nil {
			log.Printf("could not query unit of measure: %s", err.Error())
			continue out
		}
		if existingUoM != nil {
			log.Printf("\tUnit of Measure %s already exists, skipping", uom)
			continue out
		}

		log.Printf("\tAdding UoM for %s", uom)
		err = types.CreateUnitOfMeasure(ctx, client, []domain.AddUnitOfMeasureInput{unit})
		if err != nil {
			log.Printf("could not add unit of measure: %s", err.Error())
		}
	}

}

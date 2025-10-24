package commands

import (
	"context"
	"log"
	"strings"

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

		switch strings.ToLower(property.UnitOfMeasure.DataType) {
		case "number":
			fallthrough
		case "double":
			dataType = domain.DataTypeFloat
		case "uint16":
			dataType = domain.DataTypeUINt16
		case "uint32":
			dataType = domain.DataTypeUINt32
		case "boolean":
			dataType = domain.DataTypeBool
		case "byte":
			dataType = domain.DataTypeByte
		case "datetime":
			dataType = domain.DataTypeDateTime
		case "localizedtext":
			fallthrough
		case "string":
			dataType = domain.DataTypeString
		default:
			// Skip any with an unknown data type
			log.Printf("\tUnknown data type \"%s\", skipping this Unit of Measure\n", property.UnitOfMeasure.DataType)
			continue out
		}

		unit := domain.AddUnitOfMeasureInput{
			ID:       uom,
			DataType: &dataType,
		}
		units = append(units, unit)

		existingUoM, err := types.GetUnitOfMeasure(ctx, client, uom)
		if err != nil {
			log.Printf("\tcould not query unit of measure: %s", err.Error())
			continue out
		}
		if existingUoM != nil {
			log.Printf("\tUnit of Measure %s already exists, skipping", uom)
			continue out
		}

		log.Printf("\tAdding UoM for %s", uom)
		err = types.CreateUnitOfMeasure(ctx, client, []domain.AddUnitOfMeasureInput{unit})
		if err != nil {
			log.Printf("\tcould not add unit of measure: %s", err.Error())
		}
	}

}

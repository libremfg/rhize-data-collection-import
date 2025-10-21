package types

type ImportData struct {
	Datasource     string
	EquipmentClass ImportEquipmentClass
	Equipment      []ImportEquipment
}

type ImportEquipmentClass struct {
	Label       string
	Description string
	Properties  []ImportEquipmentClassProperty
}

type ImportEquipmentClassProperty struct {
	ID            string
	UnitOfMeasure ImportUnitOfMeasure
	Use           bool
}

type ImportUnitOfMeasure struct {
	ID       string
	DataType string
}

type ImportEquipment struct {
	ID          string
	TagBindings []ImportTagBinding
}

type ImportTagBinding struct {
	PropertyID string
	Tag        string
}

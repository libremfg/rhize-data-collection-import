package domain

type DataType string

const (
	DataTypeBool                DataType = "BOOL"
	DataTypeBoolArray           DataType = "BOOL_ARRAY"
	DataTypeByte                DataType = "BYTE"
	DataTypeByteArray           DataType = "BYTE_ARRAY"
	DataTypeByteString          DataType = "BYTE_STRING"
	DataTypeByteStringArray     DataType = "BYTE_STRING_ARRAY"
	DataTypeDataValue           DataType = "DATA_VALUE"
	DataTypeDataValueArray      DataType = "DATA_VALUE_ARRAY"
	DataTypeDateTime            DataType = "DATE_TIME"
	DataTypeDateTimeArray       DataType = "DATE_TIME_ARRAY"
	DataTypeExpandedNodeID      DataType = "EXPANDED_NODE_ID"
	DataTypeExpandedNodeIDArray DataType = "EXPANDED_NODE_ID_ARRAY"
	DataTypeFloat               DataType = "FLOAT"
	DataTypeFloat64             DataType = "FLOAT64"
	DataTypeFloat64Array        DataType = "FLOAT64_ARRAY"
	DataTypeFloatArray          DataType = "FLOAT_ARRAY"
	DataTypeGUID                DataType = "GUID"
	DataTypeGUIDArray           DataType = "GUID_ARRAY"
	DataTypeInt                 DataType = "INT"
	DataTypeInt16               DataType = "INT16"
	DataTypeInt16Array          DataType = "INT16_ARRAY"
	DataTypeInt32               DataType = "INT32"
	DataTypeInt32Array          DataType = "INT32_ARRAY"
	DataTypeInt64               DataType = "INT64"
	DataTypeInt64Array          DataType = "INT64_ARRAY"
	DataTypeLocalizedText       DataType = "LOCALIZED_TEXT"
	DataTypeLocalizedTextArray  DataType = "LOCALIZED_TEXT_ARRAY"
	DataTypeNodeID              DataType = "NODE_ID"
	DataTypeNodeIDArray         DataType = "NODE_ID_ARRAY"
	DataTypeQualifiedName       DataType = "QUALIFIED_NAME"
	DataTypeQualifiedNameArray  DataType = "QUALIFIED_NAME_ARRAY"
	DataTypeStatusCode          DataType = "STATUS_CODE"
	DataTypeStatusCodeArray     DataType = "STATUS_CODE_ARRAY"
	DataTypeString              DataType = "STRING"
	DataTypeStringArray         DataType = "STRING_ARRAY"
	DataTypeStructure           DataType = "STRUCTURE"
	DataTypeStructureArray      DataType = "STRUCTURE_ARRAY"
	DataTypeSByte               DataType = "S_BYTE"
	DataTypeSByteArray          DataType = "S_BYTE_ARRAY"
	DataTypeUINt                DataType = "UINT"
	DataTypeUINt16              DataType = "UINT16"
	DataTypeUINt16Array         DataType = "UINT16_ARRAY"
	DataTypeUINt32              DataType = "UINT32"
	DataTypeUINt32Array         DataType = "UINT32_ARRAY"
	DataTypeUInt64              DataType = "UInt64"
	DataTypeUInt64Array         DataType = "UInt64_ARRAY"
	DataTypeXMLElement          DataType = "XML_ELEMENT"
	DataTypeXMLElementArray     DataType = "XML_ELEMENT_ARRAY"
)

var AllDataType = []DataType{
	DataTypeBool,
	DataTypeBoolArray,
	DataTypeByte,
	DataTypeByteArray,
	DataTypeByteString,
	DataTypeByteStringArray,
	DataTypeDataValue,
	DataTypeDataValueArray,
	DataTypeDateTime,
	DataTypeDateTimeArray,
	DataTypeExpandedNodeID,
	DataTypeExpandedNodeIDArray,
	DataTypeFloat,
	DataTypeFloat64,
	DataTypeFloat64Array,
	DataTypeFloatArray,
	DataTypeGUID,
	DataTypeGUIDArray,
	DataTypeInt,
	DataTypeInt16,
	DataTypeInt16Array,
	DataTypeInt32,
	DataTypeInt32Array,
	DataTypeInt64,
	DataTypeInt64Array,
	DataTypeLocalizedText,
	DataTypeLocalizedTextArray,
	DataTypeNodeID,
	DataTypeNodeIDArray,
	DataTypeQualifiedName,
	DataTypeQualifiedNameArray,
	DataTypeStatusCode,
	DataTypeStatusCodeArray,
	DataTypeString,
	DataTypeStringArray,
	DataTypeStructure,
	DataTypeStructureArray,
	DataTypeSByte,
	DataTypeSByteArray,
	DataTypeUINt,
	DataTypeUINt16,
	DataTypeUINt16Array,
	DataTypeUINt32,
	DataTypeUINt32Array,
	DataTypeUInt64,
	DataTypeUInt64Array,
	DataTypeXMLElement,
	DataTypeXMLElementArray,
}

type PropertyBindingType string

const (
	PropertyBindingTypeBound      PropertyBindingType = "BOUND"
	PropertyBindingTypeCalculated PropertyBindingType = "CALCULATED"
	PropertyBindingTypeStatic     PropertyBindingType = "STATIC"
)

type Isa95PropertyType string

const (
	Isa95PropertyTypeClassType    Isa95PropertyType = "ClassType"
	Isa95PropertyTypeDefaultType  Isa95PropertyType = "DefaultType"
	Isa95PropertyTypeInstanceType Isa95PropertyType = "InstanceType"
)

type EquipmentElementLevel string

const (
	EquipmentElementLevelProcessCell EquipmentElementLevel = "ProcessCell"
)

type VersionState string

const (
	VersionStateActive VersionState = "ACTIVE"
	VersionStateDraft  VersionState = "DRAFT"
)

type StringExactFilterStringFullTextFilterStringRegExpFilter struct {
	Eq *string `json:"eq,omitempty"`
}

type UnitOfMeasure struct {
	ID       string    `json:"id,omitempty"`
	DataType *DataType `json:"dataType,omitempty"`
}

type UnitOfMeasureRef struct {
	ID *string `json:"id,omitempty"`
}

type AddUnitOfMeasureInput struct {
	DataType *DataType `json:"dataType,omitempty"`
	ID       string    `json:"id,omitempty"`
}

type EquipmentClassProperty struct {
	Iid                string                     `json:"iid,omitempty"`
	ID                 string                     `json:"id,omitempty"`
	Label              string                     `json:"label,omitempty"`
	BindingType        PropertyBindingType        `json:"bindingType,omitempty"`
	Parent             *EquipmentClassPropertyRef `json:"parent,omitempty"`
	PropertyType       Isa95PropertyType          `json:"propertyType,omitempty"`
	ValueUnitOfMeasure *UnitOfMeasureRef          `json:"valueUnitOfMeasure,omitempty"`
}

type EquipmentClassPropertyRef struct {
	Iid                   *string                    `json:"iid,omitempty"`
	BindingType           *PropertyBindingType       `json:"bindingType,omitempty"`
	EquipmentClassVersion *EquipmentClassVersionRef  `json:"equipmentClassVersion,omitempty"`
	ID                    *string                    `json:"id,omitempty"`
	Label                 *string                    `json:"label,omitempty"`
	Parent                *EquipmentClassPropertyRef `json:"parent,omitempty"`
	PropertyType          *Isa95PropertyType         `json:"propertyType,omitempty"`
	ValueUnitOfMeasure    *UnitOfMeasureRef          `json:"valueUnitOfMeasure,omitempty"`
}

type AddEquipmentClassPropertyInput struct {
	BindingType           *PropertyBindingType       `json:"bindingType,omitempty"`
	EquipmentClassVersion *EquipmentClassVersionRef  `json:"equipmentClassVersion,omitempty"`
	ID                    string                     `json:"id,omitempty"`
	Label                 string                     `json:"label,omitempty"`
	Parent                *EquipmentClassPropertyRef `json:"parent,omitempty"`
	PropertyType          Isa95PropertyType          `json:"propertyType,omitempty"`
	ValueUnitOfMeasure    *UnitOfMeasureRef          `json:"valueUnitOfMeasure,omitempty"`
}

type EquipmentClass struct {
	Iid      string                   `json:"iid,omitempty"`
	ID       string                   `json:"id,omitempty"`
	Versions []*EquipmentClassVersion `json:"versions,omitempty"`
}

type AddEquipmentClassInput struct {
	ID          string                      `json:"id,omitempty"`
	Label       string                      `json:"label,omitempty"`
	UISortIndex *int                        `json:"uiSortIndex,omitempty"`
	Versions    []*EquipmentClassVersionRef `json:"versions,omitempty"`
}

type EquipmentClassVersion struct {
	Iid           string                    `json:"iid,omitempty"`
	ID            string                    `json:"id,omitempty"`
	Version       string                    `json:"version,omitempty"`
	VersionStatus VersionState              `json:"versionStatus,omitempty"`
	Properties    []*EquipmentClassProperty `json:"properties,omitempty"`
}

type UpdateEquipmentClassInput struct {
	Filter *EquipmentClassFilter `json:"filter,omitempty"`
	Set    *EquipmentClassPatch  `json:"set,omitempty"`
}

type EquipmentClassPatch struct {
	ActiveVersion *EquipmentClassVersionRef `json:"activeVersion,omitempty"`
}

type EquipmentClassFilter struct {
	ID *StringExactFilterStringFullTextFilterStringRegExpFilter `json:"id,omitempty"`
}

type EquipmentClassVersionRef struct {
	Description    *string                      `json:"description,omitempty"`
	DisplayName    *string                      `json:"displayName,omitempty"`
	EquipmentLevel *EquipmentElementLevel       `json:"equipmentLevel,omitempty"`
	ID             *string                      `json:"id,omitempty"`
	Iid            *string                      `json:"iid,omitempty"`
	Properties     []*EquipmentClassPropertyRef `json:"properties,omitempty"`
	Version        *string                      `json:"version,omitempty"`
	VersionStatus  *VersionState                `json:"versionStatus,omitempty"`
}

// Equipment

type AddEquipmentInput struct {
	ID          string                 `json:"id,omitempty"`
	Label       string                 `json:"label,omitempty"`
	UISortIndex *int                   `json:"uiSortIndex,omitempty"`
	Versions    []*EquipmentVersionRef `json:"versions,omitempty"`
}

type EquipmentVersionRef struct {
	DataSources         []*EquipmentDataSourceRef `json:"dataSources,omitempty"`
	Description         *string                   `json:"description,omitempty"`
	DisplayName         *string                   `json:"displayName,omitempty"`
	EquipmentClasses    []*EquipmentClassRef      `json:"equipmentClasses,omitempty"`
	EquipmentLevel      *EquipmentElementLevel    `json:"equipmentLevel,omitempty"`
	ID                  *string                   `json:"id,omitempty"`
	Iid                 *string                   `json:"iid,omitempty"`
	PropertyNameAliases []*PropertyNameAliasRef   `json:"propertyNameAliases,omitempty"`
	TimeZoneName        *string                   `json:"timeZoneName,omitempty"`
	Version             *string                   `json:"version,omitempty"`
	VersionStatus       *VersionState             `json:"versionStatus,omitempty"`
}

type UpdateEquipmentInput struct {
	Filter *EquipmentFilter `json:"filter,omitempty"`
	Set    *EquipmentPatch  `json:"set,omitempty"`
}

type EquipmentFilter struct {
	ID *StringExactFilterStringFullTextFilterStringRegExpFilter `json:"id,omitempty"`
}

type EquipmentPatch struct {
	ActiveVersion *EquipmentVersionRef `json:"activeVersion,omitempty"`
}

type Equipment struct {
	ActiveVersion *EquipmentVersion   `json:"activeVersion,omitempty"`
	ID            string              `json:"id,omitempty"`
	Iid           string              `json:"iid,omitempty"`
	Label         string              `json:"label,omitempty"`
	Versions      []*EquipmentVersion `json:"versions,omitempty"`
}

type EquipmentVersion struct {
	DataSources         []*EquipmentDataSource `json:"dataSources,omitempty"`
	EquipmentClasses    []*EquipmentClass      `json:"equipmentClasses,omitempty"`
	ID                  string                 `json:"id,omitempty"`
	Iid                 string                 `json:"iid,omitempty"`
	PropertyNameAliases []*PropertyNameAlias   `json:"propertyNameAliases,omitempty"`
	Version             string                 `json:"version,omitempty"`
	VersionStatus       VersionState           `json:"versionStatus,omitempty"`
}

type EquipmentClassRef struct {
	ID *string `json:"id,omitempty"`
}

type EquipmentDataSourceRef struct {
	DataSource *DataSourceRef `json:"dataSource,omitempty"`
}

type DataSourceRef struct {
	ID *string `json:"id,omitempty"`
}

type PropertyNameAliasRef struct {
	DataSource           *DataSourceRef `json:"dataSource,omitempty"`
	DataSourceTopicLabel *string        `json:"dataSourceTopicLabel,omitempty"`
	Expression           *string        `json:"expression,omitempty"`
	PropertyLabel        *string        `json:"propertyLabel,omitempty"`
}

type EquipmentVersionFilter struct {
	ID      *StringExactFilterStringFullTextFilterStringRegExpFilter `json:"id,omitempty"`
	Version *StringExactFilterStringFullTextFilter                   `json:"version,omitempty"`
}

type StringExactFilterStringFullTextFilter struct {
	Eq *string `json:"eq,omitempty"`
}

type UpdateEquipmentVersionInput struct {
	Filter *EquipmentVersionFilter `json:"filter,omitempty"`
	Set    *EquipmentVersionPatch  `json:"set,omitempty"`
}

type EquipmentVersionPatch struct {
	DataSources         []*EquipmentDataSourceRef `json:"dataSources,omitempty"`
	PropertyNameAliases []*PropertyNameAliasRef   `json:"propertyNameAliases,omitempty"`
}

type DataSource struct {
	ActiveVersion *DataSourceVersion   `json:"activeVersion,omitempty"`
	ID            string               `json:"id,omitempty"`
	Iid           string               `json:"iid,omitempty"`
	Label         string               `json:"label,omitempty"`
	Versions      []*DataSourceVersion `json:"versions,omitempty"`
}

type DataSourceVersion struct {
	ConnectionString *string            `json:"connectionString,omitempty"`
	ID               string             `json:"id,omitempty"`
	Iid              string             `json:"iid,omitempty"`
	Topics           []*DataSourceTopic `json:"topics,omitempty"`
	Version          string             `json:"version,omitempty"`
	VersionStatus    VersionState       `json:"versionStatus,omitempty"`
}

type DataSourceTopic struct {
	Description *string `json:"description,omitempty"`
	Iid         string  `json:"iid,omitempty"`
	Label       string  `json:"label,omitempty"`
}

type PropertyNameAlias struct {
	Iid           *string `json:"iid,omitempty"`
	PropertyLabel string  `json:"propertyLabel,omitempty"`
}

type EquipmentDataSource struct {
	DataSource *DataSource `json:"dataSource,omitempty"`
}

type PropertyNameAliasFilter struct {
	Iid []string `json:"iid,omitempty"`
}

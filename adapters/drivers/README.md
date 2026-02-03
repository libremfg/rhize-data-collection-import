## Sheet Reader Drivers
Certain assumptions are made when reading from a sheet for determining where certain information is stored.

### General Assumptions
Sheet readers assume that the following data is stored in the following locations:

| Data | Location |
|:----:|:------:|
| Equipment Class Label | A3 |
| Unit of Measure ID | C |
| Unit of Measure Data Type | H |
| Equipment Class Property Data ID | B |
| Equipment Class Property Description | C |
| In Use | N |
| Equipment ID | O*1 |
| Tag Binding Property ID | B* |
| Tag Binding Tag | Q* |
| Tag Binding Expression | R* |

*- Location increments from given column, as detailed below.

#### In Use
In Use is expected to be repeated several times, once for each Equipment, however only the first instance of In Use is typically used.

#### Equipment & Tag Binding
Equipment and Tag Bindings are paired together, and are repeated throughout the sheet in blocks. This pattern repeats every 7 columns, with the first block being 8 long instead of 7.

---

### Excel Assumptions
For Excel files it is assumed that the name of the page in a sheet is the Equipment Class Description.

---

### Considerations
Certain fields are treated or used by the importer in specific ways, as detailed below.

#### Unit of Measure ID
When no ID is present in the Unit of Measure field, the Data Type field is used to create a new Unit of Measure with an ID matching the Data Type.

#### In Use
The importer only reads from the first `In Use` field, despite existing in several possible locations in the sheet.

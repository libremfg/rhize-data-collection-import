# Rhize Data Collection Import Tool

A command-line tool for importing data from Excel/CSV files into Rhize systems via the Libre backend.

## Overview

This tool facilitates the import of equipment data and related information from spreadsheet files (Excel/CSV) into Rhize systems. It supports authentication with the Libre backend and can import data based on equipment class descriptions.

## Features

- Import data from Excel/CSV files
- Authenticate with Libre backend
- Support for equipment class descriptions
- Configurable sheet selection
- Flexible authentication options

#### Usage

Flags for the importer can be seen by using the `--help` flag. Default values show examples values.
```shell
$ ./rhize-import.exe --help
2025/10/23 11:34:01 rhize-data-collection-import v1.6.1
Rhize Data Collection Import

Simple utility to import data from a CSV or XLSX.

Usage:
   [flags]
   [command]

Available Commands:
  completion     Generate the autocompletion script for the specified shell
  equipment      Import equipment
  equipmentClass Import Equipment Class from file
  help           Help about any command
  unitOfMeasure  Import UoMs from file

Flags:
  -A, --apiUrl string         URL for Rhize API (default "http://localhost:8080/graphql")
  -a, --authUrl string        URL for Keycloak authentication (default "http://localhost:8090")
  -b, --bypass                Bypass Keycloak authentication
  -c, --clientId string       Client ID (default "libreBaas")
  -s, --clientSecret string   Client Secret
  -D, --datasource string     The DataSource to bind topics with
  -f, --file string           Excel/CSV file to import data from
  -h, --help                  help for this command
  -p, --password string       Password for user/pass authentication
  -r, --realm string          Keycloak Realm (default "libre")
  -S, --sheet string          The Excel Sheet to search for data in
  -u, --username string       Username for user/pass authentication

Use " [command] --help" for more information about a command.
```

#### Example
Below is a command for importing Units of Measure from an Excel (XLSX) file "data.xlsx", a Sheet titled "Oven_A", and an otherwise default Rhize configuration. Configuration for authentication is read in from a `.env` file.

```shell
$ ./rhize-importer.exe unitOfMeasure \
  --file="./data.xlsx" \
  --sheet="Oven_A" 
2025/10/23 11:56:23 rhize-data-collection-import v1.6.1
2025/10/23 11:56:23 Loading values from environment for unset flags
2025/10/23 11:56:23 Log in Successfully
2025/10/23 11:56:23 Reading data from provided file
2025/10/23 11:56:23 Starting import for Units of Measure
2025/10/23 11:56:24     Adding UoM for °C
2025/10/23 11:56:24     Adding UoM for Amps
2025/10/23 11:56:24     Adding UoM for psi
2025/10/23 11:56:24 Finished import for Units of Measure
```

For import Equipment a DataSource must also be defined. The example assumes a DataSource "OPCUA" is configured in Rhize.

```shell
$ ./rhize-importer.exe equipment \
  --file="./data.xlsx" \
  --sheet="Oven_A" \
  --datasource="OPCUA"
```

### Command Line Arguments

| Flag | Description | Default |
|:------|:-------------|:---------|
| `--apiUrl` | URL for Rhize API | `http://localhost:8080/graphql` |
| `--authUrl` | URL for Keycloak Auth | `http://localhost:8090` |
| `--bypass` | Enable/disable authentication | `false` |
| `--clientId` | Client ID for authentication | `libreBaas` |
| `--clientSecret` | Client secret for authentication | (optional) |
| `--datasource` | The DataSource to bind topics with | (optional, required for Equipment) |
| `--file` | Path to Excel/CSV file to import | (required) |
| `--help` | help for this command | (optional) |
| `--password` | Authentication password | (optional) |
| `--realm` | Keycloak Realm | `libre` |
| `--sheet` | Name of sheet to import data from | (optional, required for Excel) |
| `--username` | Authentication username | (optional) |

## Prerequisites

- Go 1.24 or higher
- Access to a Rhize backend system
- Authentication credentials for the Rhize backend

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Change log
- v1.6.1 J.W.
  - Change Equipment Class import to setup Properties as Default Type instead of Instance Type
- v1.6.0 J.W.
  - Change to allow user/pass authentication if both flags are set
- v1.5.0 J.W.
  - Change to use sub-commands for specifying target instead of `target` flag
  - Change to authenticate with a client rather than client and user
  - Remove `target` flag
- v1.4.2 J.W.
  - Add checks for Equipment Class and Datasource to Equipment model checks
- v1.4.1 J.W.
  - Change to require adding Datasources to Equipment manually
  - Fix duplicate Datasources appearing on Equipment after binding
- v1.4.0 J.W.
  - Add flag `target` for specifying which resource to import data for
  - Fix issue where bindings were not being added
  - Fix issue where topic and binding checks were swapped
- v1.3.0 J.W.
  - Add Equipment binding setups from a given Datasource
  - Add support for loading auth information from environment or .env
- v1.2.1 T.H.
  - Change error messages to pin-point where to check the import spreadsheet
- v1.2.0 T.H.
  - Change behaviour to create new draft version (or update current draft), Otherwise create a completely new Equipment Class
- v1.1.0 T.H.
  - Change properties to Class type from Instance
  - Change properties without a datatype to Static from Bound
- v1.0.0 J.W.
  - Initial Release

## Support

For support, please open an issue on the GitHub repository.

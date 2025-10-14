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
  -apiUrl string
        URL for Rhize API (default "http://localhost:8080/graphql")
  -auth
        Authenticate Client (default true)
  -authUrl string
        URL for Keycloak Auth (default "http://localhost:8090")
  -clientId string
        Client ID (default "libreBaas")
  -clientSecret string
        Client Secret (default "FGY1N5eJQHg3EkOOc5O3IaM4op8o2anT")
  -description string
        Equipment Class description used in import for CSV
  -file string
        Excel/CSV file to import data from (default "./Copy of OPC_UA-CS_NBXT Extrusion Data Information.xlsx")
  -password string
        Authentication Password (default "admin")
  -realm string
        Keycloak Realm (default "libre")
  -sheet string
        Name of sheet to import data from (default "40084-3_Feeder")
  -user string
        Authentication Username (default "admin")
```

#### Example
Assuming for an Excel (XLSX) file "data.xlsx", a Sheet titled "Oven_A", and an otherwise default Rhize configuration.

```shell
$ ./rhize-importer.exe \
      --file="./data.xlsx" \
      --sheet="Oven_A" 
2025/09/30 08:17:06 rhize-data-collection-import v1.0.0
2025/09/30 08:17:06 Log in Successfully
2025/09/30 08:17:07 Adding Imported Unit of Measures
2025/09/30 08:17:07     Adding UoM for °C
2025/09/30 08:17:07     Adding UoM for rpm
2025/09/30 08:17:07     Adding UoM for %
2025/09/30 08:17:07     Adding UoM for Amps
2025/09/30 08:17:07     Adding UoM for psi
2025/09/30 08:17:07 Done Imported Unit of Measures
2025/09/30 08:17:07 Adding Imported Equipment model
2025/09/30 08:17:07     Adding Equipment Class
2025/09/30 08:17:07     Adding Equipment Properties
2025/09/30 08:17:17 Done Imported Equipment model
2025/09/30 08:17:17 Done Imported model
```

### Command Line Arguments

| Flag | Description | Default |
|------|-------------|---------|
| `-apiUrl` | URL for Rhize API | `http://localhost:8080/graphql` |
| `-auth` | Enable/disable authentication | `true` |
| `-authUrl` | URL for Keycloak Auth | `http://localhost:8090` |
| `-clientId` | Client ID for authentication | `libreBaas` |
| `-clientSecret` | Client secret for authentication | `FGY1N5eJQHg3EkOOc5O3IaM4op8o2anT` |
| `-description` | Equipment Class description used in import | (required) |
| `-file` | Path to Excel/CSV file to import | (required) |
| `-password` | Authentication password | `admin` |
| `-realm` | Keycloak Realm | `libre` |
| `-sheet` | Name of sheet to import data from | (optional) |
| `-user` | Authentication username | `admin` |

## Prerequisites

- Go 1.24 or higher
- Access to a Rhzie backend system
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

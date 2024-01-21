# Application for Data Validation and Storage

This Go application demonstrates data validation and storage using Dapr. It's designed to validate data for individuals and organizations and then store it in a Dapr state store.

## Overview

The application defines two data structures: `Individual` and `Organization`. It validates the data based on predefined rules:
- For an `Individual`, the name must be present, and the personal number must be an 11-digit number.
- For an `Organization`, both the name and organization number must be present.

After validation, the data is stored in a state store via Dapr.

## Requirements

- Dapr CLI
- Go programming environment

## Application Structure

- `Individual` and `Organization` structs: Represent data structures for individuals and organizations.
- `IsValid` methods: Validate the data according to the specified rules.
- `saveData` function: Handles the storage of validated data using Dapr.

## Usage

To run this application, ensure Dapr is initialized and running. Then execute:

```shell
dapr run --app-id myapp --dapr-http-port 3500 go run main.go

## Data Validation Logic

The validation logic is implemented within the `IsValid` methods of the `Individual` and `Organization` structs. It ensures the data adheres to the required format before storage.

## Dapr State Management

Dapr is used for state management. The `saveData` function communicates with the Dapr sidecar to store data in the configured state store (e.g., Azure Cosmos DB).

## Future Enhancements

Future versions of this application could include:
- Enhanced validation rules.
- Integration with additional Dapr components.
- Scalability improvements.

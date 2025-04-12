# Golang Flight Management System

## Overview

This project is a Golang-based application for managing airplanes and flights. It follows a Component-Oriented Programming (COP) approach and utilizes a modern N-Layer architecture for better separation of concerns and maintainability.

## Project Structure

The project is organized into several directories, each serving a specific purpose:

- **cmd/app**: Contains the entry point of the application.
- **internal/components**: Houses the core components of the application, including airplanes and flights.
- **internal/core**: Defines core domain entities and interfaces for repositories and services.
- **internal/storage/json**: Implements data storage using JSON files for persistence.
- **internal/utils**: Contains utility functions for data loading and other helper functions.

## Setup Instructions

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd golang-airplane
   ```

2. **Install Dependencies**
   Ensure you have Go installed, then run:

   ```bash
   go mod tidy
   ```

3. **Run the Application**
   To start the application, navigate to the `cmd/app` directory and run:
   ```bash
   go run main.go
   ```

## Usage

Once the application is running, you can interact with the API to manage airplanes and flights. The API endpoints will allow you to perform operations such as adding new airplanes, scheduling flights, and retrieving information.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any enhancements or bug fixes.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

# 🏫 School Management REST API

This repository contains a RESTful API built using Go and the Fiber web framework. The API provides endpoints for managing schools, classes, teachers, and students. This application demonstrates how to create a robust and scalable API with Go, utilizing GORM for database interactions and Fiber for handling HTTP requests.

## 🚀 Features

- **School Management:** Create, read, update, and delete school records.
- **Class Management:** Manage classes associated with schools, including adding, updating, and removing classes.
- **Teacher Management:** Handle teacher records, associating them with specific schools.
- **Student Management:** Manage student records, associating them with specific classes.
- **CRUD Operations:** Full CRUD (Create, Read, Update, Delete) functionality for all entities.
- **Relational Data Handling:** Manage relationships between schools, classes, teachers, and students with cascading updates and deletions.
- **JSON API:** JSON formatted responses for all endpoints, making it easy to integrate with front-end applications.

## 🛠️ Installation

1. **Clone the repository:**

    ```sh
    git clone https://github.com/your-username/your-repository-name.git
    cd your-repository-name
    ```

2. **Install dependencies:**

    ```sh
    go mod tidy
    ```

3. **Set up the database:**

    - Update the database configuration in `config.go`.
    - Ensure that your database server is running.

4. **Run the application:**

    ```sh
    go run main.go
    ```

## 🔧 Usage

Once the application is running, you can access the API endpoints using tools like cURL, Postman, or any HTTP client library.

Example endpoints:

- **GET /schools**: Retrieve all schools.
- **POST /schools**: Create a new school.
- **GET /schools/{id}**: Retrieve a specific school by ID.
- **PUT /schools/{id}**: Update a specific school by ID.
- **DELETE /schools/{id}**: Delete a specific school by ID.
- *(Similar endpoints for classes, teachers, and students)*



## 📖 Documentation

For detailed documentation on the API endpoints, refer to the [API Documentation](https://documenter.getpostman.com/view/29627549/2sA3QqfsNS).

## 🤝 Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request.



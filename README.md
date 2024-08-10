# AMA (Ask Me Anything)

**AMA** is a platform designed to facilitate "Ask Me Anything" sessions, where users can pose questions and receive responses. The project features a robust backend built with Go and a dynamic frontend developed using React.

## Table of Contents

- [Features](#features)
- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
- [Usage](#usage)
  - [Example](#example)
- [Contributing](#contributing)
- [License](#license)
- [Contact](#contact)

## Features

- **User Authentication**: Secure login and registration.
- **Question Posting**: Users can post questions and view responses.
- **Admin Dashboard**: Administrators can moderate questions and manage users.
- **Real-time Updates**: Receive live updates on new questions and answers.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed on your system:

- **Go** (version 1.22.5 or higher)
- **Node.js** (version 18 or higher)
- **npm** or **yarn**
- **PostgreSQL** (or another supported database)

### Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/yourusername/ama.git
   cd ama
   ```

2. **Backend setup (GO)**:
   ```bash
   cd go
   ```

   - Navigate to the go directory:
    ```bash
    cd go
    ```

  - Install dependencies:
    ```bash
    go mod tidy
    ```
  - Set up your database and environment variables:
    ```bash
    cp .env.example .env
    ```
  - Once we are using `sqlc` libraty to generate type-safe Go code from SQL queries. You'll need to run some commands
    ```bash
    sqlc generate -f ./internal/db/postgres/sqlc.yaml
    ``` 
  - Run docker to generate the database locally
    ```bash
    docker-compose up
    ```
  - Run the database migrations using a custom command once we are using `tern`. `Tern` is used to manage database migrations.
    ```bash
    go generate ./...
    ```
  - Run the backend server:
    ```bash
    go run ./cmd/ama/main.go
    ```
3. **Frontend Setup (React)**:

    Navigate to the web directory:
    ```bash
    cd web
    ```
    Install dependencies:
    ```bash
    npm install
    # or
    yarn install
    ```
    Run the development server:
    ```bash
    npm dev
    # or
    yarn dev
    ```

## Usage

Once the installation is complete, you can access the AMA platform via your browser:

- **Backend**: Runs on `http://localhost:8080` (or your configured port).
- **Frontend**: Runs on `http://localhost:5173`.

### Example

1. **Sign Up/Login**: Create an account or log in to the platform.
2. **Post a Question**: Navigate to the "Ask a Question" section and submit your query.
3. **View Responses**: Check the responses to your questions on the main dashboard.

## Contributing

We welcome contributions to the AMA project! To contribute:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes and commit them (`git commit -m 'Add new feature'`).
4. Push to the branch (`git push origin feature-branch`).
5. Open a Pull Request.

Please ensure that your code adheres to the project's coding standards and is well-documented.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Contact

If you have any questions, suggestions, or need support, please reach out to the project maintainer:

- **Email**: dev.pedrogiorgetti@gmail.com
- **GitHub**: [pedrogiorgetti](https://github.com/pedrogiorgetti)

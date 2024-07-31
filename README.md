# EWallet

This project is a REST API built with Golang to manage a simple eWallet system. It provides endpoints to create wallets, perform transactions, retrieve wallet details, and view transaction history.

## Contents

- [Getting Started](#getting-started)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [License](#license)

## Getting Started
To run the application:
```bash
make build && make run
```

### Prerequisites

Make sure you have the following software installed:

- [Go](https://golang.org/doc/install)
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/rtzgod/EWallet.git
    cd EWallet
    ```

2. Create a `.env` file with PostgreSQL environment variables:

    ```env
    POSTGRES_USER=myuser
    POSTGRES_PASSWORD=mypassword
    POSTGRES_DB=mydatabase
    ```
   
3. Configure configs/config.yml file according to your PostgresSQL settings

### Usage
Visit [this site](http://localhost:8080/swagger/index.html) to perform API interaction(before that run the app)

### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

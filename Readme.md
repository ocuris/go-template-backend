# Go Template Backend

## Overview

The Workflow Engine is a SaaS-based workflow automation engine similar to n8n. It supports multiple databases (MongoDB, MySQL, Neo4j) and integrates with external tools such as S3, Slack, Gmail, Outlook, and Calendar. The platform is designed for high availability, scalability, and extensibility.

## Features

- **Multi-Database Support**: MongoDB, MySQL, and Neo4j integration.
- **External Service Nodes**: Connect to S3, Slack, Gmail, Outlook, Calendar, etc.
- **Workflow Execution Engine**: Redis-based job queue for parallel execution.
- **Multi-Tenant Architecture**: SaaS-ready with user isolation.
- **Billing & Subscription**: Integrated Stripe for monetization.
- **CI/CD Pipelines**: Automated deployment and testing.
- **Secure & Scalable**: Designed for cloud-native applications.

## Tech Stack

- **Backend**: Go (Golang)
- **Databases**: PostgreSQL, MongoDB, MySQL, Neo4j
- **Queue Management**: Redis
- **Cloud & Storage**: AWS S3
- **API Communication**: gRPC, REST
- **Authentication**: JWT, OAuth
- **Billing**: Stripe
- **Deployment**: Kubernetes, Docker

## Installation

### Prerequisites

- Go 1.23+
- Docker & Docker Compose
- PostgreSQL, MongoDB, MySQL, Neo4j
- Redis

### Steps

1. **Clone the Repository**

   ```sh
   git clone https://github.com/ocuris/go-template-backend.git
   cd go-template-backend
   ```

2. **Setup Environment Variables**
   Copy the `app.config.sample.yml` to `app.config.local.yml` and update the required values.

   ```sh
   cd configs
   cp app.config.sample.yml app.config.local.yml
   ```

3. **Run Using Docker**

   ```sh
   docker-compose up --build
   ```

4. **Run Locally (Without Docker)**

   ```sh
   go mod tidy
   go run main.go
   ```

## Usage

- Access the API at `http://localhost:8080`
- Use the built-in **Admin Dashboard** to monitor workflows
- Integrate external services via API keys & OAuth

## Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create a new feature branch (`git checkout -b feature-name`)
3. Commit changes (`git commit -m 'Add new feature'`)
4. Push to the branch (`git push origin feature-name`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contact

For questions or support, contact [singhrohankumar7@gmail.com](mailto:singhrohankumar7@gmail.com).

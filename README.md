# Bank Transaction Management System

This project is a robust system for managing bank transactions, complete with API documentation using Swagger and Docker Compose support.

## Table of Contents

- [Getting Started](#getting-started)
  - [Prerequisites](#prerequisites)
  - [Running with Docker Compose](#running-with-docker-compose)
- [Why Use MySQL for Bank Transactions?](#why-use-mysql-for-bank-transactions)
- [Additional Considerations](#additional-considerations)

## Getting Started

### Prerequisites

- [Docker Compose](https://docs.docker.com/compose/) for containerized deployment (recommended).
- If you prefer not to use Docker, you can manually run the application by following these steps:

    1. Install and run a MySQL database.
    2. Create a `.env` file in the `cmd` folder.
    3. Run the application using `go run main.go`.

### Running with Docker Compose

In the root folder, execute the following commands to build and run the application:

```shell
docker-compose build
docker-compose up

Swagger API documentation is available by default at:

http://localhost:4238/swagger/index.html

Why Use MySQL for this source?
MySQL, or any other relational database management system (RDBMS), is a preferred choice for managing bank transactions due to the following key reasons:

Data Integrity: RDBMS like MySQL enforces data integrity through constraints and transactions, ensuring accurate and consistent recording of bank transactions. ACID properties guarantee data integrity even in case of system failures.

Transaction Handling: MySQL supports transactions, ensuring that related database operations occur as a single, atomic unit. This reduces the risk of errors in complex transactions.

Scalability: MySQL can handle large volumes of data and transactions, making it suitable for banks with numerous customers and transactions. It can be scaled horizontally or vertically to accommodate increased demand.

Security: MySQL offers robust security features, including user authentication, encryption, access control, and audit trails, crucial for protecting sensitive financial data.

Reporting and Analytics: MySQL supports complex queries and reporting, enabling banks to analyze transaction data, detect patterns, and generate reports for compliance, fraud detection, and business insights.

Audit Trails: Detailed audit trails in MySQL record every database operation, essential for tracking and investigating suspicious transactions, ensuring transparency, and complying with regulations.

Compliance: Banks must adhere to various regulations and standards. MySQL can be configured to enforce compliance through data retention, access controls, and reporting.

Data Backup and Recovery: MySQL supports efficient data backup and recovery mechanisms, ensuring data can be restored in case of disasters or system failures, essential for business continuity.

Cost-Effective: MySQL is open-source, making it a cost-effective choice for many banks, especially smaller financial institutions.

Community Support: MySQL has a large and active user community, providing abundant resources, documentation, and expertise for database management and troubleshooting.

In summary, MySQL is a popular choice for bank transactions due to its reliability, data integrity, security, scalability, and compliance capabilities. Properly configured and managed, MySQL provides a solid foundation for managing financial data and ensuring the smooth operation of banking systems.

Additional Considerations
Depending on your requirements, you can also consider using Redis or other technologies in conjunction with MySQL to enhance your system's functionality.
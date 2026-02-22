# CMPS3162 Lab #4

## YouTube Demo

https://youtu.be/YXrVipI1kw0

## Database CRUD

| Key               | Value                                          |
| ----------------- | ---------------------------------------------- |
| **Student Name**  | [Andres Hung](https://github.com/andreshungbz) |
| **Student Email** | 2018118240@ub.edu.bz                           |
| **Course**        | CMPS3162 - Advanced Databases                  |
| **Due Date**      | February 24, 2026                              |

## Running the Application

### Docker Compose

```
docker compose up
```

### Manual Method

#### Pre-requisites

- make
- curl
- golang-migrate

#### Database Setup

```
CREATE role hotel_user WITH LOGIN PASSWORD 'hotel_password';
CREATE DATABASE hotel;
ALTER DATABASE hotel OWNER TO hotel_user;
```

#### Application Setup

```
cp .envrc.example .envrc
make db/migrations/up
make run
```

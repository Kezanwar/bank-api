# Bank API

## (Learner project) GoLang / Postgres / Docker

### Initial Setup

Create .env in root using this example

```
DB_USER=example
DB_NAME=example
DB_PASSWORD=example
```

Setup a Postgres Database using Docker

```
docker run --name <DB_NAME> -e POSTGRES_PASSWORD=<DB_PASSWORD> -p :<PORT_DOCKER>:<PORT_LOCAL>  -d postgres
```

Amend Makefile to reflect chosen DB_NAME

### Build / Dev

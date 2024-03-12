# bank-api

Setup a Postgres database using docker

```
docker run --name <DB_NAME> -e POSTGRES_PASSWORD=<DB_PASSWORD> -p :<PORT_DOCKER>:<PORT_LOCAL>  -d postgres
```

Amend Makefile to reflect chosen DB_NAME

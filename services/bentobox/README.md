## BentoBox

Microservice. Polls a *parity* or *go-ethereum* JSON/RPC endpoint.
On new block header it will extract certain specified data.

### Quick Start

If you don't have a PSQL at hand, run one using docker, with the instructions
in the below section "Run and Configure a PostgreSQL using docker"

	make bentobox && ./build/bin/bentobox

### Run and Configure a PostgreSQL using docker

Just use the command (from this repository root)

```
make run-psql
```

Will mount the dir `$HOME/.psql` and have a DB for you.

You need to setup your database if it is the first time using it.
Use the instructions below.

#### Setup the database

You can have your database already, or you can use docker as above.
This script will set you up with your database.

Get in a console able to use `psql` and Run

```
psql -U postgres
```

*TIP*: If you are running PSQL with docker,
to access bash inside the container do

```
docker exec -ti <psql-container> /bin/bash
```

Once inside you create your database with

```
CREATE DATABASE bentobox;
\q
```

Now you need to restore the schema. Do from the comamnd line

```
psql -U postgres bentobox < database.sql
```

You are good to go.

### Command Line Options

| Options | Description | Default |
| --- | --- | --- |
| dbname | Postgres DB name | bentobox |
| dbuser | Postgres DB user name | postgres |
| dbpassword | Postgres DB user password | mysecretpassword |
| eth-host | URL of the ethereum JSON RPC source of data | http://127.0.0.1:8545/ |
| ipfs-host | URL of the ipfs HTTP API | http://127.0.0.1:5001/ |
| last-block-polling-interval | value in seconds for the last block polling | 1 |
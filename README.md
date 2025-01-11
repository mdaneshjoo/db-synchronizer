# How to run

## Setup configs

### Setup debezium json file
> this an example
```json
{
  "name": "promom-connector",
  "config": {
    "connector.class": "io.debezium.connector.postgresql.PostgresConnector",
    "plugin.name": "pgoutput",
    "database.hostname": "192.168.0.101",  
    "database.port": "5432",
    "database.user": "postgres",
    "database.password": "pasword",
    "database.dbname": "dbname",
    "database.server.name": "postgres",
    "table.include.list": "public.table"
  }
}

```
### Setup configuration file
> the file name most be `config.yaml` in the root of the project
```yaml
kafka:
  brokers: localhost:29092
  topic: postgres.public.banner
  group_id: "1"
  schemaregistry_url: http://192.168.0.101:8081
database:
  database_source: "sdsdsdsd"
  database_target: "sdsds"
logging:
  info: true
  debug : true
  
```
> `brokers`, `topic`, `group_id`, `schemaregistry_url` are required


## RUN
> make sure your database is up and runing then start compose file `docker compose up --build`

> To listen the changes run `go run main.go`


> if you want to override logging options you can use logging flags `--info`, `--debug`, `--file <filepath>`


> If you want to produse data manully run `go run producer.go <brocker address> <schema registry address>  <topic>`
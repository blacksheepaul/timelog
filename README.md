# timelog

Including a standard MVC structure, cfg, logger.

The database is SQLite.

# Migrate

for example:

create new migration

```bash
migrate -database "sqlite3://dev.db" create -seq -ext sql --dir model/migrations/ init_xxx_table
```

forward

```bash
migrate -database "sqlite3://dev.db" --path model/migrations/ up
```

# Launch

```bash
# Build binary
make build env=prod
# or
make build-linux env=test
# Build image
docker build -t timelog-app .
# Prod
docker run --rm -e ENV=prod -p 8080:8080 timelog-app
# Test
docker run --rm -e ENV=test -p 18080:8080 timelog-app
```

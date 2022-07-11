# Ley

Manages nodes in a Software Defined Network.

## Development

### Requirements

* [asdf](https://asdf-vm.com/)
* Golang
* Docker
* [pre-commit](http://pre-commit.com/)

### Running tests

```bash
docker compose up db
```

```bash
dagger do test
```

### Running the service

```bash
docker compose up
```

### Creating a database migration

```bash
migrate create \
  -dir internal/manager/migrations \
  -ext '.sql' \
  <name>
```

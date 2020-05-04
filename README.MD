## Requirements

- Install Golang v1.13
- Install [soda](https://gobuffalo.io/en/docs/db/toolbox/) to create migration database.

## Instalation

Install all Libraries from go mod
```
go get .
```

create `.env` file
```
cp .env.example .env
```

## Migration up

Using soda to up our migration files in the `migrations` folder

```
soda migrate -p database up
```

## Running App

`go runmain.go`
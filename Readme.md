Substitute values in configuration files using env variables

Build with command
```shell
go build -o filesubst cmd/filesubst/main.go
```

File format:
.env
```shell
some_val=first
second_val=1
third_val=third
```

Previously set environment variables
```shell
export some_val=value1
export third_val=value2
```

run:
```shell
chmod +x filesubst
filesubst -f .env
```
# Book Shop

## Description
This is a small go application that simulates a book shop. It has a list of books which are related to authors. 
For demo purposes, we store this data in an in-memory sqlite database, and automatically generate it on startup.

## Building
You can build the application by running the following command:
```shell
go build -o bookshop
```

You can also build the docker image by running the following command:
```shell
docker build -t bookshop:latest .
```
# Transportation Back

Transportation back is the backend made to serve a trucks company. It is made mainly 
with golang

## Installation

Install postgresql database. I recommend using docker to run it.

Get docker image
```bash
docker pull postgres
```

Run the container
```bash
docker run --name postgresql -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -v /data:/var/lib/postgresql/data -d postgres
```

Optional: You can make you container restart automaticaly
```bash
docker run -d --restart unless-stopped postgresql 
```

You should create two databases, one for testing called transportationtest and other for development called transportation. Additionally you should create two env file on the root of the project. 
The first file should be called .env, and could be 
```
host=localhost
port=5432
user=postgres
password=postgres
dbname=transportation
```
The second file should be called .env_test
```
host=localhost
port=5432
user=postgres
password=postgres
dbname=transportationtest
```

## Usage

Run tests using the command
```bash
go test ./..
```

Run the project executing
```bash
go run main.go
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
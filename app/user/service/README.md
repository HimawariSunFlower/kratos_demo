## Docker
```bash
# build
make docker

# run
docker run --rm -p 8000:8000 -p 9000:9000 -v </path/to/your/configs>:/data/conf <your-docker-image-name>
```

# config

server:
http:
addr: 0.0.0.0:8000
timeout: 1s
grpc:
addr: 0.0.0.0:9000
timeout: 1s
data:
database:
dsn: "aaa:123@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"
jwt:
secret: "hello"
trace:
endpoint: http://127.0.0.1:14268/api/traces
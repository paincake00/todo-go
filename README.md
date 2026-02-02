### stack

- [gin](https://github.com/gin-gonic/gin)
- [GORM](https://gorm.io/)
- [swag](https://github.com/swaggo/swag), [gin-swagger](https://github.com/swaggo/gin-swagger)
- [godotenv](https://github.com/joho/godotenv)
- [zap](https://go.uber.org/zap)
- [testcontainers](https://golang.testcontainers.org)


swagger:
```
http://localhost:8080/swagger/index.html
```


GET tasks
```sh
 curl  "http://localhost:8080/api/v1/tasks?limit=10&offset=0"
```

POST task
```shell
curl -X POST http://localhost:8080/api/v1/tasks \    
-H "Content-Type: application/json" \
-d '{
  "name": "New Task",
  "description": "some text"
}'
```

PUT task by id
```shell
 curl -X PUT http://localhost:8080/api/v1/tasks/2 \
-H "Content-Type: application/json" \
-d '{
  "name": "New Task",
  "description": "some text",
  "is_completed": false
}'
```
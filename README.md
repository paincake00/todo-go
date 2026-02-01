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
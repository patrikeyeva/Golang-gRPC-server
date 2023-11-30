
## Задание

Требования:
1) Переписать REST сервисы из [данного проекта](https://github.com/patrikeyeva/Golang-Kafka) на gRPC 
2) Переписать логирование на структурное (с исользованием go.uber.org/zap)
3) Добавить трейсы
4) Подключить gRPC-Gateway для обратной совместимости с REST


## Curl запросы

```bash
curl -X POST localhost:9000/article -d '{"name":"Сats","rating":5}' -i

curl -X PUT localhost:9000/article -d '{"id":1,"name":"New Cats","rating":4}' -i

curl -X POST localhost:9000/comment -d '{"article_id":1,"text":"Cute"}' -i
curl -X POST localhost:9000/comment -d '{"article_id":1,"text":"Funny"}' -i

curl -X GET localhost:9000/article/1 -i

curl -X DELETE localhost:9000/article/1 -i
```



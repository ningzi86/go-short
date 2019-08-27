
### 环境变量

APP_REDIS_ADDR

APP_REDIS_PASSWD 

APP_REDIS_DB

### 启动方式

docker run -d -p 80:80 -e APP_REDIS_ADDR=172.16.5.64:6379 -e APP_REDIS_DB=1 --name go-short 52334755/go-short

docker run -d -p 80:80 --name go-short 52334755/go-short
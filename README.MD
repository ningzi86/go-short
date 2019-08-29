
### 环境变量
	
	1.APP_REDIS_ADDR
	说明:redis数据库(可选),容器内置启动redis,也可使用外部redis数据库

	2.APP_REDIS_PASSWD
	说明:redis数据库密码
	
	3.APP_REDIS_DB
	说明:使用的redis数据库

### 启动方式
	
	1.使用外部redis数据库
	docker run -d -p 80:80 -e APP_REDIS_ADDR=172.16.5.64:6379 -e APP_REDIS_DB=1 --name go-short 52334755/go-short

	2.使用内置redis数据库
	docker run -d -p 80:80 --name go-short 52334755/go-short


### API


	Header	
	content-type: application/json
	
1.生成短链接

	method	post
	url		/api/shorten
	body
			{
			    "url": "www.baidu.com",		-- 长地址
			    "expired": 0					-- 失效时间 0永不过期  单位(分钟)
			}
	
	
2.短链接信息

	method 	get
	url		/api/info?shortUrl=[a-zA-Z0-9]+
	

3.短地址跳转长地址
	
	method	get
	url		/[a-zA-Z0-9]+
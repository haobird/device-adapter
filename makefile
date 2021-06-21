# 启动
start:
	docker-compose up -f docker-compose.yml --force-recreate -d 

emqx:
	docker-compose -f emqx.yml up  --force-recreate -d
# 启动
start:
	docker-compose -f docker-compose.yml  up  --force-recreate -d 

emqx:
	docker-compose -f emqx.yml up  --force-recreate -d
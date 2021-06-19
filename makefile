# 启动
start:
	docker-compose up -d --force-recreate

emqx:
	docker-compose -f emqx.yml up  --force-recreate -d
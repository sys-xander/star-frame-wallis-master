.PHONY: gen
gen:
	goctl api go -api .\wallpaper.api --dir . --style goZero --home template


.PHONY: run
run:
	goctl api swagger -api .\wallpaper.api -dir . --filename swagger
	go build -o wallpaper wallpaper.go
	./wallpaper

.PHONY: build
build:
	goctl api swagger -api .\wallpaper.api -dir . --filename swagger
	docker build -t wallpaper:latest .
	docker tag wallpaper:latest registry.cn-hangzhou.aliyuncs.com/don178/wallpaper:latest
	docker push registry.cn-hangzhou.aliyuncs.com/don178/wallpaper:latest
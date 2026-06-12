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

.PHONY: push
push:
	@if not defined VERSION (echo Please provide a version, e.g., make push VERSION=1.0.0 && exit 1)
	git add .
	git commit -m "Release v$(VERSION)"
	git push origin main
	git tag -a v$(VERSION) -m "Release v$(VERSION)"
	git push origin v$(VERSION)
.PHONY:docker
docker:
	@rm goweb || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o goweb .
	@docker rmi -f goweb:v0.0.1
	@docker build -t goweb:v0.0.1 .
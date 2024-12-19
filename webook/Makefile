.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=amd64 go build -o webook .
	@docker rmi -f dakunkun/webook:v0.0.1
	@docker build -t dakunkun/webook:v0.0.1


#windows:docker buildx build -t dakunkun/webook:v0.0.1 .
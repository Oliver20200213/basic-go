.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o webook .
	@docker rmi -f dakunkun/webook:v0.0.1
	@docker build -t dakunkun/webook:v0.0.1

#windows:手动执行的步骤
#cd webook
#$env: GOOS = "linux"
#$env: GOARCH = "amd64"
#go build -o -tags=k8s webook .
#docker rmi -f dakunkun/webook:v0.0.1
#docker buildx build -t dakunkun/webook:v0.0.1 .
#需要删除之前的webook deployment
#kubectl delete deployment webook
#kubectl apply -f .\k8s-webook-deployment.yaml
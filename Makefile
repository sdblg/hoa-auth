build: clean
	DOCKER_SCAN_SUGGEST=false docker build --build-arg env=$(env) -f Dockerfile -t hoa-auth:0.1 .
run:
	docker run -it -env=$(env) $(shell docker images hoa-auth:0.1 -q)
clean:
	docker rm hoa-auth -f; exit 0;
	docker rmi -f $(shell docker images hoa-auth:0.1 -q); exit 0;
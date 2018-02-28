.PHONY: all
.DEFAULT_GOAL := all

all:
	docker build --no-cache -t thrawn01/k8-docker-host-dns:latest .
	docker push thrawn01/k8-docker-host-dns:latest

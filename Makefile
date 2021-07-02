.PHONY: container
container:
	docker build -t concurrency-tester .

.PHONY: run
run: container
	for i in $$(seq 5); do docker run -d --env-file .env concurrency-tester; done

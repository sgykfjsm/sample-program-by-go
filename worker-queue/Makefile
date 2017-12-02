build:
	go build -o queued *.go

run:
	./queued -n 2048

post:
	for i in {1..4096}; do curl http://localhost:8000/work -d name=$USER -d delay=$(expr $i % 11)s; done

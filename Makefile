.PHONY: clean run build docker

clean:
	rm -f dopark

dopark:
	CGO_ENABLED=0 go build -a -ldflags="-w -s" -o dopark ./cmd/

build: dopark

run:
	run.sh

docker:
	docker build --no-cache -t ghcr.io/dereulenspiegel/dopark:latest .


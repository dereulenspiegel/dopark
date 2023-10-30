.PHONY: clean run

clean:
	rm -f dopark

dopark:
	go build -o dopark ./cmd/

run:
	run.sh


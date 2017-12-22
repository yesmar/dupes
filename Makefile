include config.mk

$(cmd): config.mk $(sources)
	go build -ldflags "-s -w -X main.release=$(version)+$(commit)" -o $(cmd)

clean:
	rm -f $(cmd)

count:
	@cloc $(sources)

ensure:
	@dep ensure

lint:
	@gometalinter

.PHONY: clean count ensure lint

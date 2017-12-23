include config.mk

goflags := -ldflags "-s -w -X main.release=$(release)"

$(cmd): config.mk $(sources)
	go build -o $@ $(goflags)

install: config.mk $(sources)
	go install $(goflags)

clean:
	rm -f $(cmd)

count:
	@cloc $(sources)

ensure:
	@dep ensure

lint:
	@gometalinter

.PHONY: clean count ensure lint

.PHONY: build test run tidy release

build:
	$(MAKE) -C go-codebuddy build

test:
	$(MAKE) -C go-codebuddy test

run:
	$(MAKE) -C go-codebuddy run

tidy:
	$(MAKE) -C go-codebuddy tidy

release:
	$(MAKE) -C go-codebuddy release

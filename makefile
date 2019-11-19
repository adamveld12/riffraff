APP := riffraff
PKGS := $(shell go list ./... | grep -v vendor)
GOBIN := $(GOPATH)/bin
LINTBIN := $(GOBIN)/golangci-lint
OUTDIR := .bin
BINARY := $(OUTDIR)/$(app)

GIT_SHA := $$(git rev-parse HEAD)
GIT_BRANCH := $$(git rev-parse --abbrev-ref HEAD)
VERSION := $$(git describe)


define SHORTCUT_DATA
{
	"shortcuts": {
		"*": "https://duckduckgo.com/%s",
		"fb": "https://facebook.com",
		"gh": "https://github.com",
		"gitemoji": "https://www.webfx.com/tools/emoji-cheat-sheet/"
	}
}
endef
export SHORTCUT_DATA

dev: clean $(BINARY)dev
	./$(BINARY)-dev -port 8080 -accesslog=true -data=data.json

test: lint
	go test -v -cover ./...

lint: $(LINTBIN)
	$(LINTBIN) run -p format -p unused -p bugs -p performance 

$(LINTBIN):
	@GO111MODULE=off go get github.com/golangci/golangci-lint/cmd/golangci-lint

package:
	docker build --build-arg VERSION=$(VERSION) \
				 --build-arg COMMIT=$(GIT_SHA) \
 			  	 -t vdhsn/$(APP):$(VERSION) .
	docker tag vdhsn/$(APP):$(VERSION) vdhsn/$(APP):$(GIT_SHA)
	docker tag vdhsn/$(APP):$(VERSION) vdhsn/$(APP):$(GIT_BRANCH)

package-run:
	docker run -it --rm --name riffraff -u 1000:1000 -p 8080:8080 vdhsn/$(APP):$(GIT_SHA)

publish: package
	docker push vdhsn/$(APP):$(VERSION)
	docker push vdhsn/$(APP):$(GIT_BRANCH)

build: $(BINARY)

$(OUTDIR):
	mkdir $@

$(BINARY): packr .bin
	packr build -o $@ -v .

$(BINARY)-dev: .bin
	go build -o $(BINARY)-dev -v .

clean:
	rm -rf .bin

clobber: clean
	rm -rf data.json

data.json:
	@echo "$${SHORTCUT_DATA}" > data.json

packr:
	go get -u github.com/gobuffalo/packr/packr

.PHONY: build clean clobber dev lint package package-run publish packr test

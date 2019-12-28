APP := riffraff
OUTDIR := .bin
TMPDIR := .tmp
PKGS := $(shell go list ./... | grep -v vendor)
GIT_SHA := $(shell git rev-parse HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
VERSION := $(shell git describe)

GOBIN := $(GOPATH)/bin
LINTBIN := $(GOBIN)/golangci-lint
BINARY := $(OUTDIR)/$(app)

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

dev: clean $(BINARY)-dev
	./$(BINARY)-dev -port 8080 -accesslog=true -data=data.json

package-run: package
	@docker run -it --rm --name riffraff -u 1000:1000 -p 8080:8080 -v $$PWD/.tmp:/data vdhsn/$(APP):$(GIT_SHA)

test: lint
	go test -v -cover $(PKGS)

lint: $(LINTBIN)
	@$(LINTBIN) run -p format -p unused -p bugs -p performance 

$(LINTBIN):
	@GO111MODULE=on go get github.com/golangci/golangci-lint/cmd/golangci-lint

package:
	@docker build --build-arg VERSION=$(VERSION) \
				  --build-arg COMMIT=$(GIT_SHA) \
 			  	  -t vdhsn/$(APP):$(VERSION) .
	@docker tag vdhsn/$(APP):$(VERSION) vdhsn/$(APP):$(GIT_SHA)
	@docker tag vdhsn/$(APP):$(VERSION) vdhsn/$(APP):$(GIT_BRANCH)
	docker images | grep 'vdhsn/$(APP)' 

publish: package
	docker push vdhsn/$(APP):$(VERSION)
	docker push vdhsn/$(APP):$(GIT_BRANCH)

release: tag package publish
	@docker tag vdhsn/$(APP):$(GIT_BRANCH) vdhsn/$(APP):$${TAG}
	docker push vdhsn/$(APP):$${TAG}

show-versions:
	@git tag --merged refs/heads/master

tag:
	git tag -as $${TAG}

clean:
	rm -rf $(OUTDIR)

clobber: clean
	rm -rf $(TMPDIR)

build: $(BINARY)

$(OUTDIR):
	@mkdir $@

$(TMPDIR):
	@mkdir $@

$(BINARY): packr .bin
	@packr build -o $@ -v .

$(BINARY)-dev: .bin
	@go build -o $@ -v .

data.json:
	@echo "$${SHORTCUT_DATA}" > data.json

packr:
	go get -u github.com/gobuffalo/packr/packr

.PHONY: build clean clobber dev lint package package-run publish packr show-versions tag test

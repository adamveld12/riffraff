APP := riffraff

dev: clean .bin/$(APP)-dev
	./.bin/$(APP)-dev -port 8080 -accesslog=true -data=data.json

test:
	go test -v ./...

package:
	docker build --build-arg VERSION=$${VERSION:-dev} -t vdhsn/$(APP):$${VERSION:-dev} .

package-run:
	docker run -it --rm --name riffraff -u 1000:1000 -p 8080:8080 vdhsn/$(APP):$${VERSION:-dev}

publish: package
	docker push vdhsn/$(APP):$${VERSION:-dev}

build: .bin/$(APP)

.bin:
	mkdir .bin

.bin/$(APP): packr .bin
	packr build -o .bin/$(APP) -v .

.bin/$(APP)-dev: .bin
	go build -o .bin/$(APP)-dev -v .

clean:
	rm -rf .bin

clobber: clean
	rm -rf data.json

data.json:
	cat <<EOF >> ./data.json
	{
		"shortcuts": {
			"*": "https://duckduckgo.com/%s",
			"fb": "https://facebook.com",
			"gh": "https://github.com",
			"gitemoji": "https://www.webfx.com/tools/emoji-cheat-sheet/"
		}
	}
	EOF

packr:
	go get -u github.com/gobuffalo/packr/packr

.PHONY: build clean dev package package-run publish packr
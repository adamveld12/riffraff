APP := riffraff

dev: clean .bin/$(APP)-dev
	./.bin/$(APP)-dev -port 8080 -accesslog=true

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

packr:
	go get -u github.com/gobuffalo/packr/packr

.PHONY: build clean dev package package-run publish packr
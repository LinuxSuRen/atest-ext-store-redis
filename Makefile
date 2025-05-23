fmt:
	go fmt ./...
build:
	go build -o bin/atest-store-redis .
cp: build
	install bin/atest-store-redis ~/.config/atest/bin
test:
	go test ./... -cover -v -coverprofile=coverage.out
	go tool cover -func=coverage.out
build-image:
	docker build .
hd:
	curl https://linuxsuren.github.io/tools/install.sh|bash
init-env: hd
	hd i cli/cli
	gh extension install linuxsuren/gh-dev
run-e2e:
	cd e2e && ./start.sh
cp-restart:
	systemctl stop atest
	make cp
	systemctl start atest

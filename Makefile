# Usage:
# make                                  # build go source code, build go source code for windows
# make build_reader                     # build go source code for reader
# make build_reader_windows             # build go source code for reader for windows
# make build_broadcaster                # build go source code for broadcaster
# make build_broadcaster_windows        # build go source code for broadcaster for windows
# make build_docker_image               # build docker image with both binaries
# make run_docker                       # run docker environment

all: build_reader build_reader_windows build_broadcaster build_broadcaster_windows
.PHONY: all build_reader build_reader_windows build_broadcaster build_broadcaster_windows build_docker_image run_docker

build_reader: ./cmd/reader/main.go
	@echo ""
	@echo "######################################"
	@echo "##         Building Reader          ##"
	@echo "######################################"
	@echo ""
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static" -s -w' -o dist/reader $<

build_reader_windows: ./cmd/reader/main.go
	@echo ""
	@echo "######################################"
	@echo "##    Building Reader for Windows   ##"
	@echo "######################################"
	@echo ""
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags '-extldflags "-static" -s -w' -o dist/reader.exe $<

build_broadcaster: ./cmd/broadcaster/main.go
	@echo ""
	@echo "######################################"
	@echo "##       Building Broadcaster       ##"
	@echo "######################################"
	@echo ""
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static" -s -w' -o dist/broadcaster $<

build_broadcaster_windows: ./cmd/broadcaster/main.go
	@echo ""
	@echo "######################################"
	@echo "## Building Broadcaster for Windows ##"
	@echo "######################################"
	@echo ""
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -a -ldflags '-extldflags "-static" -s -w' -o dist/broadcaster.exe $<

build_docker_image: ./build/docker/Dockerfile
	@echo ""
	@echo "######################################"
	@echo "##      Building docker image       ##"
	@echo "######################################"
	@echo ""
	@docker build -f $< -t caik/go-stream-broadcaster:latest .

run_docker: ./build/docker/docker-compose.yml
	@echo ""
	@echo "######################################"
	@echo "##    Running docker environment    ##"
	@echo "######################################"
	@echo ""
	@docker-compose -f $< up --build --force-recreate
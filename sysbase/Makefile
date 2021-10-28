export GO111MODULE=on
# export GIN_MODE=release

DISTDIR:=$(shell pwd)/sysbase
T:=$(shell mkdir -p $(DISTDIR)/bin)

all: build

build: mac

mac:
	cd backend && go mod tidy && go build -ldflags "-w -s" -o $(DISTDIR)/bin/sysbase-mac ./

amd:
	cd backend && go mod tidy && xgo --targets=linux/amd64 -v -x -ldflags "-w -s" -out sysbase-amd64 ./ && mv sysbase-amd64-linux-amd64 $(DISTDIR)/bin/sysbase-amd64

arm:
	cd backend && go mod tidy && xgo --targets=linux/arm64 -v -x -ldflags "-w -s" -out sysbase-arm64 ./ && mv sysbase-arm64-linux-arm64 $(DISTDIR)/bin/sysbase-arm64

win:
	cd backend && go mod tidy && xgo --targets=windows/amd64 -v -x -ldflags "-w -s -extldflags=-fno-PIC" -out sysbase.exe ./ && mv sysbase.exe-windows-4.0-amd64.exe $(DISTDIR)/bin/sysbase.exe

dev:
	cd frontend && yarn serve

static:
	cd frontend && yarn build && rm -rf $(DISTDIR)/static && mv dist $(DISTDIR)/static

clean:
	rm -rf $(DISTDIR)

# test:
# 	cd backend && go test ./go/... -race

run:
	/bin/cp -rpf backend/etc $(DISTDIR)/ && /bin/cp -rpf backend/etc/config.js $(DISTDIR)/static && cd $(DISTDIR) && ./bin/sysbase-mac



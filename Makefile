# This how we want to name the binary output
BINARY=legacybest

#VERSION=`cat VERSION`
VERSION=$(shell cat VERSION)
BUILD=`date "+%F-%T"`
#COMMIT=`git rev-parse HEAD`
COMMIT=$(shell git rev-parse HEAD)
GOFILES=$(wildcard src/*.go)
DEBUG=no

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS_f1=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Commit=${COMMIT} -X main.RUNTIMEDEBUG=${DEBUG}"
LDFLAGS_win=-ldflags "-w -s -X main.Version=${VERSION} -X main.Build=${BUILD} -X main.Commit=${COMMIT} -X main.RUNTIMEDEBUG=${DEBUG} -H windowsgui"
BASENEIMAS=${BINARY}-${VERSION}-${COMMIT}


all: build

run:
	go build ${LDFLAGS_f1} -o $(BINARY) -v ./...
	./${BINARY}
deps:
	go get github.com/webview/webview
	go get github.com/google/uuid
	go get github.com/gobuffalo/packr/v2/...
	go get -u github.com/gobuffalo/packr/v2/packr2
windows:
	cd src && \
		packr2 clean && \
		rm -f ../Release/Windows/x86/$(BINARY)-x86.exe; \
		GO111MODULE=on GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ packr2 build ${LDFLAGS_win} -o ../Release/Windows/x86/$(BINARY)-x86.exe && \
		rm -f ../Release/Windows/x64/$(BINARY)-x64.exe; \
		GO111MODULE=on GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ packr2 build ${LDFLAGS_win} -o ../Release/Windows/x64/$(BINARY)-x64.exe && \
		cd -
	cp legacybest.json Release/Windows/x86
	cp legacybest.json Release/Windows/x64
	if [ -f "builds/${BASENEIMAS}-win-x86.zip" ] ; then rm "builds/${BASENEIMAS}-win-x86.zip" ; fi
	if [ -f "builds/${BASENEIMAS}-win-x64.zip" ] ; then rm "builds/${BASENEIMAS}-win-x64.zip" ; fi
	7z a builds/${BASENEIMAS}-win-x86.zip ./Release/Windows/x86/*
	7z a builds/${BASENEIMAS}-win-x64.zip ./Release/Windows/x64/*
release: mac windows linux
	echo "Done making releases!"
mac:
	if [ ! -n "$$DEPLOY_HOST" ]; then cd src && packr2 clean && \
			GO111MODULE=on packr2 build ${LDFLAGS_f1} -o ../Release/Mac/legacybest.app/Contents/MacOS/${BINARY} && cd -; fi
	cp legacybest.json Release/Mac/legacybest.app/Contents/MacOS/
	if [ -f "builds/${BASENEIMAS}-mac-x64.zip" ] ; then rm "builds/${BASENEIMAS}-mac-x64.zip" ; fi
	7z a builds/${BASENEIMAS}-mac-x64.zip Release/Mac/legacybest.app/*
linux:
	cd src && packr2 clean && GO111MODULE=on GOOS=linux GOARCH=amd64 CGO_ENABLED=1 packr2 build ${LDFLAGS_f1} -o ../Release/Linux/$(BINARY)-x64 && cd -
	# building 32bit linux app is not currently supported
	if [ -f "builds/${BASENEIMAS}-linux-x64.zip" ] ; then rm "builds/${BASENEIMAS}-linux-x64.zip" ; fi
	7z a builds/${BASENEIMAS}-linux-x64.zip Release/Linux/*
build:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	$(eval DEBUG := yes)
	cd src && packr2 clean && GO111MODULE=on packr2 build ${LDFLAGS_f1} -o ../legacybest && cd -
clean:
	if [ -f Release/Mac/legacybest.app/Contents/MacOS/app.log ] ; then rm Release/Mac/legacybest.app/Contents/MacOS/app.log ; fi
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf builds/*
	find . -name app.log -exec rm {} \;
.PHONY: clean install release mac windows linux

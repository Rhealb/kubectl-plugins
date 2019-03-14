all: build

TAG?=v3.1
REGISTRY?=ihub.helium.io:29006
FLAGS=
ENVVAR=
GOOS?=linux
ROOTPATH=`pwd` 
BUILDGOPATH=/tmp/k8splugin-build
BUILDPATH=$(BUILDGOPATH)/src/github.com/Rhealb/kubectl-plugins/hostpathpv

 
.IGNORE : buildEnvClean

deps:
	@go get github.com/tools/godep
	
buildEnvClean:
	@rm -rf $(BUILDGOPATH) 1>/dev/null 2>/dev/null || true

buildEnv: buildEnvClean
	@mkdir -p $(BUILDGOPATH)/src/github.com/Rhealb/kubectl-plugins 1>/dev/null 2>/dev/null
	@ln -s $(ROOTPATH) $(BUILDPATH)
	
build: buildEnv clean deps 
	@cd $(BUILDPATH) && GOPATH=$(BUILDGOPATH) $(ENVVAR) GOOS=$(GOOS) CGO_ENABLED=0   godep go build ./...
	@cd $(BUILDPATH) && GOPATH=$(BUILDGOPATH) $(ENVVAR) GOOS=$(GOOS) CGO_ENABLED=0   godep go build -o kubectl-hostpathpv kubectl-hostpathpv.go

install: build
	sudo mv kubectl-hostpathpv /usr/bin
	
uninstall:
	@sudo rm /usr/bin/kubectl-hostpathpv || true
	
moverelease:
	docker build --pull -t ${REGISTRY}/library/hostpathscpmove:${TAG} .
	docker push ${REGISTRY}/library/hostpathscpmove:${TAG}

clean:
	@rm -f kubectl-hostpathpv

format:
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -s -d {} + | tee /dev/stderr)" || \
	test -z "$$(find . -path ./vendor -prune -type f -o -name '*.go' -exec gofmt -s -w {} + | tee /dev/stderr)"
 
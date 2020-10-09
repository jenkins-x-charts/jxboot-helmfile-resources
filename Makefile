CHART_REPO := http://jenkins-x-chartmuseum:8080
NAME := jxboot-helmfile-resources
OS := $(shell uname)

CHARTMUSEUM_CREDS_USR := $(shell cat /tekton/home/basic-auth-user.json)
CHARTMUSEUM_CREDS_PSW := $(shell cat /tekton/home/basic-auth-pass.json)

init:
	helm init --client-only

setup: init
	helm repo add jenkinsxio http://chartmuseum.jenkins-x.io

build: setup build-nosetup

build-nosetup: clean
	helm dependency build jxboot-helmfile-resources
	helm lint jxboot-helmfile-resources

build-no-dep: 
	helm dependency build jxboot-helmfile-resources
	helm lint jxboot-helmfile-resources

install: clean build
	helm upgrade ${NAME} jxboot-helmfile-resources --install

upgrade: clean build
	helm upgrade ${NAME} jxboot-helmfile-resources --install

delete:
	helm delete --purge ${NAME} jxboot-helmfile-resources

clean:
	rm -rf jxboot-helmfile-resources/charts
	rm -rf jxboot-helmfile-resources/${NAME}*.tgz
	rm -rf jxboot-helmfile-resources/requirements.lock

release: clean build
ifeq ($(OS),Darwin)
	sed -i "" -e "s/version:.*/version: $(VERSION)/" jxboot-helmfile-resources/Chart.yaml

else ifeq ($(OS),Linux)
	sed -i -e "s/version:.*/version: $(VERSION)/" jxboot-helmfile-resources/Chart.yaml
else
	exit -1
endif
	helm package jxboot-helmfile-resources
	curl --fail -u $(CHARTMUSEUM_CREDS_USR):$(CHARTMUSEUM_CREDS_PSW) --data-binary "@$(NAME)-$(VERSION).tgz" $(CHART_REPO)/api/charts
	rm -rf ${NAME}*.tgz


test:
	cd tests && go test -v

test-regen:
	cd tests && export HELM_UNIT_REGENERATE_EXPECTED=true && go test -v
VERSION=unknown
CONTAINER=gcr.io/trust-networks/creds-rest-service:${VERSION}
BIN=creds-rest-service

SRCDIR=go/src
PROJSL=${SRCDIR}/project
DEPTOOL=dep ensure -vendor-only -v
SETGOPATH=export GOPATH=$$(pwd)/go

all: godeps build container

build: ${BIN}

${BIN}: *.go
	${SETGOPATH} && cd ${PROJSL} && go build -o ${BIN}

godeps: Gopkg.lock ${PROJSL}
	${SETGOPATH} && cd ${PROJSL} && ${DEPTOOL}

${PROJSL}: ${SRCDIR}
	ln -s ../.. ${PROJSL}

${SRCDIR}:
	mkdir -p ${SRCDIR}

container:
	docker build -t ${CONTAINER} -f Dockerfile .

push:
	gcloud docker -- push ${CONTAINER}

mostlyclean:
	rm -f ${BIN}
	rm -rf ${SRCDIR} # leaves dep cache
	rm -rf vendor

clean: mostlyclean
	rm -rf go # clears dep cache

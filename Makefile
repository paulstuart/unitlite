.PHONY: docker run

MNT = /root/go/src/github.com/paulstuart/unitlite 

docker:
	docker build -t paulstuart/xenial-dqlite:latest .

run:
	docker run \
		-it --rm \
		--security-opt seccomp=unconfined \
		--workdir $(MNT) \
                --mount type=bind,src="$$PWD",dst=$(MNT) \
		paulstuart/xenial-dqlite:latest bash

.phony: ran
ran:
	docker run \
		-it --rm \
		--privileged \
		--security-opt seccomp=unconfined \
		--workdir /meta \
                --mount type=bind,src="$$PWD",dst=/meta \
		paulstuart/xenial-dqlite:latest bash

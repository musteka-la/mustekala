devp2p-node-scrapper:
	./build/modify-geth
	go build -v -o ./build/bin/devp2p-node-scrapper ./services/devp2p-node-scrapper/*.go

block-header-syncer:
	./build/modify-geth
	go build -v -o ./build/bin/block-header-syncer ./services/block-header-syncer/*.go

bentobox:
	go build -v -o ./build/bin/bentobox ./services/bentobox/*.go

custom-psql-image:
	docker build docker/custom-psql/. -t mustekala-psql

run-psql:
	docker run \
		-ti --rm \
		--net=host \
		--name psql \
		-e POSTGRES_PASSWORD=mysecretpassword \
		-v ${PWD}/services/bentobox:/workdir \
		-v ${HOME}/.psql:/var/lib/postgresql/data \
		mustekala-psql
clean:
	rm -rf build/bin/*

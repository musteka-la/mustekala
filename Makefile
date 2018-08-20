devp2p-node-scrapper:
	./build/modify-geth
	go build -v -o ./build/bin/devp2p-node-scrapper ./services/devp2p-node-scrapper/*.go

block-header-syncer:
	./build/modify-geth
	go build -v -o ./build/bin/block-header-syncer ./services/block-header-syncer/*.go

bentobox:
	go build -v -o ./build/bin/bentobox ./services/bentobox/*.go

eth-db-heatmap:
	go build -v -o ./build/bin/eth-db-heatmap ./services/eth-db-heatmap/*.go

find-large-smart-contracts:
	go build -v -o ./build/bin/find-large-smart-contracts ./services/find-large-smart-contracts/*.go

custom-psql-image:
	docker build docker/custom-psql/. -t mustekala-psql

run-psql:
	docker run \
		-ti --rm \
		-p 5432:5432 \
		--name psql \
		-e POSTGRES_PASSWORD=mysecretpassword \
		-v ${PWD}/services/bentobox:/workdir \
		-v ${HOME}/.psql:/var/lib/postgresql/data \
		postgres
clean:
	rm -rf build/bin/*

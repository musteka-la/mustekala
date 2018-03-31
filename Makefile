devp2p-node-scrapper:
	./build/modify-geth
	go build -v -o ./build/bin/devp2p-node-scrapper ./services/devp2p-node-scrapper/*.go

block-header-syncer:
	./build/modify-geth
	go build -v -o ./build/bin/block-header-syncer ./services/block-header-syncer/*.go

clean:
	rm -rf build/bin/*
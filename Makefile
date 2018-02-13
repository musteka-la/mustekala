devp2p-node-scrapper:
	./build/modify-geth
	go build -v -o ./build/bin/devp2p-node-scrapper ./services/devp2p-node-scrapper/*.go

clean:
	rm -rf build/bin/*
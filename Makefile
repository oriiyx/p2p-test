

create:
	@go run ./p2psharing -mode=create -file=./testfile.txt

seed:
	@go run ./p2psharing -mode=seed -file=./testfile.txt -torrent=./testfile.txt.torrent -listen=:50007

download:
	@go run ./p2psharing -mode=download -torrent=./testfile.txt.torrent.encrypted -user=user1 -key=44551708d75ad6e0ad7c22fc277a2d886c8c30ef89b04a0c31bd4c8e9ddb424e -output=./downloads -listen=:50008 -seed=127.0.0.1:50007
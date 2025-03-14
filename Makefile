

create:
	@go run ./p2psharing -mode=create -file=./testfile.txt

seed:
	@go run ./p2psharing -mode=seed -file=./testfile.txt -torrent=./testfile.txt.torrent -listen=:50007

download:
	@go run ./p2psharing -mode=download -torrent=./testfile.txt.torrent.encrypted -user=user1 -key=<key_from_step_1> -output=./downloads -listen=:50008
release:
	go build -o build/sherlock -ldflags="-X 'github.com/KonstantinGasser/sherlock/cmd.Version=$(version)'" main.go
	tar -zcvf sherlock-darwin.tar.gz build/sherlock
	shasum -a 256 sherlock-darwin.tar.gz

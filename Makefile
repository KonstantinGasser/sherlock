release:
	go build -ldflags="-X 'github.com/KonstantinGasser/sherlock/cmd.Version=$(version)'"
	tar -zcvf sherlock-darwin.tar.gz sherlock
compile:
	export GOOS=darwin && go build -o search-mac main.go
	export GOOS=linux && go build -o search-linux main.go
	export GOOS=windows && go build -o search-windows main.go
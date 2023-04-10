all:
	@echo "Building converter ..."
	# CGO_ENABLED=0 go build -o main -ldflags="-X main.DefaultClashUrl=http://other.source/xxx" main.go
	CGO_ENABLED=0 go build -o main main.go

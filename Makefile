
default:
	@echo Testing...
	go test -v
	@echo Building...
	@go build -o ./dist/htc.exe .
	@echo =================================================================================
	./dist/htc

test:
	@go test -v

notest:
	go build -o ./dist/htc.exe .
	./dist/htc

run: 
	./dist/htc

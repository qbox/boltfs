QCHECKSTYLE=go run $(QBOXROOT)/base/qiniu/src/github.com/qiniu/checkstyle/gocheckstyle/gocheckstyle.go -config=.qcodestyle

all:
	cd src; go install -v ./...

rebuild:
	cd src; go install -a -v ./...

install: all
	@echo

test:
	cd src; go test ./...

testv:
	cd src; go test -v ./...

clean:
	go clean -i ./...

style:
	@$(QCHECKSTYLE) src


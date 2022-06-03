# note: call scripts from /scripts
.DEFAULT_GOAL := all

# Flag for building in same OS
GOBUILD=GOOS=linux GOARCH=amd64 go build 

# Flag for build in cross linux distribution
GOSTATICBUILD=CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '-w -extldflags "-static"'

clean:
	rm -f build/package/*

go-update:
	go get -u ./...
	
agi:
	$(GOSTATICBUILD) -o build/package/sangoma-agi cmd/sangoma-agi/sangoma-agi.go

install:
	if [ -d /apps/spoofing-agi ] ; then \
		install -m 755 -d /apps/spoofing-agi; \
	fi
	install -m 755 spoofing-agi /juno/apps/spoofing-agi/
	if [ -d /juno/apps/spoofing-agi/configs ] ; then \
		install -m 755 -d /juno/apps/spoofing-agi/configs; \
	fi
	install -m 644 configs/application.yaml /juno/apps/spoofing-agi/configs
	install -m 644 deployments/spoofing-agi.service /etc/systemd/system/
	systemctl daemon-reload

deploy-agi:
	rsync -avP spoofing-agi vnptspoofing:/juno/apps/spoofing-agi/

all: clean agi
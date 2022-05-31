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
	$(GOSTATICBUILD) -o build/package/spoofing-agi cmd/spoofing-agi/spoofing-agi.go

mnp-checker:
	$(GOSTATICBUILD) -o build/package/mnp-checker cmd/mnp-checker/mnp-checker.go

roaming-checker:
	$(GOSTATICBUILD) -o build/package/roaming-checker cmd/roaming-checker/roaming-checker.go

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

deploy-mnp-checker:
	rsync -avP build/package/mnp-checker vnptspoofing:/juno/apps/mnp-checker/

deploy-roaming-checker:
	rsync -avP build/package/roaming-checker  vnptspoofing:/juno/apps/roaming-checker/

all: clean agi
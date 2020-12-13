# get a list of supported patforms
# go tool dist list
goarch_x86=amd64
goarch_arm=arm
goos=linux

ODIR=/tmp/fangotasia
WORKDIR=${PWD}

SOURCES = $(wildcard *.go) $(wildcard */*.go)

fangotasia.arm: ${SOURCES}
	GOOS=${goos} GOARCH=${goarch_arm} go build -o fangotasia.arm fangotasia.go

fangotasia.x86: ${SOURCES}
	GOOS=${goos} GOARCH=${goarch_x86} go build -o fangotasia.x86 fangotasia.go

release: release.arm release.x86

current:
	echo ${WORKDIR}

release.arm: clean
	mkdir -p $(ODIR)/save
	tar --transform "s/^./fangotasia/" -cvf ${WORKDIR}/fangotasia.arm.tar ./fangotasia.arm ./config
	(cd /tmp && tar -rvf ${WORKDIR}/fangotasia.arm.tar fangotasia )
	gzip ${WORKDIR}/fangotasia.arm.tar

release.x86: clean
	mkdir -p $(ODIR)/save
	tar --transform "s/^./fangotasia/" -cvf ${WORKDIR}/fangotasia.x86.tar ./fangotasia.x86 ./config
	(cd /tmp && tar -rvf ${WORKDIR}/fangotasia.x86.tar fangotasia )
	gzip ${WORKDIR}/fangotasia.x86.tar

.PHONY: clean

clean:
	rm -rf $(ODIR)/*
	rmdir $(ODIR) || echo -n ""
	rm ${WORKDIR}/fangotasia.*.tar.gz || echo -n ""

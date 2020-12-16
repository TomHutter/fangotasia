# get a list of supported patforms
# go tool dist list
goarch_x86=amd64
goarch_arm=arm
goos=linux

SOURCES = $(wildcard *.go) $(wildcard */*.go)

fangotasia.arm: ${SOURCES}
	GOOS=${goos} GOARCH=${goarch_arm} go build -o fangotasia.arm fangotasia.go

fangotasia.x86: ${SOURCES}
	GOOS=${goos} GOARCH=${goarch_x86} go build -o fangotasia.x86 fangotasia.go

release: release.arm release.x86

current:
	echo ${WORKDIR}

release.arm: clean fangotasia.arm
	tar --transform "s/^./fangotasia/" -cvf fangotasia.arm.tar ./fangotasia.arm ./config
	gzip fangotasia.arm.tar

release.x86: clean fangotasia.x86
	tar --transform "s/^./fangotasia/" -cvf fangotasia.x86.tar ./fangotasia.x86 ./config
	gzip fangotasia.x86.tar

.PHONY: clean

clean:
	rm fangotasia.arm || echo -n ""
	rm fangotasia.arm.tar.gz || echo -n ""
	rm fangotasia.x86 || echo -n ""
	rm fangotasia.x86.tar.gz || echo -n ""

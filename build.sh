mkdir -p ./bin

build_all() {
    # Linux
    CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_arm ./            
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_arm64 ./    
    CGO_ENABLED=0 GOOS=linux GOARCH=ppc64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_ppc64 ./    
    CGO_ENABLED=0 GOOS=linux GOARCH=ppc64le go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_ppc64le ./    
    CGO_ENABLED=0 GOOS=linux GOARCH=mips go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_mips ./
    CGO_ENABLED=0 GOOS=linux GOARCH=mipsle go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_mipsle ./
    CGO_ENABLED=0 GOOS=linux GOARCH=mips64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_mips64 ./
    CGO_ENABLED=0 GOOS=linux GOARCH=mips64le go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_mips64le ./
    CGO_ENABLED=0 GOOS=linux GOARCH=s390x go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_linux_s390x ./

    # Darwin
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_darwin_amd64 ./  
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_darwin_arm64 ./

    # Others
    CGO_ENABLED=0 GOOS=freebsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_freebsd_amd64 ./
    CGO_ENABLED=0 GOOS=freebsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_freebsd_386 ./
    CGO_ENABLED=0 GOOS=openbsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_openbsd_amd64 ./
    CGO_ENABLED=0 GOOS=openbsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_openbsd_386 ./
    CGO_ENABLED=0 GOOS=openbsd GOARCH=arm64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_openbsd_arm64 ./
    CGO_ENABLED=0 GOOS=netbsd GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_netbsd_amd64 ./
    CGO_ENABLED=0 GOOS=netbsd GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_netbsd_386 ./
    CGO_ENABLED=0 GOOS=netbsd GOARCH=arm go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_netbsd_arm ./
    CGO_ENABLED=0 GOOS=dragonfly GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_dragonfly_amd64 ./
    CGO_ENABLED=0 GOOS=solaris GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_solaris_amd64 ./
    CGO_ENABLED=0 GOOS=plan9 GOARCH=386 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_plan9_386 ./
    CGO_ENABLED=0 GOOS=plan9 GOARCH=amd64 go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog_plan9_amd64 ./
}

usage() {
    echo "Usage: $0 [all]"
    echo
    echo "Build MailDirWatchDog in different modes:"
    echo "  ./build.sh       Build a single binary for the current system (default)."
    echo "  ./build.sh all   Cross-compile for all supported OS/architectures."
    echo
    echo "Examples:"
    echo "  ./build.sh       -> builds ./bin/MailDirWatchDog"
    echo "  ./build.sh all   -> builds multiple binaries into ./bin/"
    exit 1
}

if [ $# -eq 0 ]; then
    CGO_ENABLED=0  go build -trimpath -buildvcs=false -o ./bin/MailDirWatchDog ./            
elif [ $# -eq 1 ] && [ "$1" = "all" ]; then
    build_all
else
    usage
fi

# Build automation for semver-manager

binary := "smgr"
cmd_pkg := "./src/cmd/smgr"

# Build the CLI binary at the repository root
build:
    go build -o {{binary}} {{cmd_pkg}}

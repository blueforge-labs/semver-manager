# Build automation for semver-manager

binary := "smgr"
cmd_pkg := "./src/cmd/smgr"
prefix := env_var_or_default("PREFIX", "/usr/local")
destdir := env_var_or_default("DESTDIR", "")

# Build the CLI binary at the repository root
build:
    go build -o {{binary}} {{cmd_pkg}}

install: build
    install -d "{{destdir}}{{prefix}}/bin"
    install -m 755 {{binary}} "{{destdir}}{{prefix}}/bin/{{binary}}"

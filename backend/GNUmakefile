.PHONY: install

install:
	go install golang.org/x/tools/gopls@v0.6.5
	go install golang.org/x/tools/cmd/stringer@v0.1.0

sqlboiler:
	sqlboiler psql --pkgname dbmodels --add-global-variants --no-hooks --struct-tag-casing camel --output "./pkg/dbmodels" --wipe --config ./sqlboiler.toml 
	go test ./pkg/dbmodels/*.go -test.config ../../sqlboiler.toml

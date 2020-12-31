.PHONY: middleware

DSN = postgres://workstream:workstream@localhost/workstream?sslmode=disable

SERVICES = accounts
MIGRATIONS := $(SERVICES:%=%-migrate)

all: middleware $(SERVICES)

middleware:
	cd middleware; go build -ldflags "-w" -o ./bin/middleware .; cd ..

$(SERVICES):
	cd service-$@; go build -ldflags "-w" -o ./bin/$@ .; cd ..

migrate: $(MIGRATIONS)

accounts-migrate:
	migrate -source file://./service-accounts/migrations -database ${DSN} up

clean:
	rm -rf logs bin

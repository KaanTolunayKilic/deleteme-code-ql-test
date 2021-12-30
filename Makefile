GO := CGO_ENABLED=0 go

CONTAINER = tweet-extractor:dev-tag

.PHONY: build clean fmtcheck staticcheck

all: clean fmtcheck staticcheck build

gql:
	go run github.com/99designs/gqlgen generate

build: gql
	@echo ">> building"
	$(GO) build -a -o build/tweet-extractor main.go

clean:
	rm -rf build

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

staticcheck:
	@sh -c "'$(CURDIR)/scripts/staticcheck.sh'"

docker:
	@echo ">> building docker image"
	docker build --force-rm -t $(CONTAINER) .

start-compose:
	docker-compose -f deployments/docker-compose.yml up

rm-compose:
	docker-compose -f deployments/docker-compose.yml rm -f -v -s

run: docker start-compose


.PHONY: docs

docs:
	@docker run -v $$PWD/:/docs pandoc/latex -f markdown /docs/README.md -o /docs/build/output/README.pdf

serve:
	docker-compose -f $$PWD/deployments/docker-compose.yml up

stop:
	docker-compose -f $$PWD/deployments/docker-compose.yml stop

test:
	go test -v ./... -coverprofile=cover.out

test-all:
	go test -v ./... -tags=integration -coverprofile=cover.out

ci:
	docker-compose -f $$PWD/deployments/docker-compose.yml down
	docker-compose -f $$PWD/deployments/docker-compose.yml up -d --build form3-client-ci
	docker-compose -f $$PWD/deployments/docker-compose.yml run form3-client-ci ./scripts/ci/ci.sh
	docker-compose -f $$PWD/deployments/docker-compose.yml down

example:
	CGO_ENABLED=0 GO111MODULE=on go build -a -installsuffix cgo -o ./account-api-client ./cmd/account-api-client/main.go
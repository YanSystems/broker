APP_NAME=broker

BOLD=$(shell tput bold)
AQUA=$(shell tput setaf 6)
YELLOW=$(shell tput setaf 3)
RESET=$(shell tput sgr0)

test:
	@echo "$(BOLD)$(AQUA)YAN >$(RESET)$(YELLOW) Running test cases for ${APP_NAME} microservice...$(RESET)"
	@exec go test -v -coverprofile=reports/coverage.out ./...

cov: test
	@echo "$(BOLD)$(AQUA)YAN >$(RESET)$(YELLOW) Creating test coverage report for ${APP_NAME} microservice...$(RESET)"
	@exec go tool cover -html=reports/coverage.out -o reports/coverage.html
	@echo "$(BOLD)$(AQUA)YAN >$(RESET)$(YELLOW) Test coverage report for the ${APP_NAME} microservice has successfully been created.$(RESET)"

image:
	@echo "$(BOLD)$(AQUA)YAN >$(RESET)$(YELLOW) Building docker image for ${APP_NAME} microservice...$(RESET)"
	@exec docker build -t ${APP_NAME} .

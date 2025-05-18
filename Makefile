SHELL			 := cmd.exe
FRONT_END_BINARY := frontApp.exe
BROKER_BINARY    := brokerApp
AUTH_BINARY 	 := authApp
LOGGER_BINARY    := loggerApp
MAIL_BINARY      := mailerApp
LISTENER_BINARY := listenerApp
.PHONY: up up_build down start build_broker build_listener build_auth build_logger build_mailer build_front start stop

up:
	@echo Starting Docker images...
	docker-compose up -d
	@echo Docker images started!

up_build: build_broker build_auth build_logger build_mailer build_listener
	@echo Rebuilding and starting Docker images...
	docker-compose down
	docker-compose up --build -d

down:
	@echo Stopping Docker images…
	docker-compose down

build_broker:
	@echo Building broker binary…
	cd "broker service" && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o $(BROKER_BINARY) ./cmd/api
	@echo Broker built -> $(BROKER_BINARY)

build_listener:
	@echo Building listener binary…
	cd "listener-service" && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o $(LISTENER_BINARY) .
	@echo Listener built -> $(LISTENER_BINARY)

build_auth:
	@echo Building Authentication binary…
	cd authentication-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o $(AUTH_BINARY) ./cmd/api
	@echo Authentication built -> $(AUTH_BINARY)

build_logger:
	@echo Building logger binary…
	cd ./logger-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o $(LOGGER_BINARY) ./cmd/api
	@echo Logger built -> $(LOGGER_BINARY)

build_mailer:
	@echo Building mail binary…
	cd ./mail-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o $(MAIL_BINARY) ./cmd/api
	@echo Logger built -> $(MAIL_BINARY)

build_front:
	@echo Building frontend binary…
	cd frontend && \
	set CGO_ENABLED=0&& set GOOS=windows&&\
	go build -o $(FRONT_END_BINARY) .
	@echo Frontend built -> $(FRONT_END_BINARY)

start: build_front
	@echo Starting front end…
	cd frontend && start /B "" $(FRONT_END_BINARY)

stop:
	@echo Stopping front end…
	taskkill /IM "$(FRONT_END_BINARY)" /F


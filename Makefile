BINARY_NAME=gortisan.exe

build:
	@go mod vendor
	@go build -o tmp/${BINARY_NAME}
	@echo Gortisan built!

run: build 	
	@echo "Staring Gortisan..."
	@.\tmp\${BINARY_NAME} &
	@echo "Gortisan started!"

clean:
	@echo "Cleaning..."
	@go clean
	@del .\tmp\${BINARY_NAME}
	@echo "Cleaned!"

start_docker:
	docker-compose up -d

stop_docker:
	docker-compose down
	
test:
	@echo "Testing..."
	@go test ./...
	@echo "Done!"

start: run
    
stop:
	@echo "Starging the front end..."
	@taskkill /IM ${BINARY_NAME} /F
	@echo "Stopped Gortisan!"

restart: stop start
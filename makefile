build:
	@cd cmd/server && go build -o main 
	@cd cmd/db && go build -o db_init
run:
	@cd cmd/server && go run .
server:
	@cd cmd/server && ./main
initdb:
	@cd cmd/db && ./db_init
liondb:
	@cd cmd/db && go run .

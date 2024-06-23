build:
	go build -o bin/main *.go

run: build
	./bin/main

help:
	@echo  "Swan's Cache Manual\n"
	@echo "   _________                       "
	@echo "  /   _____/_  _  _______    ____  "
	@echo "  \_____  \\ \/ \/ /\__  \  /    \ "
	@echo "  /        \\     /  / __ \|   |  \\"
	@echo " /_______  / \/\_/  (____  /___|  / "
	@echo "         \/              \/     \/  "
	@echo "___________________________________________________"
	@echo "___________________________________________________"
	@echo "Run with makefile			make run "
	@echo "Call with 		 			nc localhost 57"
	@echo "___________________________________________________"

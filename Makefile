dev:
	go build cmd/main.go && mv main tennis_scoreboard
	./tennis_scoreboard

build:
	go build cmd/main.go && mv main tennis_scoreboard
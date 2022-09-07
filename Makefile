build: cards

run: cards
	@./cards

cards: main.go
	@go build

clean:
	@rm -f cards

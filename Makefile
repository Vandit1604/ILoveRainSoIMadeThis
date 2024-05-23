build:
	@go build -o rain

run: build
	@./rain


prepare-lambdas:
	GOOS=linux GOARCH=amd64 go build -o post-messages   ./lambdas/new-service/post-messages/main.go
	zip -r post-messages.zip ./post-messages
	rm ./post-messages

	GOOS=linux GOARCH=amd64 go build -o get-healthcheck ./lambdas/new-service/get-healthcheck/main.go
	zip -r get-healthcheck.zip ./get-healthcheck
	rm ./get-healthcheck

prepare-canary:
	zip -r canary.zip ./nodejs/*
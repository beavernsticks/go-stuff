# Run generators
.PHONY: generate
generate:
	find ./proto -name "*.pb.go" -delete
	protoc --go_out=proto --go_opt=module=github.com/beavernsticks/go-stuff/proto ./proto/definitions/*.proto
	find ./events -name "*.pb.go" -delete
	protoc --go_out=events --go_opt=module=github.com/beavernsticks/go-stuff/events ./events/*.proto

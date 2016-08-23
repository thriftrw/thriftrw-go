package envelope

//go:generate mockgen -source server.go -destination mock_server_test.go -package envelope
//go:generate mockgen -source client.go -destination mock_client_test.go -package envelope

//go:generate mockgen -destination mock_protocol_test.go -package envelope github.com/thriftrw/thriftrw-go/protocol Protocol

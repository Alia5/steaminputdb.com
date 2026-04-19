package steamapi

import "os"

func init() {
	if os.Getenv("GOLANG_PROTOBUF_REGISTRATION_CONFLICT") == "" {
		_ = os.Setenv("GOLANG_PROTOBUF_REGISTRATION_CONFLICT", "ignore")
	}
}

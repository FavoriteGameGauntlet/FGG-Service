package api

//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package auth --target generated/auth --config config.yaml --clean specification/mains/auth.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package games --target generated/games --config config.yaml --clean specification/mains/games.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package points --target generated/points --config config.yaml --clean specification/mains/points.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package timers --target generated/timers --config config.yaml --clean specification/mains/timers.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package users --target generated/users --config config.yaml --clean specification/mains/users.yaml
//go:generate go run github.com/ogen-go/ogen/cmd/ogen@latest --package wheel_effects --target generated/wheel_effects --config config.yaml --clean specification/mains/wheel_effects.yaml

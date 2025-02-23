package common

import "github.com/go-fuego/fuego"

type Controller interface {
	MountRoutes(s *fuego.Server)
}

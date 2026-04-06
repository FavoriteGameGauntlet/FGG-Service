package srvusers

import "FGG-Service/src/users/database"

type Service struct {
	Database dbusers.Database
}

func NewService() *Service {
	s := new(Service)

	return s
}

func (s *Service) ChangeDisplayName(userId int, displayName string) error {
	err := s.Database.ChangeDisplayNameCommand(userId, displayName)

	return err
}

package srvusers

import (
	"FGG-Service/src/common"
	"FGG-Service/src/users/database"
	"database/sql"
	"errors"
)

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

func (s *Service) GetDisplayName(userId int) (displayName *string, err error) {
	displayName, err = s.Database.GetDisplayNameCommand(userId)

	if errors.Is(err, sql.ErrNoRows) || displayName == nil {
		err = common.NewDisplayNameNotFoundError()
		return
	}

	return
}

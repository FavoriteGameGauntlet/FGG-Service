package srvusers

import (
	"FGG-Service/src/common"
	"FGG-Service/src/users/database"
	"FGG-Service/src/users/types"
	"database/sql"
	"errors"
)

type Service struct {
	Database dbusers.IDatabase
}

func NewService() *Service {
	s := new(Service)

	return s
}

func (s *Service) GetAllUserNames() (users typeusers.Users, err error) {
	users, err = s.Database.GetAllUserNamesCommand()

	if errors.Is(err, sql.ErrNoRows) {
		err = nil
	}

	return
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

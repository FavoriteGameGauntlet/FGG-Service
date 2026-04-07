package srvgames_test

import (
	"FGG-Service/src/common"
	"FGG-Service/src/games/service"
	"FGG-Service/src/games/types"
	"FGG-Service/tests/games/mock"
	"errors"
	"testing"
)

// The wishlist game already exists. The error will return.
func TestSrvGames_AddWishlistGame_AlreadyExists(test *testing.T) {
	// Arrange
	userId := 1
	wishlistGame := typegames.WishlistGame{GameId: 1, Name: "Half-Life 1"}

	databaseMock := Initialize_DatabaseMock_AlreadyExists()
	sut := srvgames.Service{Database: databaseMock}

	//Act
	err := sut.AddWishlistGame(userId, wishlistGame)

	// Assert
	if err == nil {
		test.Fatalf("No error found")
	}

	if errors.Is(err, common.NewWishlistGameAlreadyExistsConflictError(wishlistGame.Name)) {
		test.Fatalf("Unexpected error: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_DatabaseMock_AlreadyExists() *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"DoesWishlistGameExistCommand",
			1, "Half-Life 1").
		Return(
			true, nil)

	return databaseMock
}

// The wishlist game doesn't exist. The game exists. The wishlist will be returned successfully.
func TestSrvGames_AddWishlistGame_SuccessReturn(test *testing.T) {
	// Arrange
	userId := 1
	wishlistGame := typegames.WishlistGame{GameId: 1, Name: "Half-Life 1"}

	databaseMock := Initialize_DatabaseMock_SuccessReturn()
	sut := srvgames.Service{Database: databaseMock}

	//Act
	err := sut.AddWishlistGame(userId, wishlistGame)

	// Assert
	if err != nil {
		test.Fatalf("Unexpected error: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func Initialize_DatabaseMock_SuccessReturn() *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"DoesWishlistGameExistCommand",
			1, "Half-Life 1").
		Return(
			false, nil)
	databaseMock.
		On(
			"DoesGameExistCommand",
			"Half-Life 1").
		Return(
			true, nil)
	databaseMock.
		On(
			"GetWishlistGameCommand",
			"Half-Life 1").
		Return(
			typegames.WishlistGame{GameId: 1, Name: "Half-Life 1"}, nil)
	databaseMock.
		On(
			"CreateWishlistGameCommand",
			1, 1).
		Return(
			nil)

	return databaseMock
}

// The wishlist game doesn't exist. The game will be created and returned successfully.
func TestSrvGames_AddWishlistGame_SuccessCreateAndReturn(test *testing.T) {
	// Arrange
	userId := 1
	wishlistGame := typegames.WishlistGame{GameId: 1, Name: "Half-Life 1"}

	databaseMock := InitializeMock_SuccessCreateAndReturn()
	sut := srvgames.Service{Database: databaseMock}

	// Act
	err := sut.AddWishlistGame(userId, wishlistGame)

	// Assert
	if err != nil {
		test.Fatalf("Unexpected error: %v", err)
	}

	databaseMock.AssertExpectations(test)
}

func InitializeMock_SuccessCreateAndReturn() *dbgamesmock.DatabaseMock {
	databaseMock := new(dbgamesmock.DatabaseMock)

	databaseMock.
		On(
			"DoesWishlistGameExistCommand",
			1, "Half-Life 1").
		Return(
			false, nil)
	databaseMock.
		On(
			"DoesGameExistCommand",
			"Half-Life 1").
		Return(
			false, nil)
	databaseMock.
		On(
			"CreateGameCommand",
			"Half-Life 1").
		Return(
			nil)
	databaseMock.
		On(
			"GetWishlistGameCommand",
			"Half-Life 1").
		Return(
			typegames.WishlistGame{GameId: 1, Name: "Half-Life 1"}, nil)
	databaseMock.
		On(
			"CreateWishlistGameCommand",
			1, 1).
		Return(
			nil)

	return databaseMock
}

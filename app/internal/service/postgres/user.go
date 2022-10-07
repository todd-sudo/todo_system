package service_pg

import (
	"context"

	"github.com/mashingan/smapping"
	"github.com/todd-sudo/todo_system/internal/dto"
	"github.com/todd-sudo/todo_system/internal/entity"
	"github.com/todd-sudo/todo_system/internal/hasher"
	"github.com/todd-sudo/todo_system/internal/storage/postgres"
	"github.com/todd-sudo/todo_system/pkg/logging"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	InsertUser(ctx context.Context, user *dto.InsertUserDTO) (*entity.User, error)
	UpdateUser(ctx context.Context, user *dto.UpdateUserDTO) (*entity.User, error)
	DeleteUser(ctx context.Context, username string) error
	FindUserByUsername(ctx context.Context, username string) (*entity.User, error)
	ProfileUser(ctx context.Context, username string) (*entity.User, error)
	ComparePassword(hashedPassword string, password string) error
}

type userService struct {
	ctx     context.Context
	storage *postgres.Storage
	log     logging.Logger
	hasher  hasher.PasswordHasher
}

func NewUserService(ctx context.Context, log logging.Logger, storage *postgres.Storage, hasher hasher.PasswordHasher) UserService {
	return &userService{
		ctx:     ctx,
		log:     log,
		storage: storage,
		hasher:  hasher,
	}
}

// InsertUser - insert user to db
func (s *userService) InsertUser(ctx context.Context, user *dto.InsertUserDTO) (*entity.User, error) {
	userDB := entity.User{}
	err := smapping.FillStruct(&userDB, smapping.MapFields(user))
	if err != nil {
		s.log.Errorf("smapping user-create struct error: %v", err)
		return nil, err
	}
	userDB.Password, err = s.hasher.Hash(userDB.Password)
	if err != nil {
		s.log.Errorf("hashed password error: %v", err)
		return nil, err
	}
	userCreate, err := s.storage.InsertUser(ctx, &userDB)
	if err != nil {
		s.log.Errorf("create user error: %v", err)
		return nil, err
	}
	return userCreate, nil
}

// UpdateUser - update user in db
func (s *userService) UpdateUser(ctx context.Context, user *dto.UpdateUserDTO) (*entity.User, error) {
	userDB := entity.User{}
	err := smapping.FillStruct(&userDB, smapping.MapFields(user))
	if err != nil {
		s.log.Errorf("smapping user-update struct error: %v", err)
		return nil, err
	}
	userUpdate, err := s.storage.UpdateUser(ctx, &userDB)
	if err != nil {
		s.log.Errorf("update user error: %v", err)
		return nil, err
	}
	return userUpdate, nil
}

// DeleteUser - delete user in db
func (s *userService) DeleteUser(ctx context.Context, username string) error {
	err := s.storage.DeleteUser(ctx, username)
	if err != nil {
		s.log.Errorf("delete user error: %v", err)
		return err
	}
	return nil
}

// FindUserByUsername - find user by username in db
func (s *userService) FindUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	userDb, err := s.storage.FindUserByUsername(ctx, username)
	if err != nil {
		s.log.Errorf("find user by username error: %v", err)
		return nil, err
	}
	return userDb, nil
}

func (s *userService) ComparePassword(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		s.log.Errorln(err)
		return err
	}
	return nil
}

// ProfileUser - get profile user in db
func (s *userService) ProfileUser(ctx context.Context, username string) (*entity.User, error) {
	userDb, err := s.storage.ProfileUser(ctx, username)
	if err != nil {
		s.log.Errorf("get profile user error: %v", err)
		return nil, err
	}
	return userDb, nil
}

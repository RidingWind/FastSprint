package service

import (
	"errors"

	"github.com/fastsprint/common/middleware"
	"github.com/fastsprint/common/utils"
	"github.com/fastsprint/user-service/internal/model"
	"github.com/fastsprint/user-service/internal/repository"
)

type UserService interface {
	Register(req model.RegisterRequest) (*model.User, error)
	Login(req model.LoginRequest) (*model.LoginResponse, error)
	GetUserInfo(id uint) (*model.UserInfo, error)
	UpdateUser(id uint, req model.UpdateUserRequest) (*model.User, error)
	ChangePassword(id uint, req model.ChangePasswordRequest) error
	ListUsers(query model.UserQuery) ([]model.User, int64, error)
	GetUserByID(id uint) (*model.User, error)
	GetUsersByIDs(ids []uint) ([]model.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) Register(req model.RegisterRequest) (*model.User, error) {
	existingUser, _ := s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("用户名已存在")
	}

	existingEmail, _ := s.userRepo.GetByEmail(req.Email)
	if existingEmail != nil {
		return nil, errors.New("邮箱已被注册")
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}

	user := &model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		DisplayName:  req.DisplayName,
		Status:       "active",
	}

	if user.DisplayName == "" {
		user.DisplayName = user.Username
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(req model.LoginRequest) (*model.LoginResponse, error) {
	user, err := s.userRepo.GetByUsernameOrEmail(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != "active" {
		return nil, errors.New("账户已被禁用")
	}

	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("用户名或密码错误")
	}

	token, err := middleware.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, errors.New("生成令牌失败")
	}

	return &model.LoginResponse{
		Token: token,
		User: model.UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Email:       user.Email,
			DisplayName: user.DisplayName,
			AvatarURL:   user.AvatarURL,
			Status:      user.Status,
		},
	}, nil
}

func (s *userService) GetUserInfo(id uint) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &model.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName,
		AvatarURL:   user.AvatarURL,
		Status:      user.Status,
	}, nil
}

func (s *userService) UpdateUser(id uint, req model.UpdateUserRequest) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) ChangePassword(id uint, req model.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("原密码错误")
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func (s *userService) ListUsers(query model.UserQuery) ([]model.User, int64, error) {
	return s.userRepo.List(query)
}

func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) GetUsersByIDs(ids []uint) ([]model.User, error) {
	return s.userRepo.GetByIDs(ids)
}

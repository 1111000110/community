package user

import "community/service/user/rpc/client/userservice"

func ModelUserToUser(user *userservice.User) *User {
	return &User{
		UserId:    user.UserId,
		Password:  user.Password,
		Phone:     user.Phone,
		Gender:    user.Gender,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		BirthDate: user.BirthDate,
		Role:      user.Role,
		Status:    user.Status,
		Email:     user.Email,
		Ct:        user.CreateAt,
		Ut:        user.UpdateAt,
	}
}

func UserToModelUser(user *User) *userservice.User {
	return &userservice.User{
		UserId:    user.UserId,
		Password:  user.Password,
		Phone:     user.Phone,
		Gender:    user.Gender,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		BirthDate: user.BirthDate,
		Role:      user.Role,
		Status:    user.Status,
		Email:     user.Email,
		CreateAt:  user.Ct,
		UpdateAt:  user.Ut,
	}
}

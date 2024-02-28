package main

type UserService interface {
	Register(user User) (insertedID string, err error)
}

type MongoUserService struct {
}

func NewMongoUserService() *MongoUserService {
	// DB URL
	// Create connection pool
	return &MongoUserService{}
}

func (m MongoUserService) Register(user User) (insertedID string, err error) {
	//TODO implement me
	panic("implement me")
}

type MockUserService struct {
	registeredUsers []User
	RegisterFunc    func(user User) (insertedID string, err error)
}

func (s *MockUserService) Register(user User) (insertedID string, err error) {
	s.registeredUsers = append(s.registeredUsers, user)
	if s.RegisterFunc != nil {
		return s.RegisterFunc(user)
	}
	return "This is UserID", err
}

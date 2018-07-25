package service

import (
	"resulguldibi/color-api/repository"
	httpClient "resulguldibi/http-client/entity"
	redisClient "resulguldibi/redis-client/entity"
)

type ColorService struct {
	colorRepository repository.ColorRepository
	redisClient     redisClient.IRedisClient
	httpClient      httpClient.IHttpClient
}

type UserService struct {
	redisClient redisClient.IRedisClient
	httpClient  httpClient.IHttpClient
}

type SocketService struct {
}

func NewColorService(colorRepository repository.ColorRepository, redisClient redisClient.IRedisClient) ColorService {
	return ColorService{colorRepository: colorRepository, redisClient: redisClient}
}

func NewColorServiceHttpClient(colorRepository repository.ColorRepository, redisClient redisClient.IRedisClient, httpClient httpClient.IHttpClient) ColorService {
	return ColorService{colorRepository: colorRepository, redisClient: redisClient, httpClient: httpClient}
}

func NewUserService(redisClient redisClient.IRedisClient) UserService {
	return UserService{redisClient: redisClient}
}

func NewUserServiceWithHttpClient(redisClient redisClient.IRedisClient, httpClient httpClient.IHttpClient) UserService {
	return UserService{redisClient: redisClient, httpClient: httpClient}
}

func NewSocketService() SocketService {
	return SocketService{}
}

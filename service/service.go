package service

import (
	"resulguldibi/color-api/repository"
	httpClient "resulguldibi/http-client/entity"
	redisClient "resulguldibi/redis-client/entity"
)

type ColorService struct {
	colorRepository repository.ColorRepository
	redisClient     redisClient.IRedisClient
}

type UserService struct {
	redisClient redisClient.IRedisClient
	httpClient  httpClient.IHttpClient
}

func NewColorService(colorRepository repository.ColorRepository, redisClient redisClient.IRedisClient) ColorService {
	return ColorService{colorRepository: colorRepository, redisClient: redisClient}
}

func NewUserService(redisClient redisClient.IRedisClient) UserService {
	return UserService{redisClient: redisClient}
}

func NewUserServiceWithHttpClient(redisClient redisClient.IRedisClient, httpClient httpClient.IHttpClient) UserService {
	return UserService{redisClient: redisClient, httpClient: httpClient}
}

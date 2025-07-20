package service

import (
	"uptimemonitor/store"
)

type Service struct {
	CheckService
}

func New(store *store.Store) *Service {
	return &Service{
		CheckService: CheckService{
			Store: store,
		},
	}
}

package interfaces

import "github.mpi-internal.com/Yapo/rabbit2kafka/pkg/domain"

//StorageHandler handler for the storage reader
type StorageHandler interface {
	Start(Asyc bool)
	SetReader(reader domain.Reader)
}

//StorageRepo repository that contains the StorageHandler
type StorageRepo struct {
	storageHandler StorageHandler
}

//NewStorageRepo constructor for a StorageRepo
func NewStorageRepo(storageHandler StorageHandler) *StorageRepo {
	storageRepo := new(StorageRepo)
	storageRepo.storageHandler = storageHandler
	return storageRepo
}

//Start starts the underlying StorageHandler
func (s *StorageRepo) Start(async bool) {
	s.storageHandler.Start(async)
}

//SetReader sets the reader function for the StorageHandler
func (s *StorageRepo) SetReader(reader domain.Reader) {
	s.storageHandler.SetReader(reader)
}

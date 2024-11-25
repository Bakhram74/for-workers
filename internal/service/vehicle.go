package service

import (
	"context"

	db "github.com/ShamilKhal/shgo/db/sqlc"
)

type VehicleService struct {
	store db.Store
}

func NewVehicleService(store db.Store) *VehicleService {
	return &VehicleService{
		store: store,
	}
}

func (service *VehicleService) CreateVehicle(ctx context.Context, arg db.CreateVehicleParams) (db.Vehicle, error) {
	return service.store.CreateVehicle(ctx, arg)
}

func (service *VehicleService) GetVehicleByID(ctx context.Context, arg db.GetVehicleByIDParams) (db.Vehicle, error) {
	return service.store.GetVehicleByID(ctx, arg)
}

func (service *VehicleService) DeleteVehicle(ctx context.Context, arg db.DeleteVehicleParams) error {
	return service.store.DeleteVehicle(ctx, arg)
}

func (service *VehicleService) UpdateVehicle(ctx context.Context, arg db.UpdateVehicleParams) (db.Vehicle, error) {
	return service.store.UpdateVehicle(ctx, arg)
}

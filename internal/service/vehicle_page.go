package service

import (
	"context"
	"fmt"

	db "github.com/ShamilKhal/shgo/db/sqlc"
)

func (service *VehicleService) FindVehicle(ctx context.Context, number string, region, limit, offset int) ([]db.Vehicle, int, error) {

	gotNumber := fmt.Sprintf("%%%s%%", number)

	count, err := service.countVehicle(ctx, gotNumber, region)
	if err != nil {
		return []db.Vehicle{}, 0, err
	}

	vehicles, err := service.getVehicle(ctx, gotNumber, region, limit, offset)
	if err != nil {
		return []db.Vehicle{}, 0, err
	}

	return vehicles, int(count), nil
}

func (service *VehicleService) getVehicle(ctx context.Context, number string, region, limit, offset int) ([]db.Vehicle, error) {
	if region <= 0 {
		arg := db.GetVehicleByNumberParams{
			Number: number,
			Limit:  int32(limit),
			Offset: int32(offset),
		}
		return service.store.GetVehicleByNumber(ctx, arg)
	}

	arg := db.GetVehicleByRegionParams{
		Number: number,
		Region: int32(region),
		Limit:  int32(limit),
		Offset: int32(offset),
	}
	return service.store.GetVehicleByRegion(ctx, arg)
}

func (service *VehicleService) countVehicle(ctx context.Context, number string, region int) (int64, error) {
	if region <= 0 {
		return service.store.CountVehicleByNumber(ctx, number)
	}

	arg := db.CountVehicleByRegionParams{
		Number: number,
		Region: int32(region),
	}
	return service.store.CountVehicleByRegion(ctx, arg)
}

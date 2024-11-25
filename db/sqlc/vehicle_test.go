package db

import (
	"context"

	"strings"
	"testing"

	"github.com/ShamilKhal/shgo/pkg/utils"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/require"
)

func createRandomVehicle(t *testing.T, arg CreateVehicleParams, user User) Vehicle {

	vehicle, err := testQueries.CreateVehicle(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, user.ID, vehicle.UserID)
	require.Equal(t, arg.Type, vehicle.Type)
	require.Equal(t, arg.Brand, vehicle.Brand)
	require.Equal(t, arg.Number, vehicle.Number)
	require.Equal(t, arg.Region, vehicle.Region)
	require.Equal(t, arg.Country, vehicle.Country)
	require.NotZero(t, vehicle.CreatedAt)

	return vehicle
}

func TestCreateVehicle(t *testing.T) {
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")

	userArg := CreateUserParams{
		ID:    userID,
		Name:  utils.RandomString(6),
		Phone: utils.RandomNumbers(11),
	}
	user := createRandomUser(t, userArg)

	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	vehicleArg := CreateVehicleParams{
		ID:      id,
		UserID:  user.ID,
		Type:    "легковой",
		Brand:   "toyota",
		Number:  "H777TT",
		Region:  174,
		Country: "RUS",
	}
	vehicle := createRandomVehicle(t, vehicleArg, user)
	arg := GetVehicleByIDParams{
		ID:     vehicle.ID,
		UserID: user.ID,
	}
	got, err := testQueries.GetVehicleByID(context.Background(), arg)
	require.NoError(t, err)
	require.Equal(t, got.UserID, vehicle.UserID)
	require.Equal(t, got.ID, vehicle.ID)
	require.Equal(t, got.Type, vehicle.Type)
	require.Equal(t, got.Brand, vehicle.Brand)
	require.Equal(t, got.Number, vehicle.Number)
	require.Equal(t, got.Region, vehicle.Region)
	require.Equal(t, got.Country, vehicle.Country)
	require.NotZero(t, got.CreatedAt)

}

func TestWithNullableFields(t *testing.T) {
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	userArg := CreateUserParams{
		ID:    userID,
		Name:  utils.RandomString(6),
		Phone: utils.RandomNumbers(11),
	}
	user := createRandomUser(t, userArg)
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	vehicleArg := CreateVehicleParams{
		ID:      id,
		UserID:  user.ID,
		Number:  "H999TT",
		Type:     "грузовой",
		Country: "RUS",
	}
	createRandomVehicle(t, vehicleArg, user)
}

func TestDeleteVehicle(t *testing.T) {
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	userArg := CreateUserParams{
		ID:    userID,
		Name:  utils.RandomString(6),
		Phone: utils.RandomNumbers(11),
	}
	user := createRandomUser(t, userArg)

	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	vehicleArg := CreateVehicleParams{
		ID:      id,
		UserID:  user.ID,
		Type:   "грузовой",
		Brand:  "toyota",
		Number:  "H999TT",
		Region:  174,
		Country: "RUS",
	}
	vehicle := createRandomVehicle(t, vehicleArg, user)

	arg := DeleteVehicleParams{
		ID:     vehicle.ID,
		UserID: user.ID,
	}

	err := testQueries.DeleteVehicle(context.Background(), arg)
	require.NoError(t, err)
}

func TestUpdateVehicleAllFields(t *testing.T) {

	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	userArg := CreateUserParams{
		ID:    userID,
		Name:  utils.RandomString(6),
		Phone: utils.RandomNumbers(11),
	}
	user := createRandomUser(t, userArg)
	id := strings.ReplaceAll(uuid.NewString(), "-", "")
	vehicleArg := CreateVehicleParams{
		ID:      id,
		UserID:  user.ID,
		Type:    "легковой",
		Brand:   "toyota",
		Number:  "H777TT",
		Region:  174,
		Country: "RUS",
	}
	vehicle := createRandomVehicle(t, vehicleArg, user)

	arg := UpdateVehicleParams{
		Type:    pgtype.Text{String: "грузовой", Valid: true},
		Brand:   pgtype.Text{String: "mers", Valid: true},
		Country: pgtype.Text{String: "KAZ", Valid: true},
		ID:      vehicle.ID,
		UserID:  user.ID,
	}
	newVehicle, err := testQueries.UpdateVehicle(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, newVehicle.ID, vehicle.ID)
	require.Equal(t, newVehicle.UserID, user.ID)
	require.Equal(t, newVehicle.Number, "H777TT")
	require.Equal(t, newVehicle.Country, pgtype.Text(pgtype.Text{String: "KAZ", Valid: true}))

	require.NotEqual(t, newVehicle.Type, vehicle.Type)
	require.NotEqual(t, newVehicle.Brand, vehicle.Brand)
	require.NotEqual(t, newVehicle.Country, vehicle.Country)

}

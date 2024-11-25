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

func createRandomUser(t *testing.T, arg CreateUserParams) User {
	user, err := testQueries.CreateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.Equal(t, arg.Name, user.Name)
	require.Equal(t, "guest", user.Role)
	require.Equal(t, arg.Phone, user.Phone)
	require.Equal(t, arg.ImageUrl, user.ImageUrl)
	require.Equal(t, arg.StatusText, user.StatusText)
	require.NotZero(t, user.CreatedAt)
	return user
}

func TestCreateUser(t *testing.T) {
	phone, err := utils.ValidatePhone("+79193273091")
	require.NoError(t, err)
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	arg := CreateUserParams{
		ID:         userID,
		Name:       utils.RandomString(6),
		Phone:      phone,
		StatusText: "Im busy",
		ImageUrl:   "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTAe9NZZk7nUE_anJir2Scf7tsqMHRdEpCbJg&s",
	}

	createRandomUser(t, arg)
	t.Cleanup(func() {
		_, err = testQueries.db.Exec(context.Background(), `TRUNCATE TABLE users CASCADE`)
		require.NoError(t, err)
	})
}

func TestGetUser(t *testing.T) {

	arg := CreateUserParams{
		Name:  "Alex",
		Phone: "79193273091",
	}

	user := createRandomUser(t, arg)

	gotUser, err := testQueries.GetUserByPhone(context.Background(), user.Phone)
	require.NoError(t, err)

	require.Equal(t, gotUser.Name, user.Name)
	require.Equal(t, gotUser.Role, user.Role)
	require.Equal(t, gotUser.Phone, user.Phone)
	require.Equal(t, gotUser.CreatedAt, user.CreatedAt)
	require.NotZero(t, gotUser.CreatedAt)

	//get empty user
	emptyUser, err := testQueries.GetUserByPhone(context.Background(), "79993273099")
	require.Empty(t, emptyUser)
	require.Equal(t, emptyUser.ID, "")
	require.Error(t, err)
}

func TestUpdateUserAllFields(t *testing.T) {
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	arg := CreateUserParams{
		ID:         userID,
		Name:       utils.RandomString(6),
		Phone:      utils.RandomNumbers(11),
		StatusText: "Im busy",
		ImageUrl:   utils.RandomString(9),
	}
	oldUser := createRandomUser(t, arg)

	newName := "NewName"
	newPhone := utils.RandomNumbers(11)
	newImageUrl := utils.RandomString(9)
	newStatus := "New status"

	arg2 := UpdateUserParams{
		ID:         oldUser.ID,
		Role:       pgtype.Text{Valid: true, String: "user"},
		Name:       pgtype.Text{Valid: true, String: newName},
		Phone:      pgtype.Text{Valid: true, String: newPhone},
		ImageUrl:   pgtype.Text{Valid: true, String: newImageUrl},
		StatusText: pgtype.Text{Valid: true, String: newStatus},
	}

	updatedUser, err := testQueries.UpdateUser(context.Background(), arg2)
	require.NoError(t, err)
	require.Equal(t, oldUser.ID, updatedUser.ID)

	require.NotEqual(t, oldUser.Name, updatedUser.Name)
	require.Equal(t, newName, updatedUser.Name)
	require.NotEqual(t, oldUser.Phone, updatedUser.Phone)
	require.Equal(t, newPhone, updatedUser.Phone)

	require.NotEqual(t, oldUser.ImageUrl, updatedUser.ImageUrl)
	require.Equal(t, newImageUrl, updatedUser.ImageUrl)
	require.NotEqual(t, oldUser.StatusText, updatedUser.StatusText)
	require.Equal(t, newStatus, updatedUser.StatusText)

	require.NotEqual(t, oldUser.Role, updatedUser.Role)
	require.Equal(t, "user", updatedUser.Role)
}

func TestUpdateUserOnlyImg(t *testing.T) {
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	arg := CreateUserParams{
		ID:         userID,
		Name:       utils.RandomString(6),
		Phone:      utils.RandomNumbers(11),
		StatusText: "Im busy",
		ImageUrl:   "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTAe9NZZk7nUE_anJir2Scf7tsqMHRdEpCbJg&s",
	}
	oldUser := createRandomUser(t, arg)

	newImg := "https://png.pngtree.com/png-vector/20230918/ourmid/pngtree-man-in-shirt-smiles-and-gives-thumbs-up-to-show-approval-png-image_10094381.png"
	arg2 := UpdateUserParams{
		ID:       oldUser.ID,
		ImageUrl: pgtype.Text{Valid: true, String: newImg},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg2)
	require.NoError(t, err)
	require.Equal(t, oldUser.ID, updatedUser.ID)
	require.NotEqual(t, oldUser.ImageUrl, updatedUser.ImageUrl)

	require.Equal(t, oldUser.Phone, updatedUser.Phone)
	require.Equal(t, updatedUser.Role, "guest")
	require.Equal(t, oldUser.Name, updatedUser.Name)
	require.Equal(t, oldUser.StatusText, updatedUser.StatusText)
}

func TestUpdateUserOnlyPhone(t *testing.T) {
	phone, err := utils.ValidatePhone("+79193273091")
	require.NoError(t, err)
	userID := strings.ReplaceAll(uuid.NewString(), "-", "")
	arg := CreateUserParams{
		ID:         userID,
		Name:       utils.RandomString(6),
		Phone:      phone,
		StatusText: "Im busy",
		ImageUrl:   utils.RandomString(9),
	}
	oldUser := createRandomUser(t, arg)

	newPhone := "79193273085"
	arg2 := UpdateUserParams{
		ID:    oldUser.ID,
		Phone: pgtype.Text{Valid: true, String: newPhone},
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg2)
	require.NoError(t, err)
	require.Equal(t, oldUser.ID, updatedUser.ID)
	require.NotEqual(t, oldUser.Phone, updatedUser.Phone)
	require.Equal(t, newPhone, updatedUser.Phone)
	require.Equal(t, updatedUser.Role, "guest")

	require.Equal(t, oldUser.Name, updatedUser.Name)
	require.Equal(t, oldUser.ImageUrl, updatedUser.ImageUrl)
	require.Equal(t, oldUser.StatusText, updatedUser.StatusText)
}

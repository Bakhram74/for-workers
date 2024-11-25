// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0
// source: vehicle.sql

package db

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const countVehicleByNumber = `-- name: CountVehicleByNumber :one
SELECT COUNT(id) FROM vehicle
 WHERE number LIKE $1
`

func (q *Queries) CountVehicleByNumber(ctx context.Context, number string) (int64, error) {
	row := q.db.QueryRow(ctx, countVehicleByNumber, number)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countVehicleByRegion = `-- name: CountVehicleByRegion :one
SELECT COUNT(id) FROM vehicle
WHERE number LIKE $1 AND region = $2
`

type CountVehicleByRegionParams struct {
	Number string `json:"number"`
	Region int32  `json:"region"`
}

func (q *Queries) CountVehicleByRegion(ctx context.Context, arg CountVehicleByRegionParams) (int64, error) {
	row := q.db.QueryRow(ctx, countVehicleByRegion, arg.Number, arg.Region)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createVehicle = `-- name: CreateVehicle :one
INSERT INTO vehicle (
  id,
  user_id,
  type,
  brand,
  number,
  region,
  country
) VALUES (
 $1, $2, $3, $4, $5, $6,$7
) RETURNING id, user_id, type, brand, number, region, country, created_at
`

type CreateVehicleParams struct {
	ID      string `json:"id"`
	UserID  string `json:"user_id"`
	Type    string `json:"type"`
	Brand   string `json:"brand"`
	Number  string `json:"number"`
	Region  int32  `json:"region"`
	Country string `json:"country"`
}

func (q *Queries) CreateVehicle(ctx context.Context, arg CreateVehicleParams) (Vehicle, error) {
	row := q.db.QueryRow(ctx, createVehicle,
		arg.ID,
		arg.UserID,
		arg.Type,
		arg.Brand,
		arg.Number,
		arg.Region,
		arg.Country,
	)
	var i Vehicle
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Brand,
		&i.Number,
		&i.Region,
		&i.Country,
		&i.CreatedAt,
	)
	return i, err
}

const deleteVehicle = `-- name: DeleteVehicle :exec
DELETE FROM vehicle
WHERE id = $1 AND user_id = $2
`

type DeleteVehicleParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) DeleteVehicle(ctx context.Context, arg DeleteVehicleParams) error {
	_, err := q.db.Exec(ctx, deleteVehicle, arg.ID, arg.UserID)
	return err
}

const getVehicleByID = `-- name: GetVehicleByID :one
SELECT id, user_id, type, brand, number, region, country, created_at FROM vehicle
WHERE id = $1 AND user_id = $2
`

type GetVehicleByIDParams struct {
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

func (q *Queries) GetVehicleByID(ctx context.Context, arg GetVehicleByIDParams) (Vehicle, error) {
	row := q.db.QueryRow(ctx, getVehicleByID, arg.ID, arg.UserID)
	var i Vehicle
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Brand,
		&i.Number,
		&i.Region,
		&i.Country,
		&i.CreatedAt,
	)
	return i, err
}

const getVehicleByNumber = `-- name: GetVehicleByNumber :many
SELECT id, user_id, type, brand, number, region, country, created_at FROM vehicle
WHERE number LIKE $1 LIMIT $2 OFFSET $3
`

type GetVehicleByNumberParams struct {
	Number string `json:"number"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetVehicleByNumber(ctx context.Context, arg GetVehicleByNumberParams) ([]Vehicle, error) {
	rows, err := q.db.Query(ctx, getVehicleByNumber, arg.Number, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Vehicle{}
	for rows.Next() {
		var i Vehicle
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Brand,
			&i.Number,
			&i.Region,
			&i.Country,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVehicleByRegion = `-- name: GetVehicleByRegion :many
SELECT id, user_id, type, brand, number, region, country, created_at FROM vehicle
WHERE number LIKE $1 AND region = $2
LIMIT $3 OFFSET $4
`

type GetVehicleByRegionParams struct {
	Number string `json:"number"`
	Region int32  `json:"region"`
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
}

func (q *Queries) GetVehicleByRegion(ctx context.Context, arg GetVehicleByRegionParams) ([]Vehicle, error) {
	rows, err := q.db.Query(ctx, getVehicleByRegion,
		arg.Number,
		arg.Region,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Vehicle{}
	for rows.Next() {
		var i Vehicle
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Type,
			&i.Brand,
			&i.Number,
			&i.Region,
			&i.Country,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateVehicle = `-- name: UpdateVehicle :one
UPDATE vehicle
SET
  type = COALESCE($1, type),
  brand = COALESCE($2, brand),
  number = COALESCE($3, number),
  region = COALESCE($4, region),
  country = COALESCE($5, country)
WHERE id = $6 AND user_id = $7
RETURNING id, user_id, type, brand, number, region, country, created_at
`

type UpdateVehicleParams struct {
	Type    pgtype.Text `json:"type"`
	Brand   pgtype.Text `json:"brand"`
	Number  pgtype.Text `json:"number"`
	Region  pgtype.Int4 `json:"region"`
	Country pgtype.Text `json:"country"`
	ID      string      `json:"id"`
	UserID  string      `json:"user_id"`
}

func (q *Queries) UpdateVehicle(ctx context.Context, arg UpdateVehicleParams) (Vehicle, error) {
	row := q.db.QueryRow(ctx, updateVehicle,
		arg.Type,
		arg.Brand,
		arg.Number,
		arg.Region,
		arg.Country,
		arg.ID,
		arg.UserID,
	)
	var i Vehicle
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Type,
		&i.Brand,
		&i.Number,
		&i.Region,
		&i.Country,
		&i.CreatedAt,
	)
	return i, err
}
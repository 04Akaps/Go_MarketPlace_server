// Code generated by sqlc. DO NOT EDIT.
// source: launchpad.sql

package goServer

import (
	"context"
	"database/sql"
	"encoding/json"
)

const createNewLaunchpad = `-- name: CreateNewLaunchpad :execresult
INSERT INTO launchpad (
    hash_value,
    first_owner_email,
    ca_address,
    chain_id,
    price,
    airdrop_address,
    whitelist_address
) VALUES (
      ?,?,?,?,?,?,?
)
`

type CreateNewLaunchpadParams struct {
	HashValue        string          `json:"hash_value"`
	FirstOwnerEmail  string          `json:"first_owner_email"`
	CaAddress        string          `json:"ca_address"`
	ChainID          string          `json:"chain_id"`
	Price            int32           `json:"price"`
	AirdropAddress   json.RawMessage `json:"airdrop_address"`
	WhitelistAddress json.RawMessage `json:"whitelist_address"`
}

func (q *Queries) CreateNewLaunchpad(ctx context.Context, arg CreateNewLaunchpadParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, createNewLaunchpad,
		arg.HashValue,
		arg.FirstOwnerEmail,
		arg.CaAddress,
		arg.ChainID,
		arg.Price,
		arg.AirdropAddress,
		arg.WhitelistAddress,
	)
}

const getLaunchpadByChainId = `-- name: GetLaunchpadByChainId :many
SELECT  id, hash_value, first_owner_email, ca_address, chain_id, price, airdrop_address, whitelist_address, created_at FROM launchpad
WHERE  chain_id = ? ORDER BY created_at ASC
`

func (q *Queries) GetLaunchpadByChainId(ctx context.Context, chainID string) ([]Launchpad, error) {
	rows, err := q.db.QueryContext(ctx, getLaunchpadByChainId, chainID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Launchpad{}
	for rows.Next() {
		var i Launchpad
		if err := rows.Scan(
			&i.ID,
			&i.HashValue,
			&i.FirstOwnerEmail,
			&i.CaAddress,
			&i.ChainID,
			&i.Price,
			&i.AirdropAddress,
			&i.WhitelistAddress,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLaunchpadByHash = `-- name: GetLaunchpadByHash :one
SELECT id, hash_value, first_owner_email, ca_address, chain_id, price, airdrop_address, whitelist_address, created_at FROM launchpad
WHERE hash_value = ?
`

func (q *Queries) GetLaunchpadByHash(ctx context.Context, hashValue string) (Launchpad, error) {
	row := q.db.QueryRowContext(ctx, getLaunchpadByHash, hashValue)
	var i Launchpad
	err := row.Scan(
		&i.ID,
		&i.HashValue,
		&i.FirstOwnerEmail,
		&i.CaAddress,
		&i.ChainID,
		&i.Price,
		&i.AirdropAddress,
		&i.WhitelistAddress,
		&i.CreatedAt,
	)
	return i, err
}
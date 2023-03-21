-- name: CreateNewLaunchpad :execresult
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
);

-- name: GetLaunchpadByHash :one
SELECT * FROM launchpad
WHERE hash_value = ?;

-- name: GetLaunchpadByChainId :many
SELECT  * FROM launchpad
WHERE  chain_id = ? ORDER BY created_at ASC;

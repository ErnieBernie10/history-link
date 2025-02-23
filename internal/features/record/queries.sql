-- name: GetRecords :many
select * from record;

-- name: GetRecord :one
select * from record where id = $1;

-- name: CreateRecord :one
insert into record (title, description, location, significance, url, start_date, end_date, status, type)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9)
returning *;

-- name: UpdateRecord :exec
update record
set title = $2, description = $3, location = $4, significance = $5, url = $6, start_date = $7, end_date = $8, status = $9
where id = $1;

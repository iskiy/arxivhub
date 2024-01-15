-- name: CreatePaper :one
INSERT INTO papers (
    arxiv_id,
    title,
    abstract,
    authors,
    date
) VALUES ($1, $2, $3, $4, $5
) RETURNING *;

-- name: GetPaper :one
SELECT * FROM papers
WHERE arxiv_id = $1
LIMIT 1;

-- name: GetPapers :many
SELECT * FROM papers
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeletePaper :exec
DELETE FROM papers
WHERE arxiv_id = $1;

-- name: SavePaperForUser :one
INSERT INTO saved_papers(
     user_id,
     paper_id
) VALUES ($1, $2
) RETURNING *;

-- -- name: GetSavedPapersForUser :many
SELECT p.*
FROM papers AS p
INNER JOIN saved_papers AS sv
ON c.arxiv_id = sv.paper_id
WHERE user_id = $1
ORDER BY c.title
LIMIT $2
OFFSET $3;

-- name: GetSavedPaperIDs :many
SELECT p.paper_id
FROM saved_papers AS p
WHERE user_id = $1
ORDER BY p.created_at;

-- name: DeleteSavedPaper :exec
DELETE from saved_papers
WHERE user_id = $1 AND paper_id = $2;


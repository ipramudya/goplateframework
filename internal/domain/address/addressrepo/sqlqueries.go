package addressrepo

const (
	getOneByIDQuery = `
	SELECT * FROM addresses
	WHERE id = $1
	LIMIT 1`

	updateQuery = `
	UPDATE addresses
	SET street = $1, city = $2, province = $3, postal_code = $4
	WHERE id = $5
	RETURNING *`
)

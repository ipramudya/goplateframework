package addressrepo

const (
	getOneByIDQuery = `
	SELECT * FROM addresses
	WHERE id = $1
	LIMIT 1`

	addOneQuery = `
	INSERT INTO addresses(street, city, province, postal_code)
	VALUES($1, $2, $3, $4)
	RETURNING *`

	updateQuery = `
	UPDATE addresses
	SET street = $1, city = $2, province = $3, postal_code = $4
	WHERE id = $5`
)

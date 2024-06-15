package outletrepo

const (
	getOneByIDQuery = `
	SELECT o.*, a.street, a.city, a.province, a.postal_code 
		FROM outlets o
	INNER JOIN addresses a 
		ON o.address_id = a.id
	WHERE o.id = $1
	LIMIT 1`

	createOutletQuery = `
	WITH new_address AS (
		INSERT INTO addresses(street, city, province, postal_code)
		VALUES($1, $2, $3, $4)
		RETURNING *
	)
	INSERT INTO outlets(name, phone, opening_time, closing_time, address_id)
	VALUES($5, $6, $7, $8, (SELECT id FROM new_address))
	RETURNING
		*,
		(SELECT street FROM new_address),
		(SELECT city FROM new_address),
		(SELECT province FROM new_address),
		(SELECT postal_code FROM new_address)
	`
)

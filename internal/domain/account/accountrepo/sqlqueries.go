package accountrepo

const (
	getOneByEmailQuery = `
	SELECT * FROM accounts 
	WHERE email=$1
	LIMIT 1`

	getOneByIDQuery = `
	SELECT * FROM accounts 
	WHERE id=$1
	LIMIT 1`

	createAccountQuery = `
	INSERT INTO accounts(firstname, lastname, email, password, phone)
	VALUES($1, $2, $3, $4, $5)
	RETURNING *`

	changePasswordQuery = `
	UPDATE accounts
	SET password=$1, updated_at = CURRENT_TIMESTAMP
	WHERE email=$2`
)

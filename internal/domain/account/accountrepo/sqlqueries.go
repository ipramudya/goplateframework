package accountrepo

const (
	GetOneByEmailQuery = `
	SELECT * FROM accounts 
	WHERE email=$1
	LIMIT 1`

	GetOneByIDQuery = `
	SELECT * FROM accounts 
	WHERE id=$1
	LIMIT 1`

	CreateAccountQuery = `
	INSERT INTO accounts(firstname, lastname, email, password, phone)
	VALUES($1, $2, $3, $4, $5)
	RETURNING *`

	ChangePasswordQuery = `
	UPDATE accounts
	SET password=$1
	WHERE email=$2`
)

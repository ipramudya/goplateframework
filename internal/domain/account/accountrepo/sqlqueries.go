package accountrepo

const (
	GetOneByEmailQuery = `
	SELECT id, firstname, lastname, email, password, phone, role, created_at, updated_at
	FROM accounts 
	WHERE email=$1
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

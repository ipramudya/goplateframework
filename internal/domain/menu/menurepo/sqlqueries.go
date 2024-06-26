package menurepo

const (
	createMenuQuery = `
		INSERT INTO menus (name, description, price, is_available, image_url, outlet_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING *
	`

	updateMenuQuery = `
		UPDATE menus
		SET name = $1, description = $2, price = $3, is_available = $4, image_url = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $6
		RETURNING *
	`

	getAllByOutletIDQuery = `
		SELECT * FROM menus
		WHERE outlet_id = $1
		ORDER BY created_at DESC
	`
)

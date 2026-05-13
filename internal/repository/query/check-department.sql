SELECT EXISTS (
	SELECT 1 FROM department WHERE id = ?
)

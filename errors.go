package cfg

type GetEnvError struct {
	field string
}

func (e GetEnvError) Error() string {
	return "cant find env tag for field: " + e.field
}

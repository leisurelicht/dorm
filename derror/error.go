package derror

import "errors"

// Model Error
var (
	DBClientNotFound        = errors.New("database client provider not found")
	DBClientNotValid        = errors.New("database client provider not valid")
	DoesNotExist            = errors.New("does not exist")
)

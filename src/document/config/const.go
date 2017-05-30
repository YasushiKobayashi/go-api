package config

// 400 StatusBadRequest
var BadRequest = map[string]string{
	"error": "BadRequest",
}

// 401 StatusUnauthorized
var Unauthorized = map[string]string{
	"error": "Unauthorized",
}

// 404 StatusNotFound
var NotFound = map[string]string{
	"error": "NotFound",
}

// 406 StatusNotAcceptable
var NotAcceptable = map[string]string{
	"error": "NotAcceptable",
}

// 409 StatusConflict
var ValdError = map[string]string{
	"error": "ValdError",
}

// 500 StatusInternalServerError
var ServerError = map[string]string{
	"error": "InternalServerError",
}

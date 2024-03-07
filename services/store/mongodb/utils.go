package mongodb

import "github.com/dchest/uniuri"

const (
	idLength = 24
)

func randomID() string {
	return uniuri.NewLen(idLength)
}

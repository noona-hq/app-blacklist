package entity

import "time"

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Email     string    `bson:"email,omitempty"`
	CompanyID string    `bson:"companyID,omitempty"`
	Token     Token     `bson:"token,omitempty"`
	CreatedAt time.Time `bson:"createdAt,omitempty"`
	UpdatedAt time.Time `bson:"updatedAt,omitempty"`
	DeletedAt time.Time `bson:"deletedAt,omitempty"`
}

type Token struct {
	RefreshToken         string    `bson:"refreshToken,omitempty"`
	AccessToken          string    `bson:"accessToken,omitempty"`
	AccessTokenExpiresAt time.Time `bson:"expiresAt,omitempty"`
}

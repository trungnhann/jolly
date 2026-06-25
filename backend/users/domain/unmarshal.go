package domain

import "time"

func UnmarshalUser(
	id UserUUID,
	email string,
	name string,
	passwordHash string,
	role Role,
	avatarURL string,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		id:           id,
		email:        email,
		name:         name,
		passwordHash: passwordHash,
		role:         role,
		avatarURL:    avatarURL,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

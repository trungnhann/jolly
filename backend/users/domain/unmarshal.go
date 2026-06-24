package domain

import "time"

func UnmarshalUser(
	id UserUUID,
	email string,
	name string,
	role Role,
	createdAt time.Time,
	updatedAt time.Time,
) User {
	return User{
		id:        id,
		email:     email,
		name:      name,
		role:      role,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}
}

package models

type User struct {
	ID        int    `json:"id"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name"`
}

func (u *User) Save() error {
	_, err := db.Exec("INSERT INTO users (username, password, first_name, last_name) VALUES (?, ?, ?, ?)", u.Username, u.Password, u.FirstName, u.LastName)
	return err
}

func GetUserByUsername(username string) (*User, error) {
	var user User
	err := db.QueryRow("SELECT id, username, password, first_name, last_name FROM users WHERE username = ?", username).Scan(&user.ID, &user.Username, &user.Password, &user.FirstName, &user.LastName)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetAllUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, username, first_name, last_name FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.FirstName, &user.LastName); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

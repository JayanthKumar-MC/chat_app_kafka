package models

import "log"

type Message struct {
	ID          int    `json:"id"`
	SenderID    int    `json:"sender_id"`
	ReceiverID  int    `json:"receiver_id"`
	MessageText string `json:"message_text"`
	Status      string `json:"status"`
	Timestamp   string `json:"timestamp"`
}

func SaveMessage(senderID, receiverID int, messageText string) (int, error) {
	result, err := db.Exec("INSERT INTO messages (sender_id, receiver_id, message_text) VALUES (?, ?, ?)", senderID, receiverID, messageText)
	if err != nil {
		return 0, err
	}
	messageID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(messageID), nil
}

func GetMessagesBetweenUsers(senderID, receiverID int) ([]Message, error) {
	rows, err := db.Query("SELECT id, sender_id, receiver_id, message_text, status, timestamp FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", senderID, receiverID, receiverID, senderID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []Message
	for rows.Next() {
		var message Message
		if err := rows.Scan(&message.ID, &message.SenderID, &message.ReceiverID, &message.MessageText, &message.Status, &message.Timestamp); err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

// UpdateMessageStatus updates the status of a message identified by its ID.
func UpdateMessageStatus(messageID int, status string) error {
	// SQL statement to update the status of a message
	query := "UPDATE messages SET status = ? WHERE id = ?"

	// Execute the SQL statement
	_, err := db.Exec(query, status, messageID)
	if err != nil {
		log.Printf("Error updating message status: %v", err)
		return err
	}

	return nil
}

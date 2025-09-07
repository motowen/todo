package database

// Todo represents a todo item in the database
type Todo struct {
	ID          string			   `bson:"id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" validate:"required"`
	Description string             `bson:"description" json:"description"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   int64          	   `bson:"created_at" json:"created_at"`
	UpdatedAt   int64              `bson:"updated_at" json:"updated_at"`
}
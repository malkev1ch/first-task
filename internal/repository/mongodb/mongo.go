package mongodb

import "go.mongodb.org/mongo-driver/mongo"

// RepositoryMongo type replies for accessing to mongodb database.
type RepositoryMongo struct {
	DB *mongo.Client
}

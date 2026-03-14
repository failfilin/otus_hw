package repository

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/failfilin/otus_hw/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	CollectionRestaurants = "restaurants"
	CollectionMenus       = "menus"
	CollectionDishes      = "dishes"
	CollectionCarts       = "carts"
)

var mongoDB *mongo.Database

var uuidType = reflect.TypeOf(uuid.UUID{})

func encodeUUID(_ bson.EncodeContext, vw bson.ValueWriter, val reflect.Value) error {
	u, ok := val.Interface().(uuid.UUID)
	if !ok {
		return bson.ValueEncoderError{Name: "UUIDEncodeValue", Types: []reflect.Type{uuidType}, Received: val}
	}
	return vw.WriteString(u.String())
}

func decodeUUID(_ bson.DecodeContext, vr bson.ValueReader, val reflect.Value) error {
	if vr.Type() != bson.TypeString {
		return bson.ValueDecoderError{Name: "UUIDDecodeValue", Types: []reflect.Type{uuidType}, Received: val}
	}
	str, err := vr.ReadString()
	if err != nil {
		return err
	}
	u, err := uuid.Parse(str)
	if err != nil {
		return fmt.Errorf("cannot decode %q as UUID: %w", str, err)
	}
	val.Set(reflect.ValueOf(u))
	return nil
}

func buildRegistry() *bson.Registry {
	reg := bson.NewRegistry()
	reg.RegisterTypeEncoder(uuidType, bson.ValueEncoderFunc(encodeUUID))
	reg.RegisterTypeDecoder(uuidType, bson.ValueDecoderFunc(decodeUUID))
	return reg
}

func InitMongo() error {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		uri = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB")
	if dbName == "" {
		dbName = "otus"
	}

	opts := options.Client().
		ApplyURI(uri).
		SetRegistry(buildRegistry()).
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	mongoDB = client.Database(dbName)
	log.Println("Connected to MongoDB:", uri, "db:", dbName)
	return nil
}

func loadCollection[T models.EatType](collection string, safeSl *models.SafeSlice[T]) error {
	safeSl.Mu.Lock()
	defer safeSl.Mu.Unlock()

	cursor, err := mongoDB.Collection(collection).Find(context.Background(), bson.D{})
	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())

	var data []T
	if err := cursor.All(context.Background(), &data); err != nil {
		return err
	}
	if data == nil {
		data = []T{}
	}
	safeSl.Slice = data
	return nil
}

func InsertDocument(collection string, item any) error {
	if mongoDB == nil {
		return nil
	}
	_, err := mongoDB.Collection(collection).InsertOne(context.Background(), item)
	return err
}

func UpdateDocument(collection string, filter, update bson.D) error {
	if mongoDB == nil {
		return nil
	}
	_, err := mongoDB.Collection(collection).UpdateOne(context.Background(), filter, update)
	return err
}

func DeleteDocument(collection string, filter bson.D) error {
	if mongoDB == nil {
		return nil
	}
	_, err := mongoDB.Collection(collection).DeleteOne(context.Background(), filter)
	return err
}

func DeleteManyDocuments(collection string, filter bson.D) error {
	if mongoDB == nil {
		return nil
	}
	_, err := mongoDB.Collection(collection).DeleteMany(context.Background(), filter)
	return err
}

func InitRepository() {
	if err := InitMongo(); err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	if err := loadCollection(CollectionRestaurants, &RestSlice); err != nil {
		log.Printf("Warning: failed to load restaurants: %v", err)
		RestSlice.Slice = []models.Restaurant{}
	}
	if err := loadCollection(CollectionMenus, &MenuSlice); err != nil {
		log.Printf("Warning: failed to load menus: %v", err)
		MenuSlice.Slice = []models.Menu{}
	}
	if err := loadCollection(CollectionDishes, &DishSlice); err != nil {
		log.Printf("Warning: failed to load dishes: %v", err)
		DishSlice.Slice = []models.Dish{}
	}
	if err := loadCollection(CollectionCarts, &CartSlice); err != nil {
		log.Printf("Warning: failed to load carts: %v", err)
		CartSlice.Slice = []models.Cart{}
	}
}

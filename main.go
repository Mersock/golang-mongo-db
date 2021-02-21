package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Product struct {
	// ID          primitive.ObjectID `json:"_id" bson:"_id"`
	Name        string   `json:"product_name" bson:"product_name"`
	Price       int      `json:"price" bson:"price"`
	Currency    string   `json:"currency" bson:"currency"`
	Quantity    string   `json:"quantity" bson:"quantity"`
	Discount    int      `json:"discount,omitempty" bson:"discount,omitempty"`
	Vendor      string   `json:"vendor" bson:"vendor"`
	Accessories []string `json:"accessories,omitempty" bson:"accessories,omitempty"`
	SkuID       string   `json:"sku_id" bson:"sku_id"`
}

var iPhone10 = Product{
	// ID:          primitive.NewObjectID(),
	Name:        "Iphone 10",
	Price:       900,
	Currency:    "USD",
	Quantity:    "40",
	Vendor:      "apple",
	Accessories: []string{"charger", "headset", "slotpenar"},
	SkuID:       "1234",
}

var trimmer = Product{
	// ID:          primitive.NewObjectID(),
	Name:        "easy trimter",
	Price:       120,
	Currency:    "USD",
	Quantity:    "300",
	Vendor:      "Philips",
	Discount:    7,
	Accessories: []string{"charger", "comb", "bladeset", "clean oil"},
	SkuID:       "2345",
}

func main() {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://root:123456@mongo/?authSource=admin"))
	if err != nil {
		fmt.Println(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	if err != nil {
		fmt.Println(err)
	}
	db := client.Database("tronics")
	collection := db.Collection("products")

	//Binary Encoded Json.
	//bson.D = ordered document bson.D{{"name", "knz"},{"surname", "phumthawan"}}
	//bson.M = unordered document/map bson.M{{"hello":"world","foo","bar"}}
	//bson.A = array bson.A{"hello","world"}
	//bson.E = usually used as a element inside bson.D

	// ***** Create *****

	// Using struct
	res1, err := collection.InsertOne(context.Background(), trimmer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Insert Using Strut")
	fmt.Println(res1.InsertedID.(primitive.ObjectID).Timestamp())

	// Using bson.D
	res2, err := collection.InsertOne(context.Background(), bson.D{
		{"name", "knz"},
		{"surname", "phumthawan"},
		{"hobby", bson.A{"game", "soccer", "run"}},
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Insert Using bson.D")
	fmt.Println(res2.InsertedID.(primitive.ObjectID).Timestamp())

	//Using bson.M
	res3, err := collection.InsertOne(context.Background(), bson.M{
		"name":    "eric",
		"surname": "carname",
		"hobbies": bson.A{"videogame", "music", "soccer"},
	})
	fmt.Println("Insert Using bson.M")
	fmt.Println(res3.InsertedID.(primitive.ObjectID).Timestamp())

	resMany, err := collection.InsertMany(context.Background(), []interface{}{iPhone10, trimmer})
	fmt.Println("Insert Many")
	fmt.Println(resMany.InsertedIDs)

	// ***** Read *****

	//Equality operator using FindOne
	var findOne Product
	err = collection.FindOne(context.Background(), bson.M{"price": 900}).Decode(&findOne)
	fmt.Println("Equality operator using FindOne")
	fmt.Println(findOne)

	//Comparison operator using Find
	var find Product
	fmt.Println("----- Comparison operator using Find -----")
	findCursor, err := collection.Find(context.Background(), bson.M{"price": bson.M{"$gt": 100}})
	for findCursor.Next(context.Background()) {
		err := findCursor.Decode(&find)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(find.Name)
	}

	// Logical operator using Find
	var findLogic Product
	logicFilter := bson.M{
		"$and": bson.A{
			bson.M{"price": bson.M{"$gt": 100}},
			bson.M{"quantity": bson.M{"$gt": 30}},
		},
	}
	fmt.Println("----- Logical operator using Find -----")
	findLogicRes, err := collection.Find(context.Background(), logicFilter)
	for findLogicRes.Next(context.Background()) {
		err := findLogicRes.Decode(&findLogic)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(findLogic.Name)
	}

	// Element operator using Find
	var findElement Product
	elementFilter := bson.M{
		"accessories": bson.M{"$exists": true},
	}
	fmt.Println("----- Element operator using Find -----")
	findElementRes, err := collection.Find(context.Background(), elementFilter)
	for findElementRes.Next(context.Background()) {
		err := findElementRes.Decode(&findElement)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(findElement.Name)
	}

	// Array operator using Find
	var findArray Product
	arrayFilter := bson.M{
		"accessories": bson.M{"$all": bson.A{"charger"}},
	}
	fmt.Println("----- Array operator using Find -----")
	findArrayRes, err := collection.Find(context.Background(), arrayFilter)
	for findArrayRes.Next(context.Background()) {
		err := findArrayRes.Decode(&findArray)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(findArray.Name)
	}

	//Update operator for Field
	fmt.Println("----- Update operator for field -----")
	updateFieldCon := bson.M{"$set": bson.M{"IsEssential": "false"}}
	updateFieldRes, err := collection.UpdateMany(context.Background(), bson.M{}, updateFieldCon)
	fmt.Println(updateFieldRes.ModifiedCount)

	//Update operator for Array
	updateArrayCon := bson.M{"$addToSet": bson.M{"accessories": "manual"}}
	updateArrayRes, err := collection.UpdateMany(context.Background(), arrayFilter, updateArrayCon)
	fmt.Println(updateArrayRes.MatchedCount)

	//Update operator for field - multiple operators
	fmt.Println("----- Update operator for field - multiple operators -----")
	incCon := bson.M{
		"$mul": bson.M{
			"price": 1.20,
		},
	}
	incRes, err := collection.UpdateMany(context.Background(), bson.M{}, incCon)
	fmt.Println(incRes.MatchedCount)

	//Delete operation
	fmt.Println("----- Delete operation -----")
	delRes, err := collection.DeleteMany(context.Background(), arrayFilter)
	fmt.Println(delRes.DeletedCount)

	if err != nil {
		fmt.Println(err)
	}
}

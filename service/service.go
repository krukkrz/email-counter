package service

import (
	"context"
	"email-counter/connector"
	"encoding/json"
	"log"
	"strconv"

	"github.com/gofiber/fiber"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Iteration struct {
	_id             string `bson:"id,omitempty"`
	Iteration       int    `bson:"iteration,omitempty"`
	EmailsSentCount int    `bson:"emailsSentCount,omitempty"`
	Archived        bool   `bson:"archived,omitempty"`
}

func HealthCheck(c *fiber.Ctx) {
	log.Printf("health check endpoint called from %s", c.IP())
}

func CreateList(c *fiber.Ctx) {
	setDefaultHeaders(c)
	collection := getCollection()

	var list Iteration
	json.Unmarshal([]byte(c.Body()), &list)

	log.Printf("Creating list for iteration: %v", list.Iteration)
	res, err := collection.InsertOne(context.Background(), list)
	handleServerError(err, c)

	response, _ := json.Marshal(res)
	c.Send(response)
}

func GetAll(c *fiber.Ctx) {
	setDefaultHeaders(c)
	collection := getCollection()
	iteration := parseIterationFromPath(c)

	log.Printf("getting results for iteration: %v", iteration)
	cur, err := collection.Find(context.Background(), bson.M{})

	handleServerError(err, c)

	defer cur.Close(context.Background())

	var results []bson.M
	cur.All(context.Background(), &results)
	json, _ := json.Marshal(results)
	c.Send(json)
}

func GetListReportByIteration(c *fiber.Ctx) {
	setDefaultHeaders(c)
	collection := getCollection()
	iteration := parseIterationFromPath(c)

	log.Printf("getting results for iteration: %v", iteration)
	results := findByIterationNumber(iteration, c, collection)

	json, _ := json.Marshal(results)
	c.Send(json)
}

func UpdateEmailsSentCounter(c *fiber.Ctx) {
	setDefaultHeaders(c)
	collection := getCollection()
	iteration := parseIterationFromPath(c)

	log.Printf("Updating counter of iteration: %v", iteration)
	increaseCounterByIterationNumber(iteration, c, collection)
}

func ArchiveIteration(c *fiber.Ctx) {
	setDefaultHeaders(c)
	collection := getCollection()
	objID := getObjectIdFromPath(c)

	archiveIteration(objID, c, collection)
}

func getObjectIdFromPath(c *fiber.Ctx) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(c.Params("id"))
	handleBadRequestError(err, c)
	return objID
}

func setDefaultHeaders(c *fiber.Ctx) {
	c.Fasthttp.Response.Header.Add("Access-Control-Allow-Origin", "*")
	c.Fasthttp.SetContentType("application/json")
}

func getCollection() *mongo.Collection {
	return connector.GetMongoDbCollection()
}

func handleServerError(err error, c *fiber.Ctx) {
	if err != nil {
		c.Status(500).Send(err)
		return
	}
}

func handleBadRequestError(err error, c *fiber.Ctx) {
	if err != nil {
		c.Status(400).Send(err)
		return
	}
}

func handle404Error(results []bson.M, c *fiber.Ctx) {
	if results == nil {
		c.SendStatus(404)
	}
}

func parseIterationFromPath(c *fiber.Ctx) int {
	iterationString := c.Params("iteration")
	iteration, err := strconv.Atoi(iterationString)
	handleBadRequestError(err, c)
	return iteration
}

func filterByIterationNumber(iteration int) bson.M {
	return bson.M{
		"iteration": iteration,
		"archived":  nil,
	}
}

func findByIterationNumber(iteration int, c *fiber.Ctx, collection *mongo.Collection) []bson.M {
	filter := filterByIterationNumber(iteration)
	cur, err := collection.Find(context.Background(), filter)

	handleServerError(err, c)

	defer cur.Close(context.Background())

	var results []bson.M
	cur.All(context.Background(), &results)

	handle404Error(results, c)

	return results
}

func increaseCounterByIterationNumber(iteration int, c *fiber.Ctx, collection *mongo.Collection) {
	filter := filterByIterationNumber(iteration)

	update := bson.M{
		"$inc": bson.M{"emailsSentCount": 1},
	}

	result, err := collection.UpdateMany(
		c.Context(),
		filter,
		update,
	)

	handleServerError(err, c)
	log.Printf("Updated %v documents", result.ModifiedCount)
}

func archiveIteration(objID primitive.ObjectID, c *fiber.Ctx, collection *mongo.Collection) {
	log.Printf("Archiving iteration with ID: %v", objID)

	update := bson.M{
		"$set": bson.M{"archived": true},
	}

	filter := bson.M{"_id": objID}

	result, err := collection.UpdateOne(
		c.Context(),
		filter,
		update,
	)

	handleServerError(err, c)
	log.Printf("Updated %v documents", result.ModifiedCount)
}

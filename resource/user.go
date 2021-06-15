package resource

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joaoteixeira/go-auth-service-v2/helper"
	"github.com/twinj/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Name      string             `bson:"name" json:"name" validate:"required"`
	Email     string             `bson:"email" json:"email" validate:"email,required"`
	Password  string             `bson:"password" json:"password" validate:"required,min=6"`
	Salt      string             `bson:"salt" json:"salt"`
	Roles     []string           `bson:"roles" json:"roles"`
	Active    bool               `bson:"active" json:"active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

var collection = helper.ConnectDB().Collection("users")
var validate *validator.Validate

func init() {
	validate = validator.New()
}

func (rs User) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Put("/", rs.Delete)

	r.Route("/{id}", func(r chi.Router) {
		// r.Use(rs.TodoCtx) // lets have a todos map, and lets actually load/manipulate
		r.Get("/", rs.Get)       // GET /todos/{id} - read a single todo by :id
		r.Put("/", rs.Update)    // PUT /todos/{id} - update a single todo by :id
		r.Delete("/", rs.Delete) // DELETE /todos/{id} - delete a single todo by :id
	})

	return r
}

func (rs User) List(w http.ResponseWriter, r *http.Request) {

	var users []User

	cur, errFind := collection.Find(context.TODO(), bson.D{{}})

	if errFind != nil {
		helper.ResponseWithError(w, 500, errFind.Error())
		return
	}

	for cur.Next(context.TODO()) {
		var t User
		err := cur.Decode(&t)

		if err != nil {
			helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		users = append(users, t)
	}

	cur.Close(context.TODO())

	js, err := json.Marshal(users)

	if err != nil {
		helper.ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Write(js)
}

func (rs User) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	validationErr := validate.Struct(user)

	if validationErr != nil {
		helper.ResponseWithError(w, http.StatusBadRequest, validationErr.Error())
		return
	}

	user.ID = primitive.NewObjectID()
	user.Salt = uuid.NewV4().String()
	user.Password = helper.HashPassword(user.Password, user.Salt)
	user.Active = true
	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	useri, err := collection.InsertOne(context.TODO(), &user)

	if err != nil {
		helper.ResponseWithError(w, 500, "Error in insert db")
		return
	}

	helper.Response(w, 201, useri)
}

func (rs User) Get(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa get"))
}

func (rs User) Update(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa update"))
}

func (rs User) Delete(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("aaa delete"))
}

package main

import (
	"encoding/json"
	"fmt"
	"github.com/graphql-go/graphql"
	gqlhandler "github.com/graphql-go/graphql-go-handler"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

// Post Structure
type Post struct {
	UserID int `json:"userId"`
	ID int `json:"id"`
	Title string `json:"title"`
	Body string `json:"body"`
}

// Comment Structure
type Comment struct {
	PostID int `json:"postId"`
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Body string `json:"body"`
}

// Consume all Posts from jsonplaceholder
func getPosts() []Post {
	response, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject []Post
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

// Consume Comments
func getComments() []Comment {
	response, err := http.Get("https://jsonplaceholder.typicode.com/comments")
	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseObject []Comment
	json.Unmarshal(responseData, &responseObject)

	return responseObject
}

// Setup Graphql Schema
func setupGraphql() graphql.Schema {
	posts := getPosts()
	comments := getComments()
	var postType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Post",
			Fields: graphql.Fields{
				"userId": &graphql.Field{
					Type: graphql.Int,
				},
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"title": &graphql.Field{
					Type: graphql.String,
				},
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		})
	var commentType = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Comment",
			Fields: graphql.Fields{
				"postId": &graphql.Field{
					Type: graphql.Int,
				},
				"id": &graphql.Field{
					Type: graphql.Int,
				},
				"name": &graphql.Field{
					Type: graphql.String,
				},
				"email": &graphql.Field{
					Type: graphql.String,
				},
				"body": &graphql.Field{
					Type: graphql.String,
				},
			},
		})

	fields := graphql.Fields{
		"post": &graphql.Field{
			Type:        postType,
			Description: "Get Post by ID",
			Args: graphql.FieldConfigArgument{
				"post_id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				id, ok := p.Args["post_id"].(int)
				if ok {
					for _, post := range posts {
						if int(post.ID) == id {
							return post, nil
						}
					}
				}
				return nil, nil
			},
		},
		"posts": &graphql.Field{
			Type: graphql.NewList(postType),
			Description: "Get All Posts",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return posts, nil
			},
		},

		"comment": &graphql.Field{
			Type:        commentType,
			Description: "Search for a comment",
			Args: graphql.FieldConfigArgument{
				"search": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				search, ok := p.Args["search"].(string)
				if ok {
					for _, comment := range comments {
						if strings.Contains(comment.Body, search) {
							return comment, nil
						} else if strings.Contains(comment.Email, search) {
							return comment, nil
						} else if comment.Name == search {
							return comment, nil
						}
					}
				}
				return nil, nil
			},
		},
		"comments": &graphql.Field{
			Type: graphql.NewList(commentType),
			Description: "Get All Comments",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return comments, nil
			},
		},
	}

	rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: fields}
	schemaConfig := graphql.SchemaConfig{Query: graphql.NewObject(rootQuery)}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("Failed to create new GraphQl Schema. err %v", err)
	}
	return schema
}


func main()  {

	schema := setupGraphql()

	// enable or disable GraphiQL or Graphql Playground
	h := gqlhandler.New(&gqlhandler.Config{
		Schema: &schema,
		Pretty: true,
		GraphiQL: false,
		Playground: true,
	})

	// serve a GraphQL endpoint at `/graphql`
	http.Handle("/graphql", h)

	// and serve!
	log.Fatal(http.ListenAndServe(":4000", nil))

}
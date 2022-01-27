package main

import (
	"cesargdd/graph-subcriptions/graph"
	"cesargdd/graph-subcriptions/graph/generated"
	"cesargdd/graph-subcriptions/pg"
	"cesargdd/graph-subcriptions/userAuth"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
	"github.com/rs/cors"
)

const defaultPort = "8080"

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

var conn = pg.Connect()
var db = pg.New(conn)

func main() {

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)
	router.Use(userAuth.Middleware())
	// Use New instead of NewDefaultServer in order to have full control over defining transports
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: true,
		},
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			token := initPayload.Authorization()
			user := userAuth.ForContext(ctx)
			if user == nil {
				return nil, fmt.Errorf("access denied")
			}

			Id, _ := strconv.Atoi(user.ID)
			// get the user from the database
			_, err := db.GetUserById(ctx, int32(Id))
			if err != nil {
				fmt.Println("error user db for websockets")
			}

			nctx, _ := context.WithCancel(ctx)

			// put it in context
			userCtx := context.WithValue(nctx, userCtxKey, &pg.AuthResponse{
				AuthToken: &pg.AuthToken{
					AccessToken: token,
				},
				User: user,
			})

			// and return it so the resolvers can see it
			return userCtx, nil
		},
	})
	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	router.Handle("/query", (srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

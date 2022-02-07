package main

import (
	"cesargdd/graph-subcriptions/graph"
	"cesargdd/graph-subcriptions/graph/generated"
	"cesargdd/graph-subcriptions/jwt"
	"cesargdd/graph-subcriptions/pg"
	"cesargdd/graph-subcriptions/userAuth"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
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

func main() {

	var conn = pg.Connect()
	var db = pg.New(conn)

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	router := chi.NewRouter()

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
		AllowedHeaders:   []string{"*"},
	}).Handler)
	router.Use(userAuth.Middleware())
	// Use New instead of NewDefaultServer in order to have full control over defining transports
	srv := handler.New(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		Upgrader: websocket.Upgrader{
			HandshakeTimeout: time.Minute,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			EnableCompression: true,
		},
		KeepAlivePingInterval: 10 * time.Second,
		InitFunc: func(ctx context.Context, initPayload transport.InitPayload) (context.Context, error) {
			payload := initPayload.Authorization()
			if payload == "" {
				fmt.Println("Noa acces graded")
			}
			username, err := jwt.ParseToken(payload)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(username)
			user, err := db.GetIdUserByUsername(ctx, username)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(user)
			cont := context.TODO()
			usrCtx := context.WithValue(cont, userAuth.UserCtxKey, &user)
			// log.Println(usrCtx.Value(userAuth.UserCtxKey))
			return usrCtx, nil
		},
	})
	srv.Use(extension.Introspection{})

	router.Handle("/", playground.Handler("GraphQL playground", "/query"))

	router.Handle("/query", (srv))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

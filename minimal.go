package main

import (
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humachi"
	"github.com/go-chi/chi/v5"
)

func helloWorld(ctx huma.Context, input struct{}) (*struct{}, error) {
	huma.WriteTextResponse(ctx, http.StatusOK, "Hello, World!")
	return &struct{}{}, nil
}

func main() {
	r := chi.NewRouter()
	api := humachi.New(r, huma.DefaultConfig("Minimal API", "1.0.0"))

	huma.Register(api, huma.Operation{
		OperationID: "hello-world",
		Summary:     "Hello World endpoint",
		Method:      http.MethodGet,
		Path:        "/hello",
	}, helloWorld)

	http.ListenAndServe(":8080", r)
}

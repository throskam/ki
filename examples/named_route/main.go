package main

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/throskam/ki"
	"github.com/throskam/ki/middlewares"
)

func main() {
	router := ki.NewRouter()

	router.Use(middlewares.Locator(router))

	noopHandler := func(w http.ResponseWriter, r *http.Request) {}

	router.Get("/posts", noopHandler, ki.WithName("list-posts"))
	router.Post("/posts", noopHandler, ki.WithName("create-post"))
	router.Get("/posts/{postID}", noopHandler, ki.WithName("get-post"))
	router.Delete("/posts/{postID}", noopHandler, ki.WithName("delete-post"))
	router.Get("/posts/{postID}/comments", noopHandler, ki.WithName("list-post-comments"))
	router.Post("/posts/{postID}/comments", noopHandler, ki.WithName("create-post-comment"))
	router.Get("/posts/{postID}/comments/{commentID}", noopHandler, ki.WithName("get-post-comment"))
	router.Delete("/posts/{postID}/comments/{commentID}", noopHandler, ki.WithName("delete-post-comment"))

	router.Get("/named-routes", func(w http.ResponseWriter, r *http.Request) {
		listPostsLocation := ki.GetLocation(r.Context(), "list-posts").WithQuery(url.Values{"page": []string{"1"}, "limit": []string{"10"}})
		createPostLocation := ki.GetLocation(r.Context(), "create-post")
		getPostLocation := ki.GetLocation(r.Context(), "get-post").WithPathParams("1234")
		deletePostLocation := ki.GetLocation(r.Context(), "delete-post").WithPathParams("1234")
		listPostCommentsLocation := ki.GetLocation(r.Context(), "list-post-comments").WithPathParams("1234").WithQueryParam("sort", "newest")
		createPostCommentLocation := ki.GetLocation(r.Context(), "create-post-comment").WithPathParams("1234")
		getPostCommentLocation := ki.GetLocation(r.Context(), "get-post-comment").WithPathParams("1234", "5678")
		deletePostCommentLocation := ki.GetLocation(r.Context(), "delete-post-comment").WithPathParams("1234", "5678")

		_, _ = fmt.Fprintf(w, "list posts: %s %s\n", listPostsLocation.Method(), listPostsLocation.URL())
		_, _ = fmt.Fprintf(w, "create post: %s %s\n", createPostLocation.Method(), createPostLocation.URL())
		_, _ = fmt.Fprintf(w, "get post: %s %s\n", getPostLocation.Method(), getPostLocation.URL())
		_, _ = fmt.Fprintf(w, "delete post: %s %s\n", deletePostLocation.Method(), deletePostLocation.URL())
		_, _ = fmt.Fprintf(w, "list post comments: %s %s\n", listPostCommentsLocation.Method(), listPostCommentsLocation.URL())
		_, _ = fmt.Fprintf(w, "create post comment: %s %s\n", createPostCommentLocation.Method(), createPostCommentLocation.URL())
		_, _ = fmt.Fprintf(w, "get post comment: %s %s\n", getPostCommentLocation.Method(), getPostCommentLocation.URL())
		_, _ = fmt.Fprintf(w, "delete post comment: %s %s\n", deletePostCommentLocation.Method(), deletePostCommentLocation.URL())
	})

	_ = http.ListenAndServe(":8080", router)
}

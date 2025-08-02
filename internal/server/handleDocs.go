package server

import (
	"html/template"
	"io"
	"net/http"

	_ "embed"
)

//go:embed openapi.json
var spec []byte

const docsPageTemplate string = `<!DOCTYPE html>
<html>
  <head>
    <title>Recurse Print API</title>
    <meta charset="utf-8" />
    <meta
      name="viewport"
      content="width=device-width, initial-scale=1" />
  </head>

  <body>
    <div id="app"></div>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
    <script>
      Scalar.createApiReference('#app', {
        // The URL of the OpenAPI/Swagger document
        url: '{{ .SpecURL }}'
      })
    </script>
  </body>
</html>
`

type docsPageData struct {
	Title      string
	FaviconURL string
	SpecURL    string
}

func genDocsPage(data docsPageData, out io.Writer) error {
	t, err := template.New("docsPageHTML").Parse(docsPageTemplate)
	if err != nil {
		return err
	}

	err = t.Execute(out, data)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) handleSchema(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	_, _ = res.Write(spec)
}

func (s *Server) handleDocs(res http.ResponseWriter, req *http.Request) {
	err := genDocsPage(docsPageData{
		Title:   "Recurse Print API",
		SpecURL: "/openapi.json",
	}, res)
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
	}
}

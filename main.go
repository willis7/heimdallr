package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/graphql-go/graphql"
)

type Project struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

type Build struct {
	ID       string `json:"id,omitempty"`
	Project  string `json:"project"`
	ModuleID string `json:"moduleid"`
}

// Test Data
var projects []Project = []Project{
	Project{
		ID:   "1",
		Name: "Thor",
	},
	Project{
		ID:   "2",
		Name: "Loki",
	},
}

var builds []Build = []Build{
	Build{
		ID:       "thor_build_1",
		Project:  "1",
		ModuleID: "com.asgard:thor:0.1",
	},
	Build{
		ID:       "thor_build_2",
		Project:  "1",
		ModuleID: "com.asgard:thor:0.2",
	},
	Build{
		ID:       "loki_build_2",
		Project:  "2",
		ModuleID: "com.asgard:loki:0.1",
	},
}

func Filter(builds []Build, f func(Build) bool) []Build {
	vsf := make([]Build, 0)
	for _, v := range builds {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}

func main() {
	buildType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Build",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"project": &graphql.Field{
				Type: graphql.String,
			},
			"moduleid": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	projectType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Project",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
		},
	})
	rootQuery := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"builds": &graphql.Field{
				Type: graphql.NewList(buildType),
				Args: graphql.FieldConfigArgument{
					"project": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					project := params.Args["project"].(string)
					filtered := Filter(builds, func(v Build) bool {
						return strings.Contains(v.Project, project)
					})
					return filtered, nil
				},
			},
			"projects": &graphql.Field{
				Type: graphql.NewList(projectType),
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					id := params.Args["id"].(string)
					for _, p := range projects {
						if p.ID == id {
							return p, nil
						}
					}
					return nil, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(graphql.SchemaConfig{
		Query: rootQuery,
	})
	http.HandleFunc("/graphql", func(w http.ResponseWriter, r *http.Request) {
		result := graphql.Do(graphql.Params{
			Schema:        schema,
			RequestString: r.URL.Query().Get("query"),
		})
		json.NewEncoder(w).Encode(result)
	})
	http.ListenAndServe(":12345", nil)
}

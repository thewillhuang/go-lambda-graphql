package gql

import (
	"Lambda/services/generate"
	"context"
	"database/sql"
	"encoding/json"
	"go-lambda-graphql/config"
	"io/ioutil"
	"testing"

	_ "github.com/lib/pq"
	"github.com/malisit/kolpa"
	graphql "github.com/neelance/graphql-go"
	"github.com/neelance/graphql-go/gqltesting"
	"github.com/tidwall/gjson"
	"github.com/volatiletech/sqlboiler/boil"
)

func TestUserResolvers(t *testing.T) {
	// setup db
	rawSchema, _ := ioutil.ReadFile("schema.gql")
	schema := graphql.MustParseSchema(string(rawSchema), &Resolver{})
	db, _ := sql.Open("postgres", config.ConnectionString)
	// boil.DebugMode = !config.Production
	boil.SetDB(db)
	fake := kolpa.C()
	t.Run("reject wrong email", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		password := fake.LoremSentence()
		name := fake.Name()
		email2 := fake.Email()

		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					query {
						viewer(jwt: "` + jwt + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"viewer":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email2 + `", password: "` + password + `") {}
			}
			`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		err := gjson.ParseBytes(result2).Get("errors.0.message").String()
		if err != "wrong email or password combination" {
			t.Errorf("expected error")
		}
	})

	t.Run("reject invalid email", func(t *testing.T) {
		t.Parallel()
		email := fake.Name()
		password := fake.LoremSentence()
		name := fake.Name()
		test := gqltesting.Test{
			Schema: schema,
			Query: `
			mutation {
				signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
					name
					email
				}
			}
		`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		err := gjson.ParseBytes(result).Get("errors.0.message").String()
		if err != "Email: must be a valid email address." {
			t.Errorf("expected error")
		}
	})

	t.Run("reject wrong password", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		password := fake.LoremSentence()
		password2 := fake.LoremSentence()
		name := fake.Name()

		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					query {
						viewer(jwt: "` + jwt + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"viewer":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})

		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email + `", password: "` + password2 + `") {}
			}
			`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		err := gjson.ParseBytes(result2).Get("errors.0.message").String()
		if err != "wrong email or password combination" {
			t.Errorf("expected error")
		}
	})

	t.Run("update user email and password", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		password := fake.LoremSentence()
		name := fake.Name()
		email2 := fake.Email()
		password2 := fake.LoremSentence()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				mutation {
					updateUser(jwt: "` + jwt + `", email: "` + email2 + `", password: "` + password2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"updateUser":{"email":"` + email2 + `","name":"` + name + `"}}
			`,
			},
		})
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email2 + `", password: "` + password2 + `") {}
			}
		`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		jwt2 := gjson.ParseBytes(result2).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				query {
					viewer(jwt: "` + jwt2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"viewer":{"email":"` + email2 + `","name":"` + name + `"}}
			`,
			},
		})
	})
	t.Run("update email", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		email2 := fake.Email()
		password := fake.LoremSentence()
		name := fake.Name()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				mutation {
					updateUser(jwt: "` + jwt + `", email: "` + email2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"updateUser":{"email":"` + email2 + `","name":"` + name + `"}}
			`,
			},
		})
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email2 + `", password: "` + password + `") {}
			}
		`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		jwt2 := gjson.ParseBytes(result2).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				query {
					viewer(jwt: "` + jwt2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"viewer":{"email":"` + email2 + `","name":"` + name + `"}}
			`,
			},
		})
	})

	t.Run("reject update if invalid email", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		email2 := generate.GenerateRandomString(32)
		password := fake.LoremSentence()
		name := fake.Name()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
				mutation {
					updateUser(jwt: "` + jwt + `", email: "` + email2 + `") {
						name
						email
					}
				}
			`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		err := gjson.ParseBytes(result2).Get("errors.0.message").String()
		if err != "Email: must be a valid email address." {
			t.Errorf("expected error")
		}
	})

	t.Run("update name", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		password := fake.LoremSentence()
		name := fake.Name()
		name2 := fake.Name()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				mutation {
					updateUser(jwt: "` + jwt + `", name: "` + name2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"updateUser":{"email":"` + email + `","name":"` + name2 + `"}}
			`,
			},
		})
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email + `", password: "` + password + `") {}
			}
		`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		jwt2 := gjson.ParseBytes(result2).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				query {
					viewer(jwt: "` + jwt2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"viewer":{"email":"` + email + `","name":"` + name2 + `"}}
			`,
			},
		})
	})

	t.Run("update inputting same values", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		password := fake.LoremSentence()
		name := fake.Name()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				mutation {
					updateUser(jwt: "` + jwt + `", password: "` + password + `", email: "` + email + `", name: "` + name + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"updateUser":{"email":"` + email + `","name":"` + name + `"}}
			`,
			},
		})
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email + `", password: "` + password + `") {}
			}
		`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		jwt2 := gjson.ParseBytes(result2).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				query {
					viewer(jwt: "` + jwt2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"viewer":{"email":"` + email + `","name":"` + name + `"}}
			`,
			},
		})
	})

	t.Run("update password", func(t *testing.T) {
		t.Parallel()
		name := fake.Name()
		email := fake.Email()
		password := fake.LoremSentence()
		password2 := fake.LoremSentence()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				mutation {
					updateUser(jwt: "` + jwt + `", password: "` + password2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"updateUser":{"email":"` + email + `","name":"` + name + `"}}
			`,
			},
		})
		test2 := gqltesting.Test{
			Schema: schema,
			Query: `
			{
				jwt(email: "` + email + `", password: "` + password2 + `") {}
			}
		`,
		}
		result2, _ := json.Marshal(test2.Schema.Exec(context.Background(), test2.Query, test2.OperationName, test2.Variables))
		jwt2 := gjson.ParseBytes(result2).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
				query {
					viewer(jwt: "` + jwt2 + `") {
						name
						email
					}
				}
			`,
				ExpectedResult: `
				{"viewer":{"email":"` + email + `","name":"` + name + `"}}
			`,
			},
		})
	})

	t.Run("user creation", func(t *testing.T) {
		t.Parallel()
		email := fake.Email()
		password := fake.LoremSentence()
		name := fake.Name()

		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					mutation {
						signup(name: "` + name + `", email: "` + email + `", password: "` + password + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"signup":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				{
					jwt(email: "` + email + `", password: "` + password + `") {}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		jwt := gjson.ParseBytes(result).Get("data.jwt").String()
		gqltesting.RunTests(t, []*gqltesting.Test{
			{
				Schema: schema,
				Query: `
					query {
						viewer(jwt: "` + jwt + `") {
							name
							email
						}
					}
				`,
				ExpectedResult: `
					{"viewer":{"email":"` + email + `","name":"` + name + `"}}
				`,
			},
		})
	})

	t.Run("reject melformed jwt", func(t *testing.T) {
		t.Parallel()
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				query {
					viewer(jwt: "` + "<<<<<<<<<<<<<<<<<<<<<<jwt" + `") {
						name
						email
					}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		err := gjson.ParseBytes(result).Get("errors.0.message").String()
		if err != "token contains an invalid number of segments" {
			t.Errorf("expected error")
		}
	})

	t.Run("reject randomly generated string as jwt", func(t *testing.T) {
		t.Parallel()
		fakejwt := generate.GenerateRandomString(32)
		test := gqltesting.Test{
			Schema: schema,
			Query: `
				query {
					viewer(jwt: "` + fakejwt + `") {
						name
						email
					}
				}
			`,
		}
		result, _ := json.Marshal(test.Schema.Exec(context.Background(), test.Query, test.OperationName, test.Variables))
		err := gjson.ParseBytes(result).Get("errors.0.message").String()
		if err != "token contains an invalid number of segments" {
			t.Errorf("expected error")
		}
	})
}

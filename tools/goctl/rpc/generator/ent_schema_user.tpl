package schema

import (
    "entgo.io/ent"
    "entgo.io/ent/dialect/entsql"
    "entgo.io/ent/schema"
    "entgo.io/ent/schema/field"
    "entgo.io/ent/schema/index"

    "github.com/toby1991/go-zero-utils/db/ent/schema/mixin"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entsql.Annotation{Table: "au_user"},
	}
}

func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		mixin.TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Uint64("id").StorageKey("user_id").Positive().Comment("用户ID"),
		field.String("username").StorageKey("user_name").Comment("用户名").MaxLen(100),
		field.String("email").StorageKey("user_email").Comment("邮箱").MaxLen(191),
		field.String("telephone").StorageKey("user_telephone").Comment("手机号").MaxLen(20),
		field.String("password").StorageKey("user_password").Comment("密码").MaxLen(64),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("username").Unique(),
		index.Fields("email").Unique(),
		index.Fields("telephone").Unique(),
	}
}

// interceptors for the User.
func (User) Interceptors() []ent.Interceptor {
	return []ent.Interceptor{
		ent.InterceptFunc(func(next ent.Querier) ent.Querier {
			return ent.QuerierFunc(func(ctx context.Context, query ent.Query) (ent.Value, error) {
				// Do something before the query execution.
				value, err := next.Query(ctx, query)
				// Do something after the query execution.
				return value, err
			})
		}),
	}
}

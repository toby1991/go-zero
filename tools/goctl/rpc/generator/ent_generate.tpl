package ent

//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/upsert,sql/lock,sql/versioned-migration,intercept,schema/snapshot,sql/execquery ./schema
//https://entgo.io/docs/feature-flags/#row-level-locks  https://entgo.io/docs/versioned-migrations https://entgo.io/docs/interceptors

// generate graph
// ent describe ./ent/schema
// go get github.com/hedwigz/entviz/cmd/entviz
// go install github.com/hedwigz/entviz/cmd/entviz
// entviz ./ent/schema

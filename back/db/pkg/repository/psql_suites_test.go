// Code generated by SQLBoiler 4.15.0 (https://github.com/volatiletech/sqlboiler). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

package repository

import "testing"

func TestUpsert(t *testing.T) {
	t.Run("SchemaMigrations", testSchemaMigrationsUpsert)

	t.Run("TBlogs", testTBlogsUpsert)

	t.Run("TUsers", testTUsersUpsert)
}
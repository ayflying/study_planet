// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// SchemaMigrations is the golang structure for table schema_migrations.
type SchemaMigrations struct {
	Version   int64       `json:"version"   orm:"version"    description:""` //
	Name      string      `json:"name"      orm:"name"       description:""` //
	AppliedAt *gtime.Time `json:"appliedAt" orm:"applied_at" description:""` //
}

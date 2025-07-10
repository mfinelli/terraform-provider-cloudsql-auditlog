-- Copyright (c) Mario Finelli
-- SPDX-License-Identifier: MPL-2.0

-- name: GetAllAuditRules :many
SELECT * FROM audit_log_rules;

-- name: CreateAuditRule :exec
CALL mysql.cloudsql_create_audit_rule(sqlc.arg(username), sqlc.arg(dbname), sqlc.arg(object), sqlc.arg(operation), sqlc.arg(op_result), 1, @outval, @outmsg);

-- name: ReadAuditRuleIDAfterCreate :one
SELECT id FROM audit_log_rules WHERE
	username = ? AND
	dbname = ? AND
	object = ? AND
	operation = ? AND
	op_result = ?;

-- name: ReadAuditLogRuleByID :one
SELECT * FROM audit_log_rules WHERE id = ?;

-- name: UpdatedAuditRuleByID :exec
CALL mysql.cloudsql_update_audit_rule(sqlc.arg(id), sqlc.arg(username), sqlc.arg(dbname), sqlc.arg(object), sqlc.arg(operation), sqlc.arg(op_result), 1, @outval, @outmsg);

-- name: DeleteAuditRuleByID :exec
CALL mysql.cloudsql_delete_audit_rule(sqlc.arg(id), 1, @outval, @outmsg);

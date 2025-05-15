-- name: GetAllAuditRules :many
SELECT * FROM audit_log_rules;

-- name: GetAuditRulesByUsername :many
SELECT * FROM audit_log_rules WHERE username = ?;

-- name: CreateAuditRule :exec
CALL mysql.cloudsql_create_audit_rule(sqlc.arg(username), sqlc.arg(dbname), sqlc.arg(object), sqlc.arg(operation), sqlc.arg(op_result), 1, @outval, @outmsg);

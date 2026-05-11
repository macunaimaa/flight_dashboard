package domain

import "time"

type AuditEntry struct {
	TenantID  string      `bson:"tenant_id" json:"tenantId"`
	UserID    string      `bson:"user_id" json:"userId"`
	Action    string      `bson:"action" json:"action"`
	Resource  string      `bson:"resource" json:"resource"`
	Detail    interface{} `bson:"detail,omitempty" json:"detail,omitempty"`
	IP        string      `bson:"ip" json:"ip"`
	RequestID string      `bson:"request_id" json:"requestId"`
	Timestamp time.Time   `bson:"timestamp" json:"timestamp"`
}

package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conditions struct {
	Field     string `json:"field"`
	Condition string `json:"condition"`
	Value     any    `json:"value"`
}

type Branches struct {
	Branch_name string            `json:"branch_name"`
	Conditions  []Conditions      `json:"conditions"`
	Actions     []ActionsWorkflow `json:"actions"`
}

type ActionsWorkflow struct {
	ID                   string                  `json:"id"`
	Type                 string                  `json:"type"`
	Subject              *string                 `json:"subject,omitempty"`
	Template             *string                 `json:"template,omitempty"`
	Time_offset          string                  `json:"time_offset"`
	Aggregation_template string                  `json:"aggregation_template"`
	Send_automatically   *bool                   `json:"send_automatically,omitempty"`
	Static_vars          *map[string]interface{} `json:"static_vars,omitempty"`
	Branches             []Branches              `json:"branches,omitempty"`
}

type Workflows struct {
	ID            primitive.ObjectID `bson:"_id"`
	Sequence_name string             `json:"sequence_name"`
	Description   string             `json:"description"`
	Actions       []ActionsWorkflow  `json:"actions"`
}

type Aggregation struct {
	Aggregation_name string      `json:"aggregation_name"`
	Collection       string      `json:"collection"`
	Aggregation      interface{} `json:"aggregation"`
}

type WorkflowHistory struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Status    string `json:"status"`
}

type WorkflowStatus struct {
	Workflow    string            `bson:"workflow"`
	Actions     []ActionsWorkflow `bson:"actions"`
	History     []WorkflowHistory `bson:"history"`
	Next_action string            `bson:"next_action"`
	Timestamp   any               `bson:"timestamp"`
}

type MessageNSQ struct {
	Name      string `json:"name"`
	Content   string `json:"content"`
	Timestamp string `json:"time_stamp"`
}

type ProcessParams struct {
	Process    ActionsWorkflow
	WorkflowId string
	StatusId   interface{}
	Exec       *bool
}

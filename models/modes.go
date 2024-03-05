package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Conditions struct {
	Field     string
	Condition string
	Value     any
}

type Branches struct {
	Branch_name string     `bson:"branch_name"`
	Conditions  []Conditions `bson:"conditions"`
	Actions     []EmailWorkflow    `bson:"actions"`
}

type EmailWorkflow struct {
	ID                   string                  `bson:"id"`
	Type                 string                  `bson:"type"`
	Subject              *string                 `bson:"subject,omitempty"`
	Template             *string                 `bson:"template,omitempty"`
	Time_offset          int                     `bson:"time_ofsset"`
	Aggregation_template string                  `bson:"aggregation_template"`
	Send_automatically   *bool                   `bson:"send_automatically,omitempty"`
	Static_vars          *map[string]interface{} `bson:"static_vars,omitempty"`
	Branches             []Branches                `json:"branches,omitempty"`
}

type Workflows struct {
	ID          primitive.ObjectID `bson:"_id"`
	Name        string             `bson:"sequence_name"`
	Description string             `bson:"description"`
	Actions     []EmailWorkflow    `bson:"actions"`
}

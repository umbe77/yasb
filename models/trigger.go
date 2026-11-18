package models

type Trigger struct {
	Endpoint string `json:"endpoint" bson:"endpoint"`
	Schema   string `json:"schema" bson:"schema"`
	Verb     string `json:"verb" bson:"verb"`
	Type     string `json:"type" bson:"type"`
}

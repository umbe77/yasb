package config

import "testing"

func TestGetConfig(t *testing.T) {
	t.Setenv("MONGO_DB_URI", "mongodb://test:27017")
	c := GetConfig()
	if c.MongoUri != "mongodb://test:27017" {
		t.Fatalf("MongoUri = %q", c.MongoUri)
	}
	if GetConfig() != c {
		t.Fatal("GetConfig should return the same instance")
	}
}

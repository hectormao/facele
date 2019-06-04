package ent

type Genericode struct {
	Genericode string `bson:"genericode" json:"genericode"`
	Code       string `bson:"code" json:"code"`
	Name       string `bson:"name" json:"name"`
}

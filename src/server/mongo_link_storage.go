package main

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

const DB_NAME = "shortly"

type Redirect struct {
	Slug string
	Url string
}

type MongoLinkStorage struct {
	session *mgo.Session
	redirects *mgo.Collection
}

func NewMongoLinkStorage(addr string) *MongoLinkStorage {
	s := MongoLinkStorage{}
	session, err := mgo.Dial(addr)
	if err != nil {
		panic(err)
	}
	s.session = session
	s.redirects = session.DB(DB_NAME).C("redirects")

	s.ensureIndexes()

	return &s
}

func (s *MongoLinkStorage) ensureIndexes() {
	index := mgo.Index{
		Key: []string{"slug"},
		Unique: true,
		Sparse: true,
		Background: false,
	}
	err := s.redirects.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

func (s *MongoLinkStorage) Store(slug string, url string) bool {
	err := s.redirects.Insert(&Redirect{Slug:slug, Url:url})
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func (s *MongoLinkStorage) GetUrl(slug string) string {
	result := Redirect{}
	err := s.redirects.Find(bson.M{"slug": slug}).One(&result)
	if err != nil {
		return ""
	}
	return result.Url
}

func (s *MongoLinkStorage) ClearAll() {
	err := s.redirects.DropCollection()
	s.session.ResetIndexCache()
	if err != nil {
		panic(err)
	}
	s.ensureIndexes()
}

func (s *MongoLinkStorage) Close() {
	s.Close()
}


package main

type Token struct {
	Id    string `gorethink:"id,omitempty"`
	Uid  string `gorethink:"uid"`
}

type User struct {
	Id string `gorethink:"id,omitempty"`
	Domain string `gorethink:"domain"`
	Tokens []string `gorethink:"tokens"`
	NamingScheme string `gorethink:"naming_scheme"`
	Encryption bool `gorethink:"encryption"`
}

type S3Bucket struct {
	Endpoint string `gorethink:"endpoint"`
	Bucket string `gorethink:"bucket"`
	AccessKeyId string `gorethink:"access_key_id"`
	SecretAccessKey string `gorethink:"secret_access_key"`
	Region string `gorethink:"region"`
}

type Domain struct {
	Id string `gorethink:"id,omitempty"`
	Public  bool `gorethink:"public"`
	Whitelist []string `gorethink:"whitelist"`
	Blacklist []string `gorethink:"blacklist"`
	Owner string `gorethink:"owner"`
	Bucket *S3Bucket `gorethink:"bucket"`
}

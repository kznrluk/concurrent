package core

import (
	"time"

	"github.com/lib/pq"
)

type Schema struct {
	ID  uint   `json:"id" gorm:"primaryKey;auto_increment"`
	URL string `json:"url" gorm:"type:text"`
}

type Key struct {
	ID              string    `json:"id" gorm:"primaryKey;type:char(42)"` // e.g. CK...
	Root            string    `json:"root" gorm:"type:char(42)"`
	Parent          string    `json:"parent" gorm:"type:char(42)"`
	EnactDocument   string    `json:"enactDocument" gorm:"type:json"`
	EnactSignature  string    `json:"enactSignature" gorm:"type:char(130)"`
	RevokeDocument  *string   `json:"revokeDocument" gorm:"type:json;default:null"`
	RevokeSignature *string   `json:"revokeSignature" gorm:"type:char(130);default:null"`
	ValidSince      time.Time `json:"validSince" gorm:"type:timestamp with time zone"`
	ValidUntil      time.Time `json:"validUntil" gorm:"type:timestamp with time zone"`
}

type SemanticID struct {
	ID        string    `json:"id" gorm:"primaryKey;type:text"`
	Owner     string    `json:"owner" gorm:"primaryKey;type:char(42)"`
	Target    string    `json:"target" gorm:"type:char(27)"`
	Document  string    `json:"document" gorm:"type:json"`
	Signature string    `json:"signature" gorm:"type:char(130)"`
	CDate     time.Time `json:"cdate" gorm:"->;<-:create;autoCreateTime"`
	MDate     time.Time `json:"mdate" gorm:"autoUpdateTime"`
}

// Association is one of a concurrent base object
// immutable
type Association struct {
	ID        string         `json:"id" gorm:"primaryKey;type:char(26)"`
	Author    string         `json:"author" gorm:"type:char(42);uniqueIndex:uniq_association"`
	Owner     string         `json:"owner" gorm:"type:char(42);"`
	SchemaID  uint           `json:"-" gorm:"uniqueIndex:uniq_association"`
	Schema    string         `json:"schema" gorm:"-"`
	Target    string         `json:"target" gorm:"type:char(27);uniqueIndex:uniq_association"`
	Variant   string         `json:"variant" gorm:"type:text;uniqueIndex:uniq_association"`
	Document  string         `json:"document" gorm:"type:json"`
	Signature string         `json:"signature" gorm:"type:char(130)"`
	CDate     time.Time      `json:"cdate" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
	Timelines pq.StringArray `json:"timelines" gorm:"type:text[]"`
}

// Profile is one of a Concurrent base object
// mutable
type Profile struct {
	ID           string        `json:"id" gorm:"primaryKey;type:char(26)"`
	Author       string        `json:"author" gorm:"type:char(42)"`
	SchemaID     uint          `json:"-"`
	Schema       string        `json:"schema" gorm:"-"`
	Document     string        `json:"document" gorm:"type:json"`
	Signature    string        `json:"signature" gorm:"type:char(130)"`
	Associations []Association `json:"associations,omitempty" gorm:"-"`
	CDate        time.Time     `json:"cdate" gorm:"->;<-:create;autoCreateTime"`
	MDate        time.Time     `json:"mdate" gorm:"autoUpdateTime"`
}

// Entity is one of a concurrent base object
// mutable
type Entity struct {
	ID                   string    `json:"ccid" gorm:"type:char(42)"`
	Domain               string    `json:"domain" gorm:"type:text"`
	Tag                  string    `json:"tag" gorm:"type:text;"`
	Score                int       `json:"score" gorm:"type:integer;default:0"`
	IsScoreFixed         bool      `json:"isScoreFixed" gorm:"type:boolean;default:false"`
	AffiliationDocument  string    `json:"affiliationDocument" gorm:"type:json"`
	AffiliationSignature string    `json:"affiliationSignature" gorm:"type:char(130)"`
	TombstoneDocument    *string   `json:"tombstoneDocument" gorm:"type:json;default:null"`
	TombstoneSignature   *string   `json:"tombstoneSignature" gorm:"type:char(130);default:null"`
	CDate                time.Time `json:"cdate" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
	MDate                time.Time `json:"mdate" gorm:"autoUpdateTime"`
}

type EntityMeta struct {
	ID        string  `json:"ccid" gorm:"type:char(42)"`
	Inviter   *string `json:"inviter" gorm:"type:char(42)"`
	Info      string  `json:"info" gorm:"type:json;default:'null'"`
	Signature string  `json:"signature" gorm:"type:char(130)"`
}

// Domain is one of a concurrent base object
// mutable
type Domain struct {
	ID           string      `json:"fqdn" gorm:"type:text"` // FQDN
	CCID         string      `json:"ccid" gorm:"type:char(42)"`
	Tag          string      `json:"tag" gorm:"type:text"`
	Score        int         `json:"score" gorm:"type:integer;default:0"`
	Meta         interface{} `json:"meta" gorm:"-"`
	IsScoreFixed bool        `json:"isScoreFixed" gorm:"type:boolean;default:false"`
	Dimension    string      `json:"dimension" gorm:"type:text"`
	CDate        time.Time   `json:"cdate" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
	MDate        time.Time   `json:"mdate" gorm:"autoUpdateTime"`
	LastScraped  time.Time   `json:"lastScraped" gorm:"type:timestamp with time zone"`
}

// Message is one of a concurrent base object
// immutable
type Message struct {
	ID              string         `json:"id" gorm:"primaryKey;type:char(26)"`
	Author          string         `json:"author" gorm:"type:char(42)"`
	SchemaID        uint           `json:"-"`
	Schema          string         `json:"schema" gorm:"-"`
	Document        string         `json:"document" gorm:"type:json"`
	Signature       string         `json:"signature" gorm:"type:char(130)"`
	CDate           time.Time      `json:"cdate" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
	Associations    []Association  `json:"associations,omitempty" gorm:"-"`
	OwnAssociations []Association  `json:"ownAssociations,omitempty" gorm:"-"`
	Timelines       pq.StringArray `json:"timelines" gorm:"type:text[]"`
}

// Timeline is one of a base object of concurrent
// mutable
type Timeline struct {
	ID          string    `json:"id" gorm:"primaryKey;type:char(26);"`
	Indexable   bool      `json:"indexable" gorm:"type:boolean;default:false"`
	Author      string    `json:"author" gorm:"type:char(42)"`
	DomainOwned bool      `json:"domainOwned" gorm:"type:boolean;default:false"`
	SchemaID    uint      `json:"-"`
	Schema      string    `json:"schema" gorm:"-"`
	Document    string    `json:"document" gorm:"type:json"`
	Signature   string    `json:"signature" gorm:"type:char(130)"`
	CDate       time.Time `json:"cdate" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
	MDate       time.Time `json:"mdate" gorm:"autoUpdateTime"`
}

// TimelineItem is one of a base object of concurrent
// immutable
type TimelineItem struct {
	ResourceID string    `json:"resourceID" gorm:"primaryKey;type:char(27);"`
	TimelineID string    `json:"timelineID" gorm:"primaryKey;type:char(26);"`
	Owner      string    `json:"owner" gorm:"type:char(42);"`
	Author     *string   `json:"author,omitempty" gorm:"type:char(42);"`
	CDate      time.Time `json:"cdate,omitempty" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
}

// CollectionItem is one of a base object of concurrent
// mutable
type CollectionItem struct {
	ID         string `json:"id" gorm:"primaryKey;type:char(26);"`
	Collection string `json:"collection" gorm:"type:char(26)"`
	Document   string `json:"document" gorm:"type:json"`
	Signature  string `json:"signature" gorm:"type:char(130)"`
}

type Ack struct {
	From      string `json:"from" gorm:"primaryKey;type:char(42)"`
	To        string `json:"to" gorm:"primaryKey;type:char(42)"`
	Document  string `json:"document" gorm:"type:json"`
	Signature string `json:"signature" gorm:"type:char(130)"`
	Valid     bool   `json:"valid" gorm:"type:boolean;default:false"`
}

// Subscription
type Subscription struct {
	ID          string             `json:"id" gorm:"primaryKey;type:char(26)"`
	Author      string             `json:"author" gorm:"type:char(42);"`
	Indexable   bool               `json:"indexable" gorm:"type:boolean;default:false"`
	DomainOwned bool               `json:"domainOwned" gorm:"type:boolean;default:false"`
	SchemaID    uint               `json:"-"`
	Schema      string             `json:"schema" gorm:"-"`
	Document    string             `json:"document" gorm:"type:json"`
	Signature   string             `json:"signature" gorm:"type:char(130)"`
	Items       []SubscriptionItem `json:"items" gorm:"foreignKey:Subscription"`
	CDate       time.Time          `json:"cdate" gorm:"->;<-:create;type:timestamp with time zone;not null;default:clock_timestamp()"`
	MDate       time.Time          `json:"mdate" gorm:"autoUpdateTime"`
}

type ResolverType uint

const (
	ResolverTypeEntity ResolverType = iota
	ResolverTypeDomain
)

type SubscriptionItem struct {
	ID           string       `json:"id" gorm:"primaryKey;type:text;"`
	Subscription string       `json:"subscription" gorm:"primaryKey;type:char(26)"`
	ResolverType ResolverType `json:"resolverType" gorm:"type:integer"`
	Entity       *string      `json:"entity" gorm:"type:char(42);"`
	Domain       *string      `json:"domain" gorm:"type:text;"`
}

// Event is websocket root packet model
type Event struct {
	Timeline  string       `json:"timeline"` // stream full id (ex: <streamID>@<domain>)
	Item      TimelineItem `json:"item,omitempty"`
	Resource  any          `json:"resource,omitempty"`
	Document  string       `json:"document"`
	Signature string       `json:"signature"`
}

type UserKV struct {
	Owner string `json:"owner" gorm:"primaryKey;type:char(42)"`
	Key   string `json:"key" gorm:"primaryKey;type:text"`
	Value string `json:"value" gorm:"type:text"`
}

type Chunk struct {
	Key   string         `json:"key"`
	Items []TimelineItem `json:"items"`
}
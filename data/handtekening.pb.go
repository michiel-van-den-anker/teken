// Code generated by protoc-gen-go.
// source: handtekening.proto
// DO NOT EDIT!

/*
Package data is a generated protocol buffer package.

It is generated from these files:
	handtekening.proto

It has these top-level messages:
	Handtekening
*/
package data

import proto "github.com/golang/protobuf/proto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal

type Handtekening struct {
	Voornaam        string `protobuf:"bytes,1,opt,name=voornaam" json:"voornaam,omitempty"`
	Tussenvoegsel   string `protobuf:"bytes,2,opt,name=tussenvoegsel" json:"tussenvoegsel,omitempty"`
	Achternaam      string `protobuf:"bytes,3,opt,name=achternaam" json:"achternaam,omitempty"`
	Geboortedatum   string `protobuf:"bytes,4,opt,name=geboortedatum" json:"geboortedatum,omitempty"`
	Geboorteplaats  string `protobuf:"bytes,5,opt,name=geboorteplaats" json:"geboorteplaats,omitempty"`
	Straat          string `protobuf:"bytes,6,opt,name=straat" json:"straat,omitempty"`
	Huisnummer      string `protobuf:"bytes,7,opt,name=huisnummer" json:"huisnummer,omitempty"`
	Postcode        string `protobuf:"bytes,8,opt,name=postcode" json:"postcode,omitempty"`
	Woonplaats      string `protobuf:"bytes,9,opt,name=woonplaats" json:"woonplaats,omitempty"`
	Handtekening    []byte `protobuf:"bytes,10,opt,name=handtekening,proto3" json:"handtekening,omitempty"`
	CaptchaResponse string `protobuf:"bytes,999,opt,name=captchaResponse" json:"captchaResponse,omitempty"`
}

func (m *Handtekening) Reset()         { *m = Handtekening{} }
func (m *Handtekening) String() string { return proto.CompactTextString(m) }
func (*Handtekening) ProtoMessage()    {}

func init() {
}

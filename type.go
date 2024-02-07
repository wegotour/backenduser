package backenduser

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Username string `bson:"username,omitempty" json:"username,omitempty"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
}

type Credential struct {
	Status  bool     `json:"status" bson:"status"`
	Token   string   `json:"token,omitempty" bson:"token,omitempty"`
	Message string   `json:"message,omitempty" bson:"message,omitempty"`
	Data    []Ticket `bson:"data,omitempty" json:"data,omitempty"`
	DataTransaksi   []Transaksi `bson:"datatransaksi,omitempty" json:"datatransaksi,omitempty"`
}

type Ticket struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama      string             `bson:"nama,omitempty" json:"nama,omitempty"`
	Harga     string             `bson:"harga,omitempty" json:"harga,omitempty"`
	Deskripsi string             `bson:"deskripsi,omitempty" json:"deskripsi,omitempty"`
	Stok      string             `bson:"stok,omitempty" json:"stok,omitempty"`
}

type Transaksi struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	NamaTicket  string             `bson:"namaticket,omitempty" json:"namaticket,omitempty"`
	Harga       string             `bson:"harga,omitempty" json:"harga,omitempty"`
	NamaPembeli string             `bson:"namapembeli,omitempty" json:"namapembeli,omitempty"`
	Email       string             `bson:"email,omitempty" json:"email,omitempty"`
	Alamat      string             `bson:"alamat,omitempty" json:"alamat,omitempty"`
	NoHP        string             `bson:"nohp,omitempty" json:"nohp,omitempty"`
	Quantity    string             `bson:"quantity,omitempty" json:"quantity,omitempty"`
	Total       string             `bson:"total,omitempty" json:"total,omitempty"`
}

type Pesan struct {
	ID     primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama   string             `bson:"nama,omitempty" json:"nama,omitempty"`
	Subjek string             `bson:"subjek,omitempty" json:"subjek,omitempty"`
	Pesan  string             `bson:"pesan,omitempty" json:"pesan,omitempty"`
}
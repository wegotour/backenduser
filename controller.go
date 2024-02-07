package backenduser

import (
	"context"
	"errors"
	"fmt"

	"github.com/aiteung/atdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func InsertOneDoc(db *mongo.Database, col string, docs interface{}) (insertedID primitive.ObjectID, err error) {
	cols := db.Collection(col)
	result, err := cols.InsertOne(context.Background(), docs)
	if err != nil {
		fmt.Printf("InsertOneDoc: %v\n", err)
	}
	insertedID = result.InsertedID.(primitive.ObjectID)
	return
}

func GetAllDocs(db *mongo.Database, col string, docs interface{}) interface{} {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &docs)
	if err != nil {
		fmt.Println(err)
	}
	return docs
}

func UpdateOneDoc(db *mongo.Database, col string, filter, update interface{}) (err error) {
	cols := db.Collection(col)
	result, err := cols.UpdateOne(context.Background(), filter, bson.M{"$set": update})
	if err != nil {
		fmt.Printf("UpdateOneDoc: %v\n", err)
	}
	if result.ModifiedCount == 0 {
		err = errors.New("no data has been changed with the specified filter")
		return err
	}
	return
}

func DeleteOneDoc(db *mongo.Database, col string, filter bson.M) (err error) {
	cols := db.Collection(col)
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		fmt.Printf("DeleteOneDoc: %v\n", err)
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("no data has been deleted with the specified filter")
		return err
	}
	return
}

// User
func InsertTicket(db *mongo.Database, col string, ticketdata Ticket) (insertedID primitive.ObjectID, err error) {
	insertedID, err = InsertOneDoc(db, col, ticketdata)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
	}
	return insertedID, err
}

func InsertTransaksi(db *mongo.Database, col string, transaksidata Transaksi) (insertedID primitive.ObjectID, err error) {
	objectid := primitive.NewObjectID()
	data := bson.M{
		"_id":         objectid,
		"namaticket":  transaksidata.NamaTicket,
		"harga":       transaksidata.Harga,
		"quantity":    transaksidata.Quantity,
		"total":       transaksidata.Total,
		"namapembeli": transaksidata.NamaPembeli,
		"email":       transaksidata.Email,
		"alamat":      transaksidata.Alamat,
		"nohp":        transaksidata.NoHP,
	}
	insertedID, err = InsertOneDoc(db, col, data)
	if err != nil {
		fmt.Printf("InsertUser: %v\n", err)
	}
	return insertedID, err
}

func GetAllDataTicket(db *mongo.Database, col string) (ticketlist []Ticket) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &ticketlist)
	if err != nil {
		fmt.Println(err)
	}
	return ticketlist
}

func GetAllDataTransaksi(db *mongo.Database, col string) (transaksilist []Transaksi) {
	cols := db.Collection(col)
	filter := bson.M{}
	cursor, err := cols.Find(context.TODO(), filter)
	if err != nil {
		fmt.Println("Error GetAllDocs in colection", col, ":", err)
	}
	err = cursor.All(context.TODO(), &transaksilist)
	if err != nil {
		fmt.Println(err)
	}
	return transaksilist
}

func InsertUser(db *mongo.Database, collection string, userdata User) string {
	hash, _ := HashPassword(userdata.Password)
	userdata.Password = hash
	atdb.InsertOneDoc(db, collection, userdata)
	return "Username : " + userdata.Username + "\nPassword : " + userdata.Password
}

func UpdateTicket(db *mongo.Database, col string, ticket Ticket) (tickets Ticket, status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": ticket.ID}
	update := bson.M{
		"$set": bson.M{
			"nama":      ticket.Nama,
			"harga":     ticket.Harga,
			"deskripsi": ticket.Deskripsi,
			"stok":      ticket.Stok,
		},
	}

	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return tickets, false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		err = fmt.Errorf("data tidak berhasil diupdate")
		return tickets, false, err
	}

	err = cols.FindOne(context.Background(), filter).Decode(&tickets)
	if err != nil {
		return tickets, false, err
	}

	return tickets, true, nil
}

func UpdateTransaksi(db *mongo.Database, col string, transaksi Transaksi) (transaksis Transaksi, status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": transaksi.ID}
	update := bson.M{
		"$set": bson.M{
			"namaticket"	:      transaksi.NamaTicket,
			"harga"			:      transaksi.Harga,
			"quantity"		: 	   transaksi.Quantity,
			"total"			:      transaksi.Total,
			"NamaPembeli"	:      transaksi.NamaPembeli,
			"email"			:      transaksi.Email,
			"alamat"		:      transaksi.Alamat,
			"nohp"			:      transaksi.NoHP,
		},
	}

	result, err := cols.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return transaksis, false, err
	}
	if result.ModifiedCount == 0 && result.UpsertedCount == 0 {
		err = fmt.Errorf("data tidak berhasil diupdate")
		return transaksis, false, err
	}

	err = cols.FindOne(context.Background(), filter).Decode(&transaksis)
	if err != nil {
		return transaksis, false, err
	}

	return transaksis, true, nil
}

func DeleteTicket(db *mongo.Database, col string, _id primitive.ObjectID) (status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("data tidak berhasil dihapus")
		return false, err
	}
	return true, nil
}

func DeleteTransaksi(db *mongo.Database, col string, _id primitive.ObjectID) (status bool, err error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}
	result, err := cols.DeleteOne(context.Background(), filter)
	if err != nil {
		return false, err
	}
	if result.DeletedCount == 0 {
		err = fmt.Errorf("data tidak berhasil dihapus")
		return false, err
	}
	return true, nil
}

func GetTicketFromID(db *mongo.Database, col string, _id primitive.ObjectID) (*Ticket, error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}

	ticketlist := new(Ticket)

	err := cols.FindOne(context.Background(), filter).Decode(ticketlist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no data found for ID %s", _id.Hex())
		}
		return nil, fmt.Errorf("error retrieving data for ID %s: %s", _id.Hex(), err.Error())
	}

	return ticketlist, nil
}

func GetTransaksiFromID(db *mongo.Database, col string, _id primitive.ObjectID) (*Transaksi, error) {
	cols := db.Collection(col)
	filter := bson.M{"_id": _id}

	transaksilist := new(Transaksi)

	err := cols.FindOne(context.Background(), filter).Decode(transaksilist)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("no data found for ID %s", _id.Hex())
		}
		return nil, fmt.Errorf("error retrieving data for ID %s: %s", _id.Hex(), err.Error())
	}

	return transaksilist, nil
}

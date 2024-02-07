package backenduser

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/whatsauth/watoken"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(PASETOPRIVATEKEYENV, MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	var Response Credential
	Response.Status = false
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	var datauser User
	err := json.NewDecoder(r.Body).Decode(&datauser)
	if err != nil {
		Response.Message = "error parsing application/json: " + err.Error()
	} else {
		if IsPasswordValid(mconn, collectionname, datauser) {
			Response.Status = true
			tokenstring, err := watoken.Encode(datauser.Username, os.Getenv(PASETOPRIVATEKEYENV))
			if err != nil {
				Response.Message = "Gagal Encode Token : " + err.Error()
			} else {
				Response.Message = "Selamat Datang"
				Response.Token = tokenstring
			}
		} else {
			Response.Message = "Password Salah"
		}
	}

	return GCFReturnStruct(Response)
}

func GCFReturnStruct(DataStuct any) string {
	jsondata, _ := json.Marshal(DataStuct)
	return string(jsondata)
}

// insert
func InsertDataTicket(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	ticketdata := new(Ticket)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&ticketdata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else if ticketdata.Nama == "" || ticketdata.Harga == "" || ticketdata.Deskripsi == "" || ticketdata.Stok == "" {
		resp.Message = "Data Tidak Boleh Kosong"
	} else {
		resp.Status = true
		insertedID, err := InsertTicket(conn, "ticket", *ticketdata)
		if err != nil {
			resp.Message = "Gagal memasukkan data ke database: " + err.Error()
		} else {
			resp.Message = "Berhasil Input data dengan ID: " + insertedID.Hex()
		}
	}
	return GCFReturnStruct(resp)
}

func InsertDataTransaksi(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	transaksidata := new(Transaksi)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&transaksidata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		insertedID, err := InsertTransaksi(conn, "transaksi", *transaksidata)
		if err != nil {
			resp.Message = "Gagal memasukkan data ke database: " + err.Error()
		} else {
			resp.Message = "Berhasil Input data dengan ID: " + insertedID.Hex()
		}
	}
	return GCFReturnStruct(resp)
}

func InsertDataPesan(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	pesandata := new(Pesan)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)
	err := json.NewDecoder(r.Body).Decode(&pesandata)
	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		insertedID, err := InsertPesanReview(conn, "review", *pesandata)
		if err != nil {
			resp.Message = "Gagal memasukkan data ke database: " + err.Error()
		} else {
			resp.Message = "Berhasil Input data dengan ID: " + insertedID.Hex()
		}
	}
	return GCFReturnStruct(resp)
}

// get
func GetAllData(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	data := GetAllDataTicket(mconn, collectionname)
	return GCFReturnStruct(data)
}

func GetOneDataTicket(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	ticketdata := new(Ticket)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&ticketdata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ticketdata.ID = ID

	// Menggunakan fungsi GetTicketFromID untuk mendapatkan data ticket berdasarkan ID
	ticketdata, err = GetTicketFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.Data = []Ticket{*ticketdata}

	return GCFReturnStruct(resp)
}

func GetDataTransaksi(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	data := GetAllDataTransaksi(mconn, collectionname)
	return GCFReturnStruct(data)
}

func GetOneDataTransaksi(MONGOCONNSTRINGENV, dbname, collectionname string, r *http.Request) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	resp := new(Credential)
	transaksidata := new(Transaksi)
	resp.Status = false
	err := json.NewDecoder(r.Body).Decode(&transaksidata)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	transaksidata.ID = ID

	// Menggunakan fungsi GetTicketFromID untuk mendapatkan data ticket berdasarkan ID
	transaksidata, err = GetTransaksiFromID(mconn, collectionname, ID)
	if err != nil {
		resp.Message = err.Error()
		return GCFReturnStruct(resp)
	}

	resp.Status = true
	resp.Message = "Get Data Berhasil"
	resp.DataTransaksi = []Transaksi{*transaksidata}

	return GCFReturnStruct(resp)
}

func GetDataReview(MONGOCONNSTRINGENV, dbname, collectionname string) string {
	mconn := SetConnection(MONGOCONNSTRINGENV, dbname)
	data := GetAllDataReview(mconn, collectionname)
	return GCFReturnStruct(data)
}


// update
func UpdateDataTicket(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	ticketdata := new(Ticket)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ticketdata.ID = ID

	err = json.NewDecoder(r.Body).Decode(&ticketdata)

	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		tickets, status, err := UpdateTicket(conn, "ticket", *ticketdata)
		if err != nil || !status {
			resp.Message = "Gagal update data : " + err.Error()
		} else {
			resp.Message = "Berhasil Update data dengan ID: " + tickets.ID.Hex()
		}
	}
	return GCFReturnStruct(resp)
}

func UpdateDataTransaksi(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	transaksidata := new(Transaksi)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	transaksidata.ID = ID

	err = json.NewDecoder(r.Body).Decode(&transaksidata)

	if err != nil {
		resp.Message = "error parsing application/json: " + err.Error()
	} else {
		resp.Status = true
		transaksis, status, err := UpdateTransaksi(conn, "transaksi", *transaksidata)
		if err != nil || !status {
			resp.Message = "Gagal update data : " + err.Error()
		} else {
			resp.Message = "Berhasil Update data dengan ID: " + transaksis.ID.Hex()
		}
	}
	return GCFReturnStruct(resp)
}

// delete
func DeleteDataTicket(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	ticketdata := new(Ticket)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ticketdata.ID = ID

	// err = json.NewDecoder(r.Body).Decode(&ticketdata)
	// if err != nil {
	// 	resp.Message = "error parsing application/json: " + err.Error()
	// } else {

	resp.Status = true
	status, err := DeleteTicket(conn, "ticket", ticketdata.ID)
	if err != nil || !status {
		resp.Message = "Gagal Delete data : " + err.Error()
	} else {
		resp.Message = "Berhasil Delete Data Ticket"
	}

	return GCFReturnStruct(resp)
}

func DeleteDataTransaksi(Mongoenv, dbname string, r *http.Request) string {
	resp := new(Credential)
	transaksidata := new(Transaksi)
	resp.Status = false
	conn := SetConnection(Mongoenv, dbname)

	id := r.URL.Query().Get("_id")
	if id == "" {
		resp.Message = "Missing '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	ID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		resp.Message = "Invalid '_id' parameter in the URL"
		return GCFReturnStruct(resp)
	}

	transaksidata.ID = ID

	// err = json.NewDecoder(r.Body).Decode(&transaksidata)
	// if err != nil {
	// 	resp.Message = "error parsing application/json: " + err.Error()
	// } else {

	resp.Status = true
	status, err := DeleteTransaksi(conn, "transaksi", transaksidata.ID)
	if err != nil || !status {
		resp.Message = "Gagal Delete data : " + err.Error()
	} else {
		resp.Message = "Berhasil Delete Data Transaksi"
	}

	return GCFReturnStruct(resp)
}
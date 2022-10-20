package postgres

import (
	"fmt"

	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateCustomer(req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	customerResp := pb.CustomerResponse{}
	fmt.Println(req)
	err := r.db.QueryRow(
		`insert into customer(firstname,lastname,bio,email,phonenumber) values($1,$2,$3,$4,$5)
			returning id,firstname,lastname,bio,email,phonenumber,created_at,updated_at
		`, req.Firstname, req.Lastname, req.Bio, req.Email, req.Phonenumber,
	).Scan(&customerResp.Id,&customerResp.FirstName, &customerResp.LastName, &customerResp.Bio,
		&customerResp.Email, &customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt)
	fmt.Println(err)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	for _,addres:=range req.Adderesses{
		addressResp:=pb.AddressResponse{}
		err=r.db.QueryRow(`
		insert into address(district,street) values($1,$2) returning id,district,street
		`,addres.District,addres.Street).Scan(&addressResp.Id,&addressResp.District,&addressResp.Street)
		if err != nil {
			return  &pb.CustomerResponse{},err
		}
		_,err=r.db.Exec(`
		insert into customer_address (customer_id,address_id) values($1,$2)`,customerResp.Id,addressResp.Id)
		if err != nil {
			return &pb.CustomerResponse{},err
		}
	customerResp.Adderesses=append(customerResp.Adderesses, &addressResp)
	}

	return &customerResp, nil

}

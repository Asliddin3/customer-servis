package postgres

import (
	"fmt"

	pb "github.com/Asliddin3/customer-servis/genproto/customer"

	"github.com/jmoiron/sqlx"
)

type customerRepo struct {
	db *sqlx.DB
}

func NewCustomerRepo(db *sqlx.DB) *customerRepo {
	return &customerRepo{db: db}
}

func (r *customerRepo) GetById(req *pb.CustomerId) (*pb.CustomerResponsePost, error) {
	var customerPostResp pb.CustomerResponsePost
	err := r.db.QueryRow(`
	select id,firstname,lastname,bio,email,phonenumber,created_at,updated_at from customer where id=$1
	`, req.Id).Scan(&customerPostResp.Id, &customerPostResp.FirstName, &customerPostResp.LastName, &customerPostResp.Bio,
		&customerPostResp.Email, &customerPostResp.PhoneNumber, &customerPostResp.CreatedAt, &customerPostResp.UpdatedAt,
	)
	if err != nil {
		return &pb.CustomerResponsePost{}, err
	}
	rows, err := r.db.Query(`
	select a.id,a.district,a.street from address a inner join customer_address ca
	on ca.customer_id=$1  where ca.address_id=a.id
	`, customerPostResp.Id)
	if err != nil {
		return &pb.CustomerResponsePost{}, err
	}
	for rows.Next() {
		addressResp := pb.AddressResponse{}
		err = rows.Scan(&addressResp.Id, &addressResp.District, &addressResp.Street)
		if err != nil {
			return &pb.CustomerResponsePost{}, err
		}
		customerPostResp.Adderesses = append(customerPostResp.Adderesses, &addressResp)
	}

	return &customerPostResp, nil
}
func (r *customerRepo) GetCustomerInfo(req *pb.CustomerId) (*pb.CustomerResponse, error) {
	customerInfo := pb.CustomerResponse{Id: req.Id}
	err := r.db.QueryRow(`
	select firstname,lastname,bio,email,phonenumber,created_at,updated_at from customer where id=$1
	`, customerInfo.Id).Scan(&customerInfo.FirstName, &customerInfo.LastName, &customerInfo.Bio,
		&customerInfo.Email, &customerInfo.PhoneNumber, &customerInfo.CreatedAt, &customerInfo.UpdatedAt)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	rows, err := r.db.Query(`
	select a.id,a.district,a.street from address a inner join customer_address ca on ca.address_id=a.id where ca.customer_id=$1
	`, customerInfo.Id)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	for rows.Next() {
		addressResp := pb.AddressResponse{}
		err = rows.Scan(&addressResp.Id, &addressResp.District, &addressResp.Street)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		customerInfo.Adderesses = append(customerInfo.Adderesses, &addressResp)
	}
	return &customerInfo, nil
}
func (r *customerRepo) DeleteCustomer(req *pb.CustomerId) (*pb.Empty, error) {
	_, err := r.db.Exec(`
	update customer set deleted_at=current_timestamp where id=$1
	`, req.Id)
	if err != nil {
		return &pb.Empty{}, err
	}
	return &pb.Empty{}, nil
}

func (r *customerRepo) GetListCustomers(req *pb.Empty) (*pb.ListCustomers, error) {
	CleanMap := func(mapOfFunc map[int]string) {
		for k := range mapOfFunc {
			delete(mapOfFunc, k)
		}
	}
	listCustomers := pb.ListCustomers{}
	rows, err := r.db.Query(`
	select id,deleted_at from customer where deleted_at is not null
	`)
	if err != nil {
		return &pb.ListCustomers{}, err
	}
	deletedCust := make(map[int]string)
	for rows.Next() {
		var id int
		var deleted_at string
		err = rows.Scan(&id, &deleted_at)
		if err != nil {
			return &pb.ListCustomers{}, err
		}
		deletedCust[id] = deleted_at
	}
	rows, err = r.db.Query(`
	select id,firstname,lastname,bio,email,phonenumber,created_at,updated_at from customer
	`)
	if err != nil {
		return &pb.ListCustomers{}, err
	}
	defer CleanMap(deletedCust)

	for rows.Next() {
		customerResp := pb.CustomerFullInfo{}
		err = rows.Scan(&customerResp.Id, &customerResp.FirstName,
			&customerResp.LastName, &customerResp.Bio, &customerResp.Email,
			&customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt)
		if err != nil {
			return &pb.ListCustomers{}, err
		}
		addreesResp := pb.AddressResponse{}
		err = r.db.QueryRow(`
		select a.id,a.district,a.street from address a inner join customer_address ca on ca.address_id=a.id where ca.customer_id=$1
		`, customerResp.Id).Scan(&addreesResp.Id, &addreesResp.District, &addreesResp.Street)
		if err != nil {
			return &pb.ListCustomers{}, err
		}
		customerResp.Adderesses = append(customerResp.Adderesses, &addreesResp)

		if val, ok := deletedCust[int(customerResp.Id)]; ok {
			customerResp.DeletedAt = val
			listCustomers.DeletedCustomers = append(listCustomers.DeletedCustomers, &customerResp)
		} else {
			listCustomers.ActiveCustomers = append(listCustomers.ActiveCustomers, &customerResp)
		}
	}

	return &listCustomers, nil
}

func (r *customerRepo) UpdateCustomer(req *pb.CustomerResponse) (*pb.CustomerResponse, error) {
	customerResp := pb.CustomerResponse{}
	err := r.db.QueryRow(`
	update customer set firstname=$1,lastname=$2,bio=$3,email=$4,phonenumber=$5,updated_at=current_timestamp
	where id=$6 returning id,firstname,lastname,bio,email,phonenumber,created_at,updated_at
	`, req.FirstName, req.LastName, req.Bio, req.Email, req.PhoneNumber, req.Id).Scan(
		&customerResp.Id, &customerResp.FirstName, &customerResp.LastName,
		&customerResp.Bio, &customerResp.Email, &customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt,
	)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	for _, address := range req.Adderesses {
		addressResp := pb.AddressResponse{}
		err = r.db.QueryRow(`
		update address set district=$1,street=$2 where id=$3
		returning id,district,street
		`, address.District, address.Street, address.Id).Scan(
			&addressResp.Id, &addressResp.District, &addressResp.Street,
		)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		customerResp.Adderesses = append(customerResp.Adderesses, &addressResp)
	}
	return &customerResp, nil
}

func (r *customerRepo) CreateCustomer(req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	customerResp := pb.CustomerResponse{}
	fmt.Println(req)
	err := r.db.QueryRow(
		`insert into customer(firstname,lastname,bio,email,phonenumber) values($1,$2,$3,$4,$5)
			returning id,firstname,lastname,bio,email,phonenumber,created_at,updated_at
		`, req.Firstname, req.Lastname, req.Bio, req.Email, req.Phonenumber,
	).Scan(&customerResp.Id, &customerResp.FirstName, &customerResp.LastName, &customerResp.Bio,
		&customerResp.Email, &customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt)
	fmt.Println(err)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	for _, addres := range req.Adderesses {
		addressResp := pb.AddressResponse{}
		err = r.db.QueryRow(`
		insert into address(district,street) values($1,$2) returning id,district,street
		`, addres.District, addres.Street).Scan(&addressResp.Id, &addressResp.District, &addressResp.Street)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		_, err = r.db.Exec(`
		insert into customer_address (customer_id,address_id) values($1,$2)`, customerResp.Id, addressResp.Id)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		customerResp.Adderesses = append(customerResp.Adderesses, &addressResp)
	}

	return &customerResp, nil

}

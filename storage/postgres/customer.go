package postgres

import (
	"fmt"
	"strings"

	pb "github.com/Asliddin3/customer-servis/genproto/customer"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
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
	select id,firstname,lastname,bio,email,phonenumber,created_at,updated_at from customer where id=$1 and deleted_at is null
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
	select firstname,lastname,bio,email,phonenumber,created_at,updated_at from customer where id=$1 and deleted_at is null
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
	rows, err := r.db.Query(`
	select id,firstname,username,lastname,bio,email,phonenumber,created_at,updated_at,deleted_at
	from customer
	`)
	if err != nil {
		return &pb.ListCustomers{}, err
	}
	customers := pb.ListCustomers{}
	var deletet_at interface{}
	for rows.Next() {
		customer := &pb.CustomerFullInfo{}
		err = rows.Scan(&customer.Id, &customer.FirstName, &customer.UserName, &customer.LastName,
			&customer.Bio, &customer.Email, &customer.PhoneNumber, &customer.CreatedAt, &customer.UpdatedAt, &deletet_at)
		if err != nil {
			return nil, err
		}
		data := fmt.Sprintf("%s", deletet_at)
		if strings.Contains(data, "nil") {
			customer.DeletedAt = ""
		} else {
			customer.DeletedAt = data
		}
		// switch v := deletet_at.(type) {
		// case nil:
		// 	customer.DeletedAt = ""
		// default:
		// 	customer.DeletedAt = fmt.Sprintf("%s", v)
		// }
		address, err := r.db.Query(`
		select a.id,a.district,a.street from address a inner join customer_address ca on ca.address_id=a.id where ca.customer_id=$1
		`, customer.Id)
		if err != nil {
			return &pb.ListCustomers{}, err
		}
		for address.Next() {
			addreesResp := pb.AddressResponse{}
			err = address.Scan(&addreesResp.Id, &addreesResp.District, &addreesResp.Street)
			if err != nil {
				return &pb.ListCustomers{}, err
			}
			customer.Adderesses = append(customer.Adderesses, &addreesResp)
		}
		customers.Customers = append(customers.Customers, customer)
	}
	// CleanMap := func(mapOfFunc map[string]string) {
	// 	for k := range mapOfFunc {
	// 		delete(mapOfFunc, k)
	// 	}
	// }
	// listCustomers := pb.ListCustomers{}
	// rows, err := r.db.Query(`
	// select id,deleted_at from customer where deleted_at is not null
	// `)
	// if err != nil {
	// 	return &pb.ListCustomers{}, err
	// }
	// deletedCust := make(map[string]string)
	// for rows.Next() {
	// 	var id string
	// 	var deleted_at string
	// 	err = rows.Scan(&id, &deleted_at)
	// 	if err != nil {
	// 		return &pb.ListCustomers{}, err
	// 	}
	// 	deletedCust[id] = deleted_at
	// }
	// rows, err = r.db.Query(`
	// select id,firstname,lastname,bio,email,phonenumber,created_at,updated_at from customer
	// `)
	// if err != nil {
	// 	return &pb.ListCustomers{}, err
	// }
	// defer CleanMap(deletedCust)

	// for rows.Next() {
	// 	customerResp := pb.CustomerFullInfo{}
	// 	err = rows.Scan(&customerResp.Id, &customerResp.FirstName,
	// 		&customerResp.LastName, &customerResp.Bio, &customerResp.Email,
	// 		&customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt)
	// 	if err != nil {
	// 		return &pb.ListCustomers{}, err
	// 	}
	// 	addreesResp := pb.AddressResponse{}
	// 	err = r.db.QueryRow(`
	// 	select a.id,a.district,a.street from address a inner join customer_address ca on ca.address_id=a.id where ca.customer_id=$1
	// 	`, customerResp.Id).Scan(&addreesResp.Id, &addreesResp.District, &addreesResp.Street)
	// 	if err != nil {
	// 		return &pb.ListCustomers{}, err
	// 	}
	// 	customerResp.Adderesses = append(customerResp.Adderesses, &addreesResp)

	// 	if val, ok := deletedCust[customerResp.Id]; ok {
	// 		customerResp.DeletedAt = val
	// 		listCustomers.DeletedCustomers = append(listCustomers.DeletedCustomers, &customerResp)
	// 	} else {
	// 		listCustomers.ActiveCustomers = append(listCustomers.ActiveCustomers, &customerResp)
	// 	}
	// }

	return &customers, nil
}

func (r *customerRepo) UpdateCustomer(req *pb.CustomerUpdate) (*pb.CustomerResponse, error) {
	customerResp := pb.CustomerResponse{}
	err := r.db.QueryRow(`
	update customer set firstname=$1,lastname=$2,bio=$3,email=$4,phonenumber=$5,updated_at=current_timestamp
	where id=$6 returning id,firstname,lastname,bio,email,phonenumber,created_at,updated_at
	`, req.FirstName, req.LastName, req.Bio, req.Email, req.PhoneNumber, req.Id).Scan(
		&customerResp.Id, &customerResp.FirstName, &customerResp.LastName,
		&customerResp.Bio, &customerResp.Email, &customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt,
	)
	fmt.Print("query error", err)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	if req.Adderesses != nil {
		addresses, err := r.db.Query(`
		select a.id from address a inner join customer_address ca
		on ca.customer_id=$1
		where a.id=ca.address_id
		`, customerResp.Id)
		fmt.Println("address query", err)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		for addresses.Next() {
			var id int
			err = addresses.Scan(&id)
			fmt.Println("error scan", err)
			if err != nil {
				return &pb.CustomerResponse{}, err
			}
			_, err = r.db.Exec(`
			delete from customer_address where customer_id=$1;
			`, customerResp.Id)
			if err != nil {
				return &pb.CustomerResponse{}, err
			}
			_, err = r.db.Exec(`
			delete from address where id =$1;
			`, id)
			if err != nil {
				return &pb.CustomerResponse{}, err
			}
			fmt.Println("error delete", err)

		}
		for _, address := range req.Adderesses {
			addressResp := pb.AddressResponse{}
			err = r.db.QueryRow(`
			insert into address (district,street) values($1,$2) returning id,district ,street
			`, address.District, address.Street).Scan(
				&addressResp.Id, &addressResp.District, &addressResp.Street,
			)
			fmt.Println("insertin address", err)
			if err != nil {
				return &pb.CustomerResponse{}, err
			}
			customerResp.Adderesses = append(customerResp.Adderesses, &addressResp)
		}
	}
	fmt.Println(customerResp)
	return &customerResp, nil
}

func (r *customerRepo) CheckField(req *pb.CheckFieldRequest) (*pb.CheckFieldResponse, error) {
	key := req.Key
	if key == "email" {
		var exists int
		err := r.db.QueryRow(`
		select count(*) from customer where email=$1 and deleted_at is null
		`, req.Value).Scan(&exists)
		if err != nil {
			return &pb.CheckFieldResponse{}, err
		}
		if exists != 0 {
			return &pb.CheckFieldResponse{Exists: true}, nil
		} else {
			return &pb.CheckFieldResponse{Exists: false}, nil
		}
	} else if key == "username" {
		var exists int
		err := r.db.QueryRow(`
		select count(*) from customer where username=$1 and deleted_at is null
		`, req.Value).Scan(&exists)
		if err != nil {
			return &pb.CheckFieldResponse{}, err
		}
		if exists != 0 {
			return &pb.CheckFieldResponse{Exists: true}, nil
		} else {
			return &pb.CheckFieldResponse{Exists: false}, nil
		}
	}
	return &pb.CheckFieldResponse{Exists: false}, nil
}

func (r *customerRepo) CreateCustomer(req *pb.CustomerRequest) (*pb.CustomerResponse, error) {
	customerResp := &pb.CustomerResponse{}
	fmt.Println(req.PassWord)
	err := r.db.QueryRow(
		`insert into customer(firstname,lastname,bio,email,phonenumber,username,password,
			refresh_token) values($1,$2,$3,$4,$5,$6,$7,$8)
			returning id,firstname,lastname,bio,email,phonenumber,created_at,updated_at
		`, req.FirstName, req.LastName, req.Bio, req.Email, req.PhoneNumber, req.UserName, req.PassWord,
		req.RefreshToken,
	).Scan(&customerResp.Id, &customerResp.FirstName, &customerResp.LastName, &customerResp.Bio,
		&customerResp.Email, &customerResp.PhoneNumber, &customerResp.CreatedAt, &customerResp.UpdatedAt)
	fmt.Println(err)
	fmt.Println(customerResp)
	if err != nil {
		return &pb.CustomerResponse{}, err
	}
	for _, addres := range req.Adderesses {
		addressResp := &pb.AddressResponse{}
		err = r.db.QueryRow(`
		insert into address(district,street) values($1,$2) returning id,district,street
		`, addres.District, addres.Street).Scan(&addressResp.Id, &addressResp.District, &addressResp.Street)
		fmt.Println("is that here", err)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		_, err = r.db.Exec(`
		insert into customer_address (customer_id,address_id) values($1,$2)`, customerResp.Id, addressResp.Id)
		if err != nil {
			return &pb.CustomerResponse{}, err
		}
		customerResp.Adderesses = append(customerResp.Adderesses, addressResp)
	}

	return customerResp, nil

}

func (r *customerRepo) Login(req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginResponse := pb.LoginResponse{}
	fmt.Println(req.UserName)
	err := r.db.QueryRow(`
	select id,firstname,lastname,email,bio,refresh_token,password from customer
	where username=$1 and deleted_at is null
	`, req.UserName).Scan(&loginResponse.Id, &loginResponse.FirstName,
		&loginResponse.LastName, &loginResponse.Email, &loginResponse.Bio,
		&loginResponse.RefreshToken, &loginResponse.PassWord)
	if err != nil {
		return &pb.LoginResponse{}, err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(loginResponse.PassWord), []byte(req.Password)); err != nil {
		return &pb.LoginResponse{}, err
	}
	return &loginResponse, nil
}

func (s *customerRepo) RefreshToken(req *pb.RefreshTokenRequest) (*pb.LoginResponse, error) {
	customerResp := pb.LoginResponse{}
	err := s.db.QueryRow(`
	update customer set refresh_token=$1 where id=$2
	returning id,firstname,lastname,bio,email,password,
	refresh_token
	`, req.RefreshToken, req.Id).Scan(
		&customerResp.Id, &customerResp.FirstName, &customerResp.LastName,
		&customerResp.Bio, &customerResp.Email, &customerResp.PassWord,
		&customerResp.RefreshToken,
	)
	if err != nil {
		return &pb.LoginResponse{}, err
	}
	return &customerResp, nil
}

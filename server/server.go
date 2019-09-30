package main

import (
	"context"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	pb "github.com/shrikar007/02-crud-grpc/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

type server struct {
	Db *gorm.DB
}

type Expenseitem struct {
	ID       string               `bson:"id"`
	Description string             `bson:"description"`
	Typeofaccount  string             `bson:"typeofaccount"`
	Amount    string             `bson:"amount"`
}
func dataToexpensePb(data *Expenseitem) *pb.Expense {
	return &pb.Expense{
		Id:       data.ID,
		Description:data.Description,
		Typeofaccount:data.Typeofaccount,
		Amount:data.Amount,
	}
}
func (db *server) CreateExpense(ctx context.Context, req *pb.CreateReq) (*pb.CreateRes, error) {
	fmt.Println("Create Expense request")
	expense := *req.GetExp()

	data := Expenseitem{
		Description: expense.GetDescription(),
		Typeofaccount:   expense.GetTypeofaccount(),
		Amount:  expense.GetAmount(),
	}
	fmt.Println(data)
    db1:=db.Db.Create(&data)

	if db1.RowsAffected != 0 {
		return &pb.CreateRes{
			Exp: &pb.Expense{
				Id:       expense.Id,
				Description: expense.GetDescription(),
				Typeofaccount:    expense.GetTypeofaccount(),
				Amount:  expense.GetAmount(),
			},

		}, nil

	} else {
		err := errors.New("Unable to update")
       return nil,err
	}
}

func (db *server) ReadExpense(ctx context.Context, req *pb.ReadReq) (*pb.ReadRes, error) {
	var expense Expenseitem
	fmt.Println("Read blog request")

	expenseId := req.GetId()
	db1:= db.Db.Table("expenseitems").Where("id = ?", expenseId).Find(&expense)
	if db1.RowsAffected != 0 {
		return &pb.ReadRes{
			Exp: &pb.Expense{
				Id:       expense.ID,
				Description: expense.Description,
				Typeofaccount:    expense.Typeofaccount,
				Amount:  expense.Amount,
			},

		}, nil

	} else {
		err := errors.New("Unable to read")
		return nil,err
	}
}

func (db *server) UpdateExpense(ctx context.Context, req *pb.UpdateReq) (*pb.UpdateRes, error) {
	var temp Expenseitem
	expense:=req.GetExp()
	db1:= db.Db.Table("expenseitems").Where("id = ?", expense.Id).Find(&temp)
	if db1.RowsAffected != 0 {
		temp.Description = expense.Description
		temp.Typeofaccount = expense.Typeofaccount
		temp.Amount = expense.Amount
		db2:=db1.Update(&temp)
		if db2.RowsAffected!=0{
			return &pb.UpdateRes{
				Exp:&pb.Expense{
					Description:temp.Description,
					Typeofaccount:temp.Typeofaccount,
					Amount:temp.Amount,
				},
			},nil
		}else{
			err := errors.New("Unable to Update")
			return nil,err

		}

	}else{
		err := errors.New("Unable to find id")
		return nil,err
	}
}
func (db *server) DeleteExpense(ctx context.Context, req *pb.DeleteReq) (*pb.DeleteRes, error) {
	var temp Expenseitem
	expenseid:=req.GetId()

	db1:= db.Db.Table("expenseitems").Where("id = ?", expenseid).Find(&temp)
	if db1.RowsAffected != 0 {
		db2:=db1.Delete(&temp)
		if db2.RowsAffected!=0{
			return &pb.DeleteRes{},nil
		}else{
			err := errors.New("Unable to delete")
			return nil,err

		}

	}else{
		err := errors.New("Unable to find id")
		return nil,err
	}
}

func ( db *server) ListExpenses(req *pb.ListReq, stream pb.ExpenseService_ListExpensesServer) error {
var exp Expenseitem
	db1,_ := db.Db.Model(&Expenseitem{}).Rows()

	for db1.Next(){

		db.Db.ScanRows(db1,&exp)
		stream.Send(&pb.ListRes{Exp:dataToexpensePb(&exp)})
	}
	if err:=db1.Err() ; err!=nil{
		return err

	}
	return nil
}

func main(){
	db, err := gorm.Open("mysql", "root:root@tcp(localhost:3306)/grpccrud?charset=utf8&parseTime=True")

	if err != nil {
		fmt.Println(err)
	}
	lis, err := net.Listen("tcp", "0.0.0.0:5005")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	pb.RegisterExpenseServiceServer(s, &server{Db:db})
	reflection.Register(s)

	go func() {
		fmt.Println("Starting Expense Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
}
package main

import (
	"context"
	"fmt"
	pb "github.com/shrikar007/02-crud-grpc/protos"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
)

func main(){
	var ch int
	fmt.Println("Expense Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:5005", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := pb.NewExpenseServiceClient(cc)


	for {


	fmt.Println("1.Create:\n2.Read:\n3:Update:\n4.Delete:\n5.Listall:\n6:Enter any key to exit:")
	fmt.Println("enter choice:")
	fmt.Scan(&ch)

	switch ch {
	case 1:
		var desc,toa,amount string
		fmt.Println("Enter Description,Type of account,Amount")
		fmt.Scan(&desc,&toa,&amount)
		expens := &pb.Expense{
			Description:   desc,
			Typeofaccount: toa,
			Amount:        amount,
		}
		createRes, err := c.CreateExpense(context.Background(), &pb.CreateReq{Exp: expens})
		//fmt.Println(createRes)
		if err != nil {
			log.Fatalf("Unexpected error: %v", err)
		}
		fmt.Printf("Expense has been created: %v", createRes)
	case 2:
		var id string
		fmt.Println("Enter id to read")
		fmt.Scan(&id)
		readexpenseReq := &pb.ReadReq{Id: id}
		readexpenseRes, readexpErr := c.ReadExpense(context.Background(), readexpenseReq)
		if readexpErr != nil {
			fmt.Printf("Error happened while reading: %v \n", readexpErr)
		}

		fmt.Printf("expense was read: %v \n", readexpenseRes)
	case 3:
		var id,desc,toa,amount string
		fmt.Println("Enter Id to be update")
		fmt.Scan(&id)
		fmt.Println("Enter Description,Type of account,Amount")
		fmt.Scan(&desc,&toa,&amount)

		newExpense := &pb.Expense{
			Id:            id,
			Description:   desc,
			Typeofaccount: toa,
			Amount:        amount,
		}
		updateRes, updateErr := c.UpdateExpense(context.Background(), &pb.UpdateReq{Exp: newExpense})
		if updateErr != nil {
			fmt.Printf("Error happened while updating: %v \n", updateErr)
		}
		fmt.Printf("expense was updated: %v\n", updateRes)
	case 4:
		var id string
		fmt.Println("Enter Id to delete")
		fmt.Scan(&id)
		deleteRes, deleteErr := c.DeleteExpense(context.Background(), &pb.DeleteReq{Id: id})

		if deleteErr != nil {
			fmt.Printf("Error happened while deleting: %v \n", deleteErr)
		}
		fmt.Printf("expense was deleted: %v \n", deleteRes)
	case 5:
		stream, err := c.ListExpenses(context.Background(), &pb.ListReq{})
		if err != nil {
			log.Fatalf("error while calling Listexpense RPC: %v", err)
		}
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Something happened: %v", err)
			}
			fmt.Println(res.GetExp())
		}
	default:
		os.Exit(3)
	 }
	}
}

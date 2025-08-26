package main

import (
	"kafkapractisel0/consumer"
	"kafkapractisel0/database"
	"kafkapractisel0/emulator"
	"kafkapractisel0/handler"
	"kafkapractisel0/mock"
	"kafkapractisel0/repo"
	"kafkapractisel0/repo/cache"
	"kafkapractisel0/services"
	"log"
	"net/http"
	"time"
)

func main() {
	envService := services.NewEnv()

	db := database.InitDB(envService)
	defer db.Close()

	// go init cache

	// repos
	orderRepo := repo.NewOrderRepo(db)
	itemsRepo := repo.NewItemsRepo(db)
	customerRepo := repo.NewCustomerRepo(db)
	deliveryRepo := repo.NewDeliveryRepo(db)
	paymentRepo := repo.NewPaymentRepo(db)
	cacheRepo := cache.NewCacheRepo(1000)

	// services
	orderService := services.NewOrderService(db, customerRepo, deliveryRepo, itemsRepo, orderRepo, paymentRepo)
	emulatorService := services.NewEmulatorService(customerRepo, deliveryRepo, itemsRepo)
	cacheService := services.NewCacheService(cacheRepo, orderRepo)

	// generate mock data
	go mock.GenerateMockCustomer(customerRepo)
	go mock.GenerateMockItem(itemsRepo)
	go mock.GenerateMockDelivery(deliveryRepo)

	// emulator and consumer
	go emulator.StartEmulate(envService, emulatorService)

	myConsumer := consumer.NewConsumer(envService, orderService, cacheRepo)
	go myConsumer.StartConsume()

	// cache update
	go cacheService.UpdateCacheNewest(10000)

	// http server
	mux := http.NewServeMux()
	orderHandler := handler.NewOrderHandler(orderService, cacheRepo)
	mux.Handle("/order/{order_uid}", middlewareOne(http.HandlerFunc(orderHandler.GetOrderById)))
	log.Println("Server is starting on :8081 ...")
	err := http.ListenAndServe(":8081", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		next.ServeHTTP(w, r)
		elapsedTime := time.Since(startTime)
		log.Printf("[%s] [%s] [%s]\n", r.Method, r.URL.Path, elapsedTime)
	})
}

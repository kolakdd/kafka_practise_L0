package mock

import (
	"kafkapractisel0/models"
	"kafkapractisel0/repo"
	"log"
	"math/rand"

	"github.com/go-faker/faker/v4"
)

func GenerateMockCustomer(s repo.CustomerRepo) {
	for range 10 {
		if err := s.CreateCustomer(faker.FirstName()); err != nil {
			log.Printf("error while generate customer, %v", err)
		}
	}
}

func GenerateMockItem(r repo.ItemsRepo) {
	for range 10 {
		var item models.Item
		item.Price = int64(rand.Intn(900_00) + 100_00)
		item.TrackNumber = "XXXLMTESTTRACK_UNIVERSAL"
		item.Rid = RidGenerator()
		item.Name = "Палочки для еды от " + faker.ChineseFirstName()
		item.Sale = rand.Intn(90) + 10
		item.Size = int16(rand.Intn(5) + 0)
		item.Currency = models.CurrencyRUB
		item.TotalPrice = item.Price - item.Price/100*int64(rand.Intn(15)+0) // скидка от 0 до 15 процентов от макс. цены
		item.NmID = 111_111
		item.Brand = "Палочки Из Поднебесной"
		item.Status = 202
		if err := r.CreateItem(item); err != nil {
			log.Printf("error while generate item, %v", err)
		}
	}
}

func GenerateMockDelivery(r repo.DeliveryRepo) {
	for range 10 {
		var delivery models.Delivery
		delivery.Name = faker.Name() + " " + faker.LastName()
		delivery.Phone = "+9720000000"
		delivery.Zip = "2639809"
		delivery.City = "Gotham City"
		delivery.Address = "Joy Ker Street"
		delivery.Region = "EnJoyKer region"
		delivery.Email = "joyker@gotham.badman"
		if err := r.CreateDelivery(delivery); err != nil {
			log.Printf("error while generate delivery, %v", err)
		}
	}

}

func RidGenerator() string {
	const hexChars = "0123456789abcdef"
	randomPart := make([]byte, 19)
	for i := range randomPart {
		randomPart[i] = hexChars[rand.Intn(len(hexChars))]
	}
	return string(randomPart) + "test"
}

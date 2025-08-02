package recurse

import "time"

type Person struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type Checkin struct {
	Person    Person    `json:"person"`
	CreatedAt time.Time `json:"created_at"`
}

type Profile struct {
	Stints []Stint `json:"stints"`
}

type Stint struct {
	Batch Batch `json:"batch"`
}

type Batch struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"short_name"`
}

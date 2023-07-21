package main

import (
	"fmt"
	"time"
)

// These structs intentionally empty for the purpose of this example.
// They are simply marker classes for the purpose of demonstration, contain no properties, and serve no other purpose.
type Bacon struct{ fried bool }
type Coffee struct{ poured bool }
type Egg struct{ fried bool }
type Juice struct{ poured bool }
type Toast struct{ toasted, butter, jam bool }

func pourOJ() *Juice {
	fmt.Println("Pouring orange juice")
	time.Sleep(time.Second)
	return &Juice{poured: true}
}

func applyJam(toast *Toast) {
	fmt.Println("Putting jam on the toast")
	time.Sleep(time.Second)
	toast.jam = true
}

func applyButter(toast *Toast) {
	fmt.Println("Putting butter on the toast")
	time.Sleep(time.Second)
	toast.butter = true
}

func toastBread(slices int) *Toast {
	for slice := 0; slice < slices; slice++ {
		fmt.Println("Putting a slices of bread in the toaster")
	}

	fmt.Println("Start toasting...")
	time.Sleep(time.Second * 3)

	return &Toast{toasted: true}
}

func fryBacon(slices int) *Bacon {
	fmt.Printf("Putting %d slices of bacon in the pan\n", slices)
	fmt.Println("Cooking first side of bacon...")
	time.Sleep(3 * time.Second)
	for slice := 0; slice < slices; slice++ {
		fmt.Println("Flipping a slice of bacon")
	}

	fmt.Println("Cooking the second side of bacon...")
	time.Sleep(3 * time.Second)
	fmt.Println("Put bacon on plate")

	return &Bacon{fried: true}
}

func fryEggs(howMany int) *Egg {
	fmt.Println("Warming the egg pan...")
	time.Sleep(3 * time.Second)
	fmt.Printf("Cracking %d eggs\n", howMany)
	time.Sleep(3 * time.Second)
	fmt.Println("Put eggs on plate")

	return &Egg{fried: true}
}

func pourCoffee() *Coffee {
	fmt.Println("Pouring coffee")
	time.Sleep(time.Second)
	return &Coffee{poured: true}
}

func main() {
	start := time.Now()

	go func() {
		cup := pourCoffee()
		fmt.Println("* Coffee is ready:", cup.poured)
	}()

	go func() {
		egg := fryEggs(2)
		fmt.Println("* Eggs are ready:", egg.fried)
	}()

	go func() {
		bacon := fryBacon(3)
		fmt.Println("* Bacon is ready:", bacon.fried)
	}()

	go func() {
		toast := toastBread(2)
		applyButter(toast)
		applyJam(toast)
		fmt.Println("* Toast is ready:", toast.toasted && toast.butter && toast.jam)
	}()

	go func() {
		oj := pourOJ()
		fmt.Println("* Orange juice is ready:", oj.poured)

	}()

	fmt.Printf("* Breakfast is ready, it took %s\n", time.Since(start))
}

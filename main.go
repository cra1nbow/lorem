package main

import (
	crand "crypto/rand"
	"fmt"
	"io"
	"log"
	"math"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
)

func generateUrl(x, y int) string {
	return fmt.Sprintf("https://picsum.photos/%d/%d", x, y)
}

func generateRandomSize() int {
	return rand.Int()%1000 + 200
}

func generateRandomUrl() string {
	return generateUrl(generateRandomSize(), generateRandomSize())
}

func init() {
	seed, _ := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	rand.Seed(seed.Int64())
}

func downloadRandomPicsum(wg *sync.WaitGroup) {
	defer wg.Done()
	res, err := http.Get(generateRandomUrl())
	if err != nil {
		log.Fatalln("fail request")
	}
	defer res.Body.Close()

	file, err := os.Create(fmt.Sprintf("./lorempicsum/%s.jpg", strings.Split(res.Request.URL.Path, "/")[2]))
	if err != nil {
		log.Fatalln("fail create file")
	}
	defer file.Close()
	io.Copy(file, res.Body)
}

func main() {
	os.Mkdir("lorempicsum", 0755)

	var wg sync.WaitGroup

	for range [100]int{} {
		wg.Add(1)
		go downloadRandomPicsum(&wg)
	}
	wg.Wait()
}

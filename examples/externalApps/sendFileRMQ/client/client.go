package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gfads/midarch/pkg/shared"
	"github.com/gfads/midarch/pkg/shared/lib"

	"github.com/streadway/amqp"
)

func clientX(ch *amqp.Channel, q amqp.Queue, msgs <-chan amqp.Delivery) {
	var FILE_SIZE string
	var SAMPLE_SIZE, AVERAGE_WAITING_TIME int
	if len(os.Args) >= 2 {
		FILE_SIZE = os.Args[1]
		SAMPLE_SIZE, _ = strconv.Atoi(os.Args[2])
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(os.Args[3])
	} else {
		FILE_SIZE = shared.EnvironmentVariableValueWithDefault("FILE_SIZE", "md")
		SAMPLE_SIZE, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("SAMPLE_SIZE", "100"))
		AVERAGE_WAITING_TIME, _ = strconv.Atoi(shared.EnvironmentVariableValueWithDefault("AVERAGE_WAITING_TIME", "60"))
	}
	fmt.Println("dateTime;info;sequential;response_time") //"FILE_SIZE / SAMPLE_SIZE / AVERAGE_WAITING_TIME:", FILE_SIZE, "/", SAMPLE_SIZE, "/", AVERAGE_WAITING_TIME)

	fileBytes := getFile(FILE_SIZE)
	rand.Seed(time.Now().UnixNano())
	// invoke remote method
	for x := 0; x < SAMPLE_SIZE; x++ {

		t1 := time.Now()

		// Encode the []byte to a base64 string
		base64String := base64.StdEncoding.EncodeToString(fileBytes)

		_, err := fibonacciRPC(ch, q, msgs, base64String)
		// failOnError(err, "Failed to handle RPC request")
		if err != nil {
			log.Fatal(";error", err, ";;\n")
		}
		//fmt.Printf("Fibo: %d => %d\n", n, res)
		t2 := time.Now()

		duration := t2.Sub(t1)
		//time.Sleep(3*time.Second)

		//durations[i] = t2.Sub(t1)

		// fmt.Printf("%v\n", float64(duration.Nanoseconds())/1000000)
		log.Printf(";ok;%d;%f\n", x+1, float64(duration.Nanoseconds())/1000000)

		// Normally distributed waiting time between calls with an average of 60 milliseconds and standard deviation of 20 milliseconds
		var rd = int(math.Round((rand.NormFloat64() + 3) * float64(AVERAGE_WAITING_TIME/3)))
		if rd > 0 {
			time.Sleep(time.Duration(rd) * time.Millisecond)
		}
	}

	//totalTime := time.Duration(0)
	//for i := range durations {
	//	totalTime += durations[i]
	//}

	//fmt.Printf("Tempo Total [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE, totalTime)
	//fmt.Printf("Tempo Médio [N=%v] [SAMPLE=%v] [TIME=%v]\n", N, shared.SAMPLE_SIZE, totalTime/shared.SAMPLE_SIZE)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func fibonacciRPC(ch *amqp.Channel, q amqp.Queue, msgs <-chan amqp.Delivery, base64File string) (res bool, err error) {
	corrId := randomString(32)

	err = ch.Publish(
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(base64File),
		})
	failOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res, err = strconv.ParseBool(string(d.Body))
			failOnError(err, "Failed to convert body to bool")
			break
		}
	}

	return
}

func main() {
	// Wait for RabbitMQ and server to get up
	timeToRun, _ := strconv.Atoi(shared.EnvironmentVariableValueWithDefault("TIME_TO_START_CLIENT", "13"))
	lib.PrintlnDebug("Waiting", timeToRun, "seconds for RabbitMQ and server to get up")
	time.Sleep(time.Duration(timeToRun) * time.Second)

	rand.Seed(time.Now().UTC().UnixNano())

	conn, err := amqp.Dial("amqp://guest:guest@rmq:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	clientX(ch, q, msgs)
}

func timeTrack(start time.Time, name string) time.Duration {
	elapsed := time.Since(start)
	return elapsed
}

//func bodyFrom(args []string) int {
//	var s string
//	if (len(args) < 2) || os.Args[1] == "" {
//		s = "30"
//	} else {
//		s = strings.Join(args[1:], " ")
//	}
//	n, err := strconv.Atoi(s)
//	failOnError(err, "Failed to convert arg to integer")
//	return n
//}

func getFile(size string) []byte {
	var fileName string
	switch size {
	case "sm":
		fileName = shared.DIR_BASE + "/examples/sendfiledistributed/client/36x36.png"
	case "md":
		fileName = shared.DIR_BASE + "/examples/sendfiledistributed/client/2k.png"
	case "lg":
		fileName = shared.DIR_BASE + "/examples/sendfiledistributed/client/4k.jpg" // Foto de <a href="https://unsplash.com/pt-br/@francesco_ungaro?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Francesco Ungaro</a> na <a href="https://unsplash.com/pt-br/fotografias/cardume-de-peixes-no-corpo-de-agua-MJ1Q7hHeGlA?utm_content=creditCopyText&utm_medium=referral&utm_source=unsplash">Unsplash</a>
	}

	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		shared.ErrorHandler(shared.GetFunction(), "Error opening file: "+fileName)
	}

	return fileBytes
}

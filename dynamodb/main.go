package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/samstradling/dynamodb-lock-client-golang"
	"os"
	"time"

	"os/exec"
	
)

func acquireLock() *lockclient.DynamoDBLockClient {
	tableName := os.Args[1]

	config := &aws.Config{
		Region: aws.String("ap-southeast-2"),
	}
	sess := session.Must(session.NewSession(config))

	lockClient := &lockclient.DynamoDBLockClient{
		LockName:        "jenkins",
		LeaseDuration:   60000 * time.Millisecond,
		HeartbeatPeriod: 1000 * time.Millisecond,
		TableName:       tableName,
		Client:          dynamodb.New(sess),
	}

	for true {
		result, err := lockClient.GetLock()
		if result {
			break
		}
		fmt.Printf("Unabe to get lock: %v\n", err)
		time.Sleep(10 * time.Second)
	}

	return lockClient
}

func main() {

	if len(os.Args) == 1 {
		fmt.Printf("Usage: lockandexec <dynamodb-table-name> <prog args..>\n")
		os.Exit(1)
	}

	lockClient := acquireLock()
	defer lockClient.RemoveLock()
	
	fmt.Printf("Acquired lock..\n")

	// do the work
	c := exec.Command(os.Args[2])
	for _, a := range os.Args[3:] {
		c.Args = append(c.Args, a)
	}

	outerr, err := c.CombinedOutput()
	if err != nil {
		fmt.Print(string(outerr))
		fmt.Print(err)
		os.Exit(1)
	} else {
		fmt.Printf("%s", outerr)
	}	
}

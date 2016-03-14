package main

/*
Schema:
create table example.samples (
  metric text,
  timestamp bigint,
  value double,
  primary key ((metric), timestamp)
);
*/

import (
  "log"
  "fmt"
  "time"

  "github.com/gocql/gocql"
)

type sample struct {
  name string
  timestamp int64
  value float64
}

func generator(metrics int, ch chan *sample) {
  for {
    log.Printf("Starting next %d at %s", metrics, time.Now())
    for i :=0; i<metrics; i++ {
      now := time.Now().UnixNano()
      ch <- &sample{
        name: fmt.Sprintf("metric012345678900123456789001234567890012345678900123456789001234567890012345678900123456789001234567890012345678900123456789001234567890012345678900123456789001234567890012345678900123456789001234567890%d", i),
        timestamp: now / 1e6,
        value: float64(i) * float64(now),
      }
    }
  }
}

func worker(session *gocql.Session, ch chan *sample) {
  for s := range ch {
    if err := session.Query(`INSERT INTO samples (metric, timestamp, value) VALUES (?, ?, ?)`,
      s.name, s.timestamp, s.value).Exec(); err != nil {
        log.Fatalf("Query error: %s", err)
      }
  }
}

func main() {
   cluster := gocql.NewCluster("127.0.0.1")
   cluster.Keyspace = "example"
   cluster.Consistency = gocql.One
   cluster.ProtoVersion = 4
   cluster.NumConns = 10
   session, err := cluster.CreateSession()
   if err != nil {
     log.Fatalf("Error creating session: %s", err)
   }
   defer session.Close()

   ch := make(chan *sample, 10240)

   go generator(100000, ch)

   for i :=0; i < 40; i++ {
     go worker(session, ch)
   }
   select{}
}

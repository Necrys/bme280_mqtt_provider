package main

import (
  "bme280_mqtt_provider/bme280"
  "log"
  "os"
  "time"
  "flag"
  "fmt"
  MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
  log.Println( "Init" )

  topic := flag.String( "topic", "", "The topic name to which to publish" )
  broker := flag.String( "broker", "", "The broker URI. ex: tcp://10.10.1.1:1883" )
  user := flag.String( "user", "", "The user (optional)" )
	password := flag.String( "password", "", "The password (optional)" )
  id := flag.String( "id", "!!!ID NOT SET!!!", "The ClientID" )
  period := flag.Uint64( "period", 60, "Sensor reading period in seconds (default: 60)" )
  flag.Parse()

  if *topic == "" {
    log.Panicln( "-topic must be set and not be empty" )
  }

  if *broker == "" {
    log.Panicln( "-broker must be set and not be empty" )
  }

  log.Printf( "topic: %s\n", *topic );
  log.Printf( "broker: %s\n", *broker );
  log.Printf( "user: %s\n", *user );
  log.Printf( "id: %s\n", *id );
  log.Printf( "period (sec): %d\n", *period );

	opts := MQTT.NewClientOptions()
	opts.AddBroker( *broker )
	opts.SetClientID( *id )
	opts.SetUsername( *user )
	opts.SetPassword( *password )

  isWorking := true
  sigs := make( chan os.Signal, 1 )
  go func() {
    sig := <-sigs
    log.Println( sig )
    isWorking = false;
  } ()

  log.Println( "Start" )

  bme280conn, err := bme280.Connect( 118, 1 )
  if err != nil {
    log.Panicln( err )
  }

  defer bme280conn.Disconnect()

  log.Println( "BME280 Connected" )

  mqttClient := MQTT.NewClient( opts )
  if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
    log.Panicln( token.Error() )
  }

  defer mqttClient.Disconnect( 250 )
  
  log.Println( "MQTT Publisher Started" )

  for isWorking {
    time.Sleep( time.Duration( *period ) * time.Second )

    timestamp := time.Now().Local().Format( "02012006150405MST" )
    temperature, humidity, pressure, err := bme280conn.ReadData()
    if err != nil {
      log.Println( err )
      continue
    }

    payload := fmt.Sprintf( "%s:%s:%f:%f:%f", *id, timestamp, temperature, humidity, pressure )

    token := mqttClient.Publish( *topic, 0, false, payload )
    token.Wait()

    if token.Error() != nil {
      log.Println( token.Error() )
      continue
    }
  }

  log.Println( "Shutdown" )
}

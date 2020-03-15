# bme280_mqtt_provider
MQTT data publisher of bme280 readings

## Usage options:
* -broker <URI> : mandatory, broker address (for ex.: tcp://my.server.com:1883)
* -topic <TOPIC> : mandatory, which topic to publish
* -id <ID> : optional, ID will be added to the payload in the beginning
* -user <USER> : optional, MQTT user
* -password <PASS> : optional, MQTT user password

## init.d service script installation steps:
1. cp ./init.d/bme280mqtt /etc/init.d/
1. chmod +x /etc/init.d/bme280mqtt
1. edit script, set correct values for the script parameters
    1. SCRIPT : path to the service executable
    1. RUNAS : run service from user
    1. ID : service -id parameter
    1. TOPIC : service -topic parameter
    1. BROKER : service -broker parameter
    1. PERIOD : service -period parameter
1. touch /var/log/bme280mqtt.log
1. chown <your_user_name> /var/log/bme280mqtt.log
1. update-rc.d bme280mqtt defaults
1. service bme280mqtt start

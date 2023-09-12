echo "system,40" | kafka-console-producer --broker-list localhost:9092 --topic system_topic
echo "verdict,1" | kafka-console-producer --broker-list localhost:9092 --topic rm_topic
echo "system,30" | kafka-console-producer --broker-list localhost:9092 --topic system_topic
echo "verdict,0" | kafka-console-producer --broker-list localhost:9092 --topic rm_topic
echo "system,35" | kafka-console-producer --broker-list localhost:9092 --topic system_topic
echo "verdict,0" | kafka-console-producer --broker-list localhost:9092 --topic rm_topic
echo "system,30" | kafka-console-producer --broker-list localhost:9092 --topic system_topic
echo "verdict,1" | kafka-console-producer --broker-list localhost:9092 --topic rm_topic
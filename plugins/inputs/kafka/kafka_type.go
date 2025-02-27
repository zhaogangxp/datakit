package kafka

var KafkaTypeMap = map[string]string{
	"ActiveControllerCount.Value":                         "int",
	"AutoLeaderBalanceRateAndTimeMs.Count":                "int",
	"ControlledShutdownRateAndTimeMs.Count":               "int",
	"ControllerChangeRateAndTimeMs.Count":                 "int",
	"ControllerShutdownRateAndTimeMs.Count":               "int",
	"EventQueueSize.Value":                                "int",
	"EventQueueTimeMs.Count":                              "int",
	"GlobalPartitionCount.Value":                          "int",
	"GlobalTopicCount.Value":                              "int",
	"IsrChangeRateAndTimeMs.Count":                        "int",
	"LeaderAndIsrResponseReceivedRateAndTimeMs.Count":     "int",
	"LeaderElectionRateAndTimeMs.Count":                   "int",
	"ListPartitionReassignmentRateAndTimeMs.Count":        "int",
	"LogDirChangeRateAndTimeMs.Count":                     "int",
	"ManualLeaderBalanceRateAndTimeMs.Count":              "int",
	"OfflinePartitionsCount.Value":                        "int",
	"PartitionReassignmentRateAndTimeMs.Count":            "int",
	"PreferredReplicaImbalanceCount.Value":                "int",
	"ReplicasIneligibleToDeleteCount.Value":               "int",
	"ReplicasToDeleteCount.Value":                         "int",
	"TopicChangeRateAndTimeMs.Count":                      "int",
	"TopicDeletionRateAndTimeMs.Count":                    "int",
	"TopicUncleanLeaderElectionEnableRateAndTimeMs.Count": "int",
	"TopicsIneligibleToDeleteCount.Value":                 "int",
	"TopicsToDeleteCount.Value":                           "int",
	"TotalQueueSize.Value":                                "int",
	"UncleanLeaderElectionEnableRateAndTimeMs.Count":      "int",
	"UncleanLeaderElectionsPerSec.Count":                  "int",
	"UpdateFeaturesRateAndTimeMs.Count":                   "int",
	"AtMinIsrPartitionCount.Value":                        "int",
	"FailedIsrUpdatesPerSec.Count":                        "int",
	"IsrExpandsPerSec.Count":                              "int",
	"IsrShrinksPerSec.Count":                              "int",
	"LeaderCount.Value":                                   "int",
	"OfflineReplicaCount.Value":                           "int",
	"PartitionCount.Value":                                "int",
	"UnderMinIsrPartitionCount.Value":                     "int",
	"UnderReplicatedPartitions.Value":                     "int",
	"AlterAcls.NumDelayedOperations":                      "int",
	"AlterAcls.PurgatorySize":                             "int",
	"DeleteRecords.NumDelayedOperations":                  "int",
	"DeleteRecords.PurgatorySize":                         "int",
	"ElectLeader.NumDelayedOperations":                    "int",
	"ElectLeader.PurgatorySize":                           "int",
	"Fetch.NumDelayedOperations":                          "int",
	"Fetch.PurgatorySize":                                 "int",
	"Heartbeat.NumDelayedOperations":                      "int",
	"Heartbeat.PurgatorySize":                             "int",
	"Produce.NumDelayedOperations":                        "int",
	"Produce.PurgatorySize":                               "int",
	"Rebalance.NumDelayedOperations":                      "int",
	"Rebalance.PurgatorySize":                             "int",
	"topic.NumDelayedOperations":                          "int",
	"topic.PurgatorySize":                                 "int",
	"LocalTimeMs.Count":                                   "int",
	"RemoteTimeMs.Count":                                  "int",
	"RequestBytes.Count":                                  "int",
	"RequestQueueTimeMs.Count":                            "int",
	"ResponseQueueTimeMs.Count":                           "int",
	"ResponseSendTimeMs.Count":                            "int",
	"ThrottleTimeMs.Count":                                "int",
	"TotalTimeMs.Count":                                   "int",
	"BytesInPerSec.Count":                                 "int",
	"BytesOutPerSec.Count":                                "int",
	"BytesRejectedPerSec.Count":                           "int",
	"FailedFetchRequestsPerSec.Count":                     "int",
	"FailedProduceRequestsPerSec.Count":                   "int",
	"FetchMessageConversionsPerSec.Count":                 "int",
	"InvalidMagicNumberRecordsPerSec.Count":               "int",
	"InvalidMessageCrcRecordsPerSec.Count":                "int",
	"InvalidOffsetOrSequenceRecordsPerSec.Count":          "int",
	"MessagesInPerSec.Count":                              "int",
	"NoKeyCompactedTopicRecordsPerSec.Count":              "int",
	"ProduceMessageConversionsPerSec.Count":               "int",
	"ReassignmentBytesInPerSec.Count":                     "int",
	"ReassignmentBytesOutPerSec.Count":                    "int",
	"ReplicationBytesInPerSec.Count":                      "int",
	"ReplicationBytesOutPerSec.Count":                     "int",
	"TotalFetchRequestsPerSec.Count":                      "int",
	"TotalProduceRequestsPerSec.Count":                    "int",
	"LogEndOffset":                                        "int",
	"LogStartOffset":                                      "int",
	"NumLogSegments":                                      "int",
	"Size":                                                "int",
	"UnderReplicatedPartitions":                           "int",
}

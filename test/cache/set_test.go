package cache_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"
	"unicode/utf8"
)

func TestHash(t *testing.T) {
	if err := redisClient.HSet(context.Background(), "test", map[string]any{"k1": "v1"}).Err(); err != nil {
		t.Fatal(err)
	}
	if err := redisClient.HSet(context.Background(), "test", map[string]any{"k2": "v2"}).Err(); err != nil {
		t.Fatal(err)
	}

	if result, err := redisClient.HGetAll(context.Background(), "test").Result(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

func TestSet(t *testing.T) {
	if err := redisClient.SAdd(context.Background(), "testSet", "v1").Err(); err != nil {
		t.Fatal(err)
	}
	if err := redisClient.SAdd(context.Background(), "testSet", "v2").Err(); err != nil {
		t.Fatal(err)
	}

	if err := redisClient.Expire(context.Background(), "testSet", 10*time.Second).Err(); err != nil {
		t.Fatal(err)
	}

	if result, err := redisClient.SMembers(context.Background(), "testSet").Result(); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(result)
	}
}

var j string = `{
    "Name": "evalsha",
    "SpanContext": {
        "TraceID": "bda69652621c1ef2db39e705c4a3bed9",
        "SpanID": "230c169f3fd4c1ae",
        "TraceFlags": "01",
        "TraceState": "",
        "Remote": false
    },
    "Parent": {
        "TraceID": "00000000000000000000000000000000",
        "SpanID": "0000000000000000",
        "TraceFlags": "00",
        "TraceState": "",
        "Remote": false
    },
    "SpanKind": 3,
    "StartTime": "2025-08-29T11:35:46.370456286+08:00",
    "EndTime": "2025-08-29T11:35:46.370657758+08:00",
    "Attributes": [
        {
            "Key": "db.system.name",
            "Value": {
                "Type": "STRING",
                "Value": "redis"
            }
        },
        {
            "Key": "db.query.text",
            "Value": {
                "Type": "STRING",
                "Value": "evalsha a8153319360adc71cdc370107f8ab9786f204765 4 asynq:{default}:pending asynq:{default}:paused asynq:{default}:active asynq:{default}:lease 1756438576 asynq:{default}:t:: evalsha a8153319360adc71cdc370107f8ab9786f204765 4 asynq:{default}:pending asynq:{default}:paused asynq:{default}:active asynq:{default}:lease 1756438576 asynq:{default}:t:"
            }
        },
        {
            "Key": "db.operation.name",
            "Value": {
                "Type": "STRING",
                "Value": "evalsha"
            }
        },
        {
            "Key": "server.address",
            "Value": {
                "Type": "STRING",
                "Value": "tc-common-test.redis.tucinfra.com:6379"
            }
        },
        {
            "Key": "db.collection.name",
            "Value": {
                "Type": "STRING",
                "Value": ""
            }
        }
    ],
    "Events": [
        {
            "Name": "exception",
            "Attributes": [
                {
                    "Key": "exception.type",
                    "Value": {
                        "Type": "STRING",
                        "Value": "github.com/redis/go-redis/v9/internal/proto.RedisError"
                    }
                },
                {
                    "Key": "exception.message",
                    "Value": {
                        "Type": "STRING",
                        "Value": "redis: nil"
                    }
                }
            ],
            "DroppedAttributeCount": 0,
            "Time": "2025-08-29T11:35:46.370661566+08:00"
        }
    ],
    "Links": null,
    "Status": {
        "Code": "Error",
        "Description": ""
    },
    "DroppedAttributes": 0,
    "DroppedEvents": 0,
    "DroppedLinks": 0,
    "ChildSpanCount": 0,
    "Resource": [
        {
            "Key": "service.name",
            "Value": {
                "Type": "STRING",
                "Value": "tc-resource-ogc.test"
            }
        },
        {
            "Key": "telemetry.sdk.language",
            "Value": {
                "Type": "STRING",
                "Value": "go"
            }
        },
        {
            "Key": "telemetry.sdk.name",
            "Value": {
                "Type": "STRING",
                "Value": "opentelemetry"
            }
        },
        {
            "Key": "telemetry.sdk.version",
            "Value": {
                "Type": "STRING",
                "Value": "1.35.0"
            }
        }
    ],
    "InstrumentationScope": {
        "Name": "pkg/rules/goredis/setup.go",
        "Version": "v0.7.0",
        "SchemaURL": "",
        "Attributes": null
    },
    "InstrumentationLibrary": {
        "Name": "pkg/rules/goredis/setup.go",
        "Version": "v0.7.0",
        "SchemaURL": "",
        "Attributes": null
    }
}`

func TestXxx(t *testing.T) {
	var m map[string]any
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		t.Fatal(err)
	}

	mtt, ok := m["Attributes"].([]any)
	if !ok {
		t.Fatal("not array")
	}

	for _, v := range mtt {
		vv, ok := v.(map[string]any)
		if !ok {
			t.Fatal("not map")
		}
		vvv, ok := vv["Value"].(map[string]any)
		if !ok {
			t.Fatal("not map")
		}
		if vvv["Type"] == "STRING" {
			if !utf8.ValidString(vvv["Value"].(string)) {
				t.Fatal("not valid string")
			}
		}
	}
}

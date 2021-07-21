# proto-dup
Read and duplicate a `.proto` file and allow to modify the original with special pre/postfixes and removing/adding optional to all fields.

## Installation & Usage
```
$ go get github.com/davidae/proto-dup
$ go install github.com/davidae/proto-dup
$ proto-dup -h
Usage of proto-dup:
  -add-optional
        add optional to all applicable fields
  -go_package string
        set a new go_package name
  -out string
        output of duplicate .proto file
  -package string
        set a new package name
  -postfix string
        add a postfix to all message and enum types
  -prefix string
        add a prefix to all message and enum types
  -remove-optional
        remove optional from all applicable fields

```
## Example command
```
$ proto-dup --package wellhellothere --prefix Hello --postfix World --add-optional $GOPATH/src/github.com/davidae/proto-dup/test/example.proto
syntax = "proto3";

package wellhellothere;
option go_package = "github.com/davidae/proto-dup;proto";

import "google/protobuf/timestamp.proto";

message HelloMessageWorld {
        optional string name = 1;
        optional int64 birthDay = 2;
        optional string phone = 3;
        optional int32 siblings = 4;
        optional bool spouse = 5;
        optional double money = 6;
        optional HelloTypeWorld type = 7;
        optional HelloAddressWorld addresss = 8;
        optional google.protobuf.Timestamp created_at = 9;
        oneof values {
                string value_s = 10;
                int32 value_i = 11;
                double value_d = 12;
        }
}

message HelloAddressWorld {
        optional string street = 1;
        optional int32 number = 2;
        optional int32 post_code = 3;
        optional int32 floor = 4;
        optional bool use_for_billing = 5;
}

enum HelloTypeWorld {
        TYPE_UNSPECIFIED = 0;
        TYPE_R = 1;
        TYPE_S = 2;
}
```
package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlowItem struct {
	/// From: "s" / "c" for server or client
	From string `json:"from" bson:"from"`
	/// Data, in a somewhat readable format
	Data string `json:"data" bson:"data"`
	/// The raw data, base64 encoded.
	B64 string `json:"b64" bson:"b64"`
	/// Timestamp of the first packet in the flow (Epoch / ms)
	Time int `json:"time" bson:"time"`
}

type FlowEntry struct {
	Id           primitive.ObjectID `json:"_id" bson:"_id"`
	Src_port     int                `json:"src_port" bson:"src_port"`
	Dst_port     int                `json:"dst_port" bson:"dst_port"`
	Src_ip       string             `json:"src_ip" bson:"src_ip"`
	Dst_ip       string             `json:"dst_ip" bson:"dst_ip"`
	Time         int                `json:"time" bson:"time"`
	Duration     int                `json:"duration" bson:"duration"`
	Num_packets  int                `json:"num_packets" bson:"num_packets"`
	Blocked      bool               `json:"blocked" bson:"blocked"`
	Filename     string             `json:"filename" bson:"filename"`
	Parent_id    primitive.ObjectID `json:"parent_id" bson:"parent_id"`
	Child_id     primitive.ObjectID `json:"child_id" bson:"child_id"`
	Fingerprints []uint32           `json:"fingerprints" bson:"fingerprints"`
	Suricata     []int              `json:"suricata" bson:"suricata"`
	Flow         []FlowItem         `json:"flow" bson:"flow"`
	Tags         []string           `json:"tags" bson:"tags"`
	Size         int                `json:"size" bson:"size"`
}

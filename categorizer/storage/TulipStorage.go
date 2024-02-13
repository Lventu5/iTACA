package storage

import "go.mongodb.org/mongo-driver/bson/primitive"

type FlowItem struct {
	/// From: "s" / "c" for server or client
	From string
	/// Data, in a somewhat readable format
	Data string
	/// The raw data, base64 encoded.
	B64 string
	/// Timestamp of the first packet in the flow (Epoch / ms)
	Time int
}

/*type FlowEntry struct {
	Src_port     int
	Dst_port     int
	Src_ip       string
	Dst_ip       string
	Time         int
	Duration     int
	Num_packets  int
	Blocked      bool
	Filename     string
	Parent_id    primitive.ObjectID
	Child_id     primitive.ObjectID
	Fingerprints []uint32
	Suricata     []int
	Flow         []FlowItem
	Tags         []string
	Size         int
}*/

type FlowEntry struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	Time         int                `json:"time" bson:"time"`
	Duration     int                `json:"duration" bson:"duration"`
	Src_IP       string             `json:"src_ip" bson:"src_ip"`
	Dst_IP       string             `json:"dst_ip" bson:"ip_dst"`
	Src_Port     uint16             `json:"src_port" bson:"src_port"`
	Dst_Port     uint16             `json:"dst_port" bson:"dst_port"`
	ContainsFlag bool               `json:"contains_flag" bson:"contains_flag"`
	Flow         []FlowItem         `json:"flow" bson:"flow"`
}

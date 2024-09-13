package protocol

// TargetType 定义消息的目标类型
type TargetType string

const (
	// SingleTarget 单发消息
	SingleTarget TargetType = "single"
	// GroupTarget 群发消息
	GroupTarget TargetType = "group"
)

// MessageType 定义消息的类型
type MessageType string

const (
	// TextMessage 文本消息
	TextMessage MessageType = "text"
	// ImageMessage 图片消息
	ImageMessage MessageType = "image"
	// VideoMessage 视频消息
	VideoMessage MessageType = "video"
	// FileMessage 文件消息
	FileMessage MessageType = "file"
	// SystemMessage 系统通知消息
	SystemMessage MessageType = "system"
)

// MessageStatus 定义消息状态
type MessageStatus string

const (
	// MessageSent 消息已发送
	MessageSent MessageStatus = "sent"
	// MessageDelivered 消息已送达
	MessageDelivered MessageStatus = "delivered"
	// MessageRead 消息已读
	MessageRead MessageStatus = "read"
)

// MessageProtocol 定义消息协议的结构体
type MessageProtocol struct {
	ID         string                 `json:"id"`                    // 消息唯一ID
	Type       MessageType            `json:"type"`                  // 消息类型 (如: "text", "image")
	Status     MessageStatus          `json:"status"`                // 消息状态
	SenderID   string                 `json:"sender_id"`             // 发送方用户ID
	ReceiverID string                 `json:"receiver_id,omitempty"` // 接收方用户ID，单发消息时使用
	GroupID    string                 `json:"group_id,omitempty"`    // 群组ID，群发消息时使用
	DeviceID   string                 `json:"device_id"`             // 发送方设备ID
	Content    string                 `json:"content"`               // 消息内容
	TargetType TargetType             `json:"target_type"`           // 消息目标类型 (如: "single", "group")
	Timestamp  string                 `json:"timestamp"`             // 时间戳，ISO8601 格式
	ExtraData  map[string]interface{} `json:"extra_data,omitempty"`  // 额外数据，用于扩展字段
}

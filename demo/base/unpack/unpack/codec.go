package unpack

import (
	"encoding/binary"
	"errors"
	"io"
)

// MMsg_Header 一共８个字节，所以这里指定写入８个字节的数据
const Msg_Header = "12345678"

// Encode 数据进行编码
func Encode(bytesBuffer io.Writer, content string) error {
	var err error
	// 定义消息的格式：msg_header+content_len+content
	// 对应字节长度分别是：8(header)+4(消息长度)+content_len

	// 写入数据到buffer里
	if err = binary.Write(bytesBuffer, binary.BigEndian, []byte(Msg_Header)); err != nil {
		return err
	}

	// 写入contentlen，是根据content的实际大小去写入的
	// 这里控制最长不超过int32的长度
	clen := int32(len([]byte(content)))
	if err = binary.Write(bytesBuffer, binary.BigEndian, clen); err != nil {
		return err
	}

	// 最后写入content大小
	if err = binary.Write(bytesBuffer, binary.BigEndian, []byte(content)); err != nil {
		return err
	}

	return nil

}

// Decode 对数据进行解码
func Decode(bytesBuffer io.Reader) (bodyBuf []byte, err error) {

	// 首先读取header大小的长度
	MagicBuf := make([]byte, len(Msg_Header))
	if _, err = io.ReadFull(bytesBuffer, MagicBuf); err != nil {
		return nil, err
	}

	// 比较读取的header和实际的header是否相同
	if string(MagicBuf) != Msg_Header {
		return nil, errors.New("msg_header error")
	}

	// 读取数据的长度，4个字节
	lengthBuf := make([]byte, 4)
	if _, err := io.ReadFull(bytesBuffer, lengthBuf); err != nil {
		return nil, err
	}

	// 把读取到的数据转换成实际的length
	// 首先拿到的是二进制数据，首先进行大段字节序解码一下拿到实际length，因为TCP实际传输的时候会用大段字节序加密
	length := binary.BigEndian.Uint32(lengthBuf)

	// 读取实际数据的长度
	bodyBuf = make([]byte, length)
	if _, err := io.ReadFull(bytesBuffer, bodyBuf); err != nil {
		return nil, err
	}
	return bodyBuf, err
}

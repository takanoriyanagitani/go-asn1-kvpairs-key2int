import asn1tools
import sys

k2int = asn1tools.compile_files("./key2int.asn")
encoded = k2int.encode(
	"KeyToIntMapItems",
	[
		dict(originalKey = "unspecified", mapdSerial = 0),
		dict(originalKey = "timestamp", mapdSerial = 1),
		dict(originalKey = "severity", mapdSerial = 2),
		dict(originalKey = "message_id", mapdSerial = 3),
		dict(originalKey = "message", mapdSerial = 4),
		dict(originalKey = "http_method", mapdSerial = 5),
		dict(originalKey = "url", mapdSerial = 6),
		dict(originalKey = "http_status", mapdSerial = 7),
		dict(originalKey = "error_code", mapdSerial = 8),
		dict(originalKey = "tag", mapdSerial = 9),
		dict(originalKey = "tags", mapdSerial = 10),
		dict(originalKey = "version", mapdSerial = 11),
		dict(originalKey = "instance", mapdSerial = 12),
		dict(originalKey = "machine", mapdSerial = 13),
		dict(originalKey = "image", mapdSerial = 14),
		dict(originalKey = "architecture", mapdSerial = 15),
	]
)
sys.stdout.buffer.write(encoded)

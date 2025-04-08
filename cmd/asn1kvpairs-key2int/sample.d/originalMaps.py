import asn1tools
import sys

k2int = asn1tools.compile_files("./key2int.asn")
encoded = k2int.encode(
	"OriginalPairs",
	[
		dict(key="timestamp", val="2025-04-07T07:09:01.0Z"),
		dict(key="severity", val="info"),
		dict(key="message_id", val="cafef00d-dead-beaf-face-864299792458"),
		dict(key="message", val="request processed."),
		dict(key="http_method", val="GET"),
		dict(key="url", val="http://example.com/"),
		dict(key="http_status", val="200"),
		dict(key="tags", val="http,read_only"),
		dict(key="version", val="1.0.0"),
		dict(key="instance", val="test-container"),
		dict(key="machine", val="bare-metal"),
		dict(key="image", val="test-image:1.0.0"),
		dict(key="architecture", val="arm64"),
	]
)
sys.stdout.buffer.write(encoded)

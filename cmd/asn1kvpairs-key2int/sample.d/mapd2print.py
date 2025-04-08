import asn1tools
import sys
import json

k2int = asn1tools.compile_files("./key2int.asn")
encoded = sys.stdin.buffer.read()
decoded = k2int.decode(
	"NormalizedPairs",
	encoded,
)
mapd = map(json.dumps, decoded)
prints = map(print, mapd)
sum(1 for _ in prints)

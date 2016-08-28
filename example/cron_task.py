import json

a = {
    "ToServices": ["server"],
    "MessageType": "Memory",
    "message": "1,2,3,4,5"
}
print(json.dumps(a))

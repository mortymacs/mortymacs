# Create a new job.
POST {{BASE_URL}}/jobs
{
    "command": "ls -lha",
    "callback_url": ""
}

HTTP 202
[Asserts]
jsonpath "$.status" == "ok"
jsonpath "$.attributes.process_id" exists
[Captures]
process_id: jsonpath "$.attributes.process_id"

# Get the job info.
GET {{BASE_URL}}/jobs/{{process_id}}

HTTP 200
[Asserts]
jsonpath "$.status" == "ok"
jsonpath "$.attributes.id" == {{process_id}}
jsonpath "$.attributes.status" == "pending"

import json

#json string data
employee_string = '{"password": "correcthorsebatterystaple123(!__+@**(A\'\"`;DROP TABLE packages;""last_name": "Rodgers", "department": "Marketing"}'

#check data type with type() method
json_object = json.loads(employee_string)
print(json_object)
import requests
import os

# ---- constants ----

global exit_value
exit_value = 0
test_tasks = ["AvayaOneXAgentStart", "AvayaOneXAgentClose"]
valid_auth_token = "7abcddbb2c74e4c0789c2c0aa6abcf5172e82e9f4916bc6409fc3989ed673e08"
another_valid_auth_token = "7cd477192d54ceb8673be093f443b8622c612896880f6879c7f8ec16fa7ba114"
invalid_auth_token = "helloworld1234"

source="gitlab/ERPK-T-IRRP/test-proxy"
url = "http://10.137.54.217:8080"

# ---- functions ----

def manageLock(auth_token):

    endpoint = "/api/lock"

    headers = {
        "auth": auth_token
    }

    r = requests.get(url + endpoint, headers=headers)
    print("Lock Manager: ", r.status_code)
    return r.status_code

test_tasks = ["AvayaOneXAgentStart", "AvayaOneXAgentClose"]
valid_auth_token = "7abcddbb2c74e4c0789c2c0aa6abcf5172e82e9f4916bc6409fc3989ed673e08"
another_valid_auth_token = "7cd477192d54ceb8673be093f443b8622c612896880f6879c7f8ec16fa7ba114"
invalid_auth_token = "helloworld1234"
source="gitlab/testingTheTest"
url = "http://10.137.54.217:8080"

# ---- functions ----

def manageLock(auth_token):

    endpoint = "/api/lock"

    headers = {
        "auth": auth_token
    }

    r = requests.get(url + endpoint, headers=headers)
    print("Lock Manager: ", r.status_code)
    return r.status_code

    
def test(Tasktype, task, source, auth_token):
    
    endpoint = "/api/post"

    payload = {
        "type": Tasktype,
        "title": "My Test Project",
        "task": task
    }

    headers = {
        "source": source,
        "auth": auth_token
    }
    
    r = requests.post(url + endpoint, json=payload, headers=headers)

    print("Task: ", r.status_code)
    print(r.json())
    
    try:
        if r.json()['results'][0]['status'] != 'OK':
            print('failed')
            global exit_value
            exit_value = 1
    except:
        print("excepted response: ", r.status_code)
        pass


        
# ---- tests ----

# set lock
if manageLock(valid_auth_token) != 200:
    print("error getting lock with valid token")


# soft regression (list of tasks)
try:
    for task in test_tasks:
        test("test", task, source, valid_auth_token)
except:
    print("an error occured in soft regression")
    pass


# hard regression (preconfigured regression file)
try:
    test("regression", "testRegression.conf", source, valid_auth_token)
except:
    print("an error occured in hard regression")
    pass


# release lock
if manageLock(valid_auth_token) != 200:
    print("error releasing lock with valid token")


# ---- invalid auth test ----

if manageLock(invalid_auth_token) == 200:
    print("error getting lock with invalid token")
    
if manageLock(invalid_auth_token) == 200:
    print("error releasing lock with invalid token")

# ---- lock stability test ----

# set test
if manageLock(valid_auth_token) != 200:
    print("error getting lock with valid token")

# set another lock => expected: 401
if manageLock(another_valid_auth_token) == 200:
    print("error lock didnt work correctly")

# checkout lock
if manageLock(valid_auth_token) != 200:
    print("error getting lock with valid token")

# ---- the end ---- 
print("Exit Value is ", exit_value)
os._exit(exit_value)

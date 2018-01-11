import requests
import random
import time

print "start at "+str(time.time())
for i in range(0, 10000):
    requests.get("http://127.0.0.1:8080/?pg1="+str(random.randint(0,1000))+"&protocol=telegram&subscriber=123345&wnumber="+str(random.randint(100000,1000000)))
print "stop at "+str(time.time())
import sys
import json
import time
import random
import urllib.request
from threading import Thread, Lock

'''
可调整的参数
BEGIN
'''

domain = 'http://127.0.0.1:18800'

urls = [
        domain+'/hkg/metatable/Healthy/Echo',
        ]

params = {
        'msg':'hello',
        }

options = {
        'content-type':'application/json;charset=UTF-8'
        }

thread_count = 100
invoke_count = 10
interval = 0


'''
可调整的参数
END
'''



failed = 0
successed = 0
thread_finished = 0
errors = {}
def onThreadFinish(_id):
    global thread_finished
    lock.acquire()
    thread_finished = thread_finished + 1
    lock.release()
    if thread_finished == thread_count:
        print('访问成功总次数：%d'%(successed))
        print('访问失败总次数：%d'%(failed))
        print('----------------------------------------------------------')
        print('错误列表：')
        for err in errors:
            print(err)
        print('##########################################################')

def call(_id, _n):
    global failed
    global successed
    global interval
    global errors
    for i in range(0, _n):
        time.sleep(interval)
        data = json.dumps(params)
        data = bytes(data, 'utf8')
        url = urls[random.randint(0, len(urls)-1)]
        request = urllib.request.Request(url, data, options)
        try:
            reply = urllib.request.urlopen(request).read().decode('utf-8')
            #print(reply)
            lock.acquire()
            successed = successed + 1
            lock.release()
        except Exception as e:
            err = str.format('%s -> %s'%(url,e))
            lock.acquire()
            failed = failed + 1
            errors[err] = 0
            lock.release()
    onThreadFinish(_id)

print('##########################################################')
print('接口列表：')
for url in urls:
    print(url)
print('----------------------------------------------------------')
print('线程数：%d '%(thread_count))
print('每线程访问次数：%d '%(invoke_count))
print('访问间隔时间：%f '%(interval))
print('**********************************************************')
# 创建100个线程
for i in range(0,thread_count):
    lock = Lock()
    # 每个线程10次访问
    t = Thread(target=call, args=(i, invoke_count,))
    t.start()

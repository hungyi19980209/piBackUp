#! /usr/bin/env python
#-*- coding:utf-8 -*-
# version : Python 2.7.13
 
import os,sys,time
import socket
from evdev import InputDevice, categorize, ecodes
import evdev 
def doConnect(host,port):
    sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    try :         
        sock.connect((host,port))
    except :
        pass
    return sock
         
def main():   
    device = InputDevice("/dev/input/event0") # my keyboard1
    number = ""
    host,port = "127.0.0.1",8069
    print host,port    
     
      
    while True :
        try :
            for event in device.read_loop():
                if event.type == ecodes.EV_KEY:
                    if event.value==1:
                        a = evdev.ecodes.KEY[event.code]
                        if a[4:] == "ENTER":
                            msg=number
                            sockLocal = doConnect(host,port)
                            sockLocal.send(msg)  
                            number=""
                            sockLocal.close()
                            sockLocal = doConnect(host,port) 
                        else:
                            print("hehe")
                            number = number + a[4:]
            
     
        except socket.error :
            print "\r\nsocket error,do reconnect "
            time.sleep(3)
            sockLocal = doConnect(host,port)   
        except :
            print '\r\nother error occur '           
            time.sleep(3) 
        time.sleep(1)
     
if __name__ == "__main__" :
    main()

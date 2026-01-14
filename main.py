import socket
import time
import random
import selectors
import struct
from client import make_packet

sel = selectors.DefaultSelector()
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind(("localhost", 9090))
sock.setblocking(False)

sel.register(sock, selectors.EVENT_READ)

TICK_RATE = 120
TICK_TIME = 1 / TICK_RATE


print("UDP port running on on 9090")
while True:
    start_time = time.time()
    events = sel.select(timeout=0)
    for key, _ in events:
        data, adr = key.fileobj.recvfrom(2048)
        time.sleep(0.05)
        key.fileobj.sendto(data, adr)

    elapsed = time.time() - start_time
    sleep_time = max(elapsed, TICK_TIME)
    time.sleep(sleep_time)




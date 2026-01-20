import socket
import time
import random
import selectors
import struct
from client import make_packet, parse_packet

sel = selectors.DefaultSelector()
sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.bind(("localhost", 9090))
sock.setblocking(False)

sel.register(sock, selectors.EVENT_READ)

TICK_RATE = 120
TICK_TIME = 1 / TICK_RATE


print("UDP port running on on 9090")
received = set()
expected_seq = 0
while True:
    start_time = time.time()
    events = sel.select(timeout=0)
    for key, _ in events:
        data, adr = key.fileobj.recvfrom(2048)
        seq, ack , ts, payload = parse_packet(data)
        received.add(seq)
        while expected_seq in received:
            expected_seq +=1
        
        ack = expected_seq - 1
        data = make_packet(seq, ack, payload)
        key.fileobj.sendto(data, adr)

    elapsed = time.time() - start_time
    sleep_time = max(0, TICK_TIME - elapsed)
    time.sleep(sleep_time)




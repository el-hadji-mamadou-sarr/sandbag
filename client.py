import socket
import time
import struct

HEADER_FMT = "!IId" #uint #double
HEADER_SIZE = struct.calcsize(HEADER_FMT)
print(HEADER_SIZE)

def make_packet(seq,ack, payload) -> bytes:
    return struct.pack(HEADER_FMT, seq, ack, time.time()) + payload

def parse_packet(data):
    seq, ack, ts = struct.unpack(HEADER_FMT, data[:HEADER_SIZE])
    payload = data[HEADER_SIZE:]
    return seq, ack , ts, payload

if __name__ == "__main__":
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    start = time.time()
    pack_lost = 0
    sock.settimeout(0.2)
    W = 5
    seq = 1
    last_acked = 0
    received_seq = []

    expected = 0
    while W > seq - last_acked:
        msg = b"x"*1024
        packet = make_packet(seq, last_acked, msg)
        sock.sendto(packet, ("127.0.0.1", 9090))
        try:
            data, _ = sock.recvfrom(2048)
            _, ack, _, _ = parse_packet(data)
            if ack > last_acked:
                last_acked = ack
            seq+=1
        except socket.timeout:
            print(f"lost seq: {seq} --- ")
            pack_lost += 1

    end= time.time()

    print(f"elepsed {end-start}, lost: {(pack_lost/10)*100}")
    print(f"received seq: {received_seq}")


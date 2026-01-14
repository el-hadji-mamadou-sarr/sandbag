import socket
import time
import struct

HEADER_FMT = "!Id" #uint #double
HEADER_SIZE = struct.calcsize(HEADER_FMT)
print(HEADER_SIZE)

def make_packet(seq, payload) -> bytes:
    return struct.pack(HEADER_FMT, seq, time.time()) + payload

def parse_packet(data):
    seq, ts = struct.unpack(HEADER_FMT, data[:HEADER_SIZE])
    payload = data[HEADER_SIZE:]
    return seq, ts, payload

if __name__ == "__main__":
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

    start = time.time()
    pack_lost = 0
    sock.settimeout(0.2)


    last_seq = -1
    for i in range(10):
        msg = b"x"*1024
        packet = make_packet(i, msg)
        sock.sendto(packet, ("127.0.0.1", 9090))
        try:
            data, _ = sock.recvfrom(2048)
            seq, ts, received_packet = parse_packet(data)
            if seq != last_seq + 1:
                print("out of order")
                last_seq = seq
        except socket.timeout:
            print(f"lost ---")
            pack_lost += 1
    end= time.time()

    print(f"elepsed {end-start}, lost: {(pack_lost/10)*100}")



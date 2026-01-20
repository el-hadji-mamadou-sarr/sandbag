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
    seq = 0
    last_acked = 0
    received_seq = []

    expected = 0
    running = True
    buffer = {}
    for i in range(5):
        # envoyer tant que la fenêtre le permet
        while seq < last_acked + W:
            packet = make_packet(seq, last_acked, b"")
            sock.sendto(packet, ("127.0.0.1", 9090))
            buffer[seq] = {"packet":packet, "send_ts":time.time()}
            seq+=1
        
        # recevoir les ack
        try:
            data, _ = sock.recvfrom(2048)
            _, ack, ts, _ = parse_packet(data)
            last_acked = max(last_acked, ack)
        except socket.timeout:
            pass
        
        # retransmettre si timeout
        for seq, info in buffer.items():
            if time.time() - info["send_ts"] > 1:
                sock.sendto(info["packet"], ("127.0.0.1", 9090))



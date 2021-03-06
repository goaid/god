from __future__ import print_function
import logging
import time
import uuid

import grpc

import session_pb2
import session_pb2_grpc

serial = 0


def generate_messages():
    messages = [
        session_pb2.Header(serial=1, messageType="pb.AuthReq"),
        session_pb2.AuthReq(cookie=uuid.uuid4().hex),
        session_pb2.Header(serial=2, messageType="pb.Heartbeat"),
        session_pb2.Heartbeat(timestamp=time.monotonic_ns()),
        session_pb2.Header(serial=3, messageType="pb.Echo"),
        session_pb2.Echo(content=uuid.uuid4().hex),
    ]
    for msg in messages:
        yield msg


def normalize(req):
    pass


def run():
    # NOTE(gRPC Python Team): .close() is possible on a channel and should be
    # used in circumstances in which the with statement does not fit the needs
    # of the code.
    with grpc.insecure_channel('localhost:50051',
                               options=[('grpc.keepalive_time_ms', 1000)]) as channel:
        stub = session_pb2_grpc.SessionStub(channel)
        for response in stub.Flow(generate_messages()):
            print(response)

if __name__ == '__main__':
    logging.basicConfig()
    run()

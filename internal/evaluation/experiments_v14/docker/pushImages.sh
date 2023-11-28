#!/bin/bash

# echo "Pushing UDP images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-udp
#docker push midarch/fibonaccidistributed:1.14.0-server-udp
#docker push midarch/fibonaccidistributed:1.14.0-client-udp
# echo
# echo

echo "Pushing TCP images"
docker push midarch/fibonaccidistributed:1.14.0-namingserver-tcp
docker push midarch/fibonaccidistributed:1.14.0-server-tcp
docker push midarch/fibonaccidistributed:1.14.0-client-tcp
echo
echo

# echo "Pushing SSL images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-ssl
#docker push midarch/fibonaccidistributed:1.14.0-server-ssl
#docker push midarch/fibonaccidistributed:1.14.0-client-ssl
# echo
# echo

# echo "Pushing QUIC images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-quic
#docker push midarch/fibonaccidistributed:1.14.0-server-quic
#docker push midarch/fibonaccidistributed:1.14.0-client-quic
# echo
# echo

#echo "Pushing HTTP images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-http
#docker push midarch/fibonaccidistributed:1.14.0-server-http
#docker push midarch/fibonaccidistributed:1.14.0-client-http
#echo
#echo

#echo "Pushing HTTPS images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-https
#docker push midarch/fibonaccidistributed:1.14.0-server-https
#docker push midarch/fibonaccidistributed:1.14.0-client-https
#echo
#echo

#echo "Pushing HTTP2 images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-http2
#docker push midarch/fibonaccidistributed:1.14.0-server-http2
#docker push midarch/fibonaccidistributed:1.14.0-client-http2
#echo
#echo

#echo "Pushing RPC images"
#docker push midarch/fibonaccidistributed:1.14.0-namingserver-rpc
#docker push midarch/fibonaccidistributed:1.14.0-server-rpc
#docker push midarch/fibonaccidistributed:1.14.0-client-rpc
#echo
#echo

#echo "Pushing non gMidArch RPC images"
#docker push midarch/fiborpc:1.14.0-server-rpc
#docker push midarch/fiborpc:1.14.0-client-rpc
#echo
#echo
#
#echo "Pushing non gMidArch gRpc images"
#docker push midarch/fibogrpc:1.14.0-server-grpc
#docker push midarch/fibogrpc:1.14.0-client-grpc
#echo
#echo
#
#echo "Pushing non gMidArch RMQ images"
#docker push midarch/fibormq:1.14.0-server-rmq
#docker push midarch/fibormq:1.14.0-client-rmq
#echo
#echo

# echo "Pushing UDPTCP images"
# docker push midarch/fibonaccidistributed:1.14.0-namingserver-udptcp
# docker push midarch/fibonaccidistributed:1.14.0-server-udptcp
# docker push midarch/fibonaccidistributed:1.14.0-client-udptcp
# echo
# echo
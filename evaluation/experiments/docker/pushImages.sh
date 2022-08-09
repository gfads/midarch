#!/bin/bash

echo "Pushing UDP images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-udp
#docker push midarch/fibomiddleware:1.0.2-server-udp
#docker push midarch/fibomiddleware:1.0.3-client-udp
echo
echo

echo "Pushing TCP images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-tcp
#docker push midarch/fibomiddleware:1.0.2-server-tcp
#docker push midarch/fibomiddleware:1.0.3-client-tcp
echo
echo

echo "Pushing SSL images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-ssl
#docker push midarch/fibomiddleware:1.0.2-server-ssl
#docker push midarch/fibomiddleware:1.0.3-client-ssl
echo
echo

echo "Pushing QUIC images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-quic
#docker push midarch/fibomiddleware:1.0.2-server-quic
#docker push midarch/fibomiddleware:1.0.3-client-quic
echo
echo

#echo "Pushing HTTP images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-http
#docker push midarch/fibomiddleware:1.0.2-server-http
#docker push midarch/fibomiddleware:1.0.2-client-http
#echo
#echo

#echo "Pushing HTTPS images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-https
#docker push midarch/fibomiddleware:1.0.2-server-https
#docker push midarch/fibomiddleware:1.0.2-client-https
#echo
#echo

#echo "Pushing HTTP2 images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-http2
#docker push midarch/fibomiddleware:1.0.2-server-http2
#docker push midarch/fibomiddleware:1.0.2-client-http2
#echo
#echo

#echo "Pushing RPC images"
#docker push midarch/fibomiddleware:1.0.2-namingserver-rpc
#docker push midarch/fibomiddleware:1.0.2-server-rpc
#docker push midarch/fibomiddleware:1.0.3-client-rpc
#echo
#echo

#echo "Pushing non gMidArch RPC images"
#docker push midarch/fiborpc:1.0.3-server-rpc
#docker push midarch/fiborpc:1.0.3-client-rpc
#echo
#echo
#
#echo "Pushing non gMidArch gRpc images"
#docker push midarch/fibogrpc:1.0.3-server-grpc
#docker push midarch/fibogrpc:1.0.3-client-grpc
#echo
#echo
#
#echo "Pushing non gMidArch RMQ images"
#docker push midarch/fibormq:1.0.3-server-rmq
#docker push midarch/fibormq:1.0.3-client-rmq
#echo
#echo

echo "Pushing new UDPTCP images"
docker push midarch/newfibomiddleware:1.0.0-namingserver-udptcp
docker push midarch/newfibomiddleware:1.0.0-server-udptcp
docker push midarch/newfibomiddleware:1.0.0-client-udptcp
echo
echo
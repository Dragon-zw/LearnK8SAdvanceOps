FROM ubuntu

RUN apt update && apt install vim curl wget -y

COPY manager a.txt /

# CMD ["/manager","-s","docker"]
CMD ["-s","docker"]

ENTRYPOINT ["/manager"]
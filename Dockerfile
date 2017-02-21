FROM scratch
ADD httpProxy /
EXPOSE 9080
EXPOSE 9090
CMD ["/httpProxy"]


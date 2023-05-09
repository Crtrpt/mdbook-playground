FROM  ubuntu
RUN ["apt","update"]
RUN ["apt","upgrade"]
RUN ["apt","install","git","-y"]
RUN ["apt","install","golang","-y"]
RUN ["apt","install","make","-y"]
WORKDIR "/var/data"
CMD ["go","run","main.go"]

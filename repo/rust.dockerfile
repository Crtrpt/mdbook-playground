FROM  ubuntu
RUN ["apt","update"]
RUN ["apt","upgrade"]
RUN ["apt","install","git","-y"]
RUN ["apt","install","curl","-y"]
RUN ["apt","install","gcc","-y"]
RUN ["curl","--proto", "=https", "--tlsv1.2","-sSf","-o","install.sh","https://sh.rustup.rs"]
RUN ["sh","install.sh","-y"]
RUN ["apt","install","make","-y"]
WORKDIR "/var/data"
CMD ["make"]

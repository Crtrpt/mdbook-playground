FROM  ubuntu
RUN ["apt","update"]
RUN ["apt","upgrade"]
RUN ["apt","install","golang","-y"]
RUN ["apt","install","make","-y"]
RUN ["apt","install","git","-y"]
RUN ["git","clone","https://github.com/Crtrpt/mdbook-playground.git"]
CMD ["make"]

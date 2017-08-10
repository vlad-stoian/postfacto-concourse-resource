FROM concourse/buildroot:git

RUN ln -s /usr/bin/gpg2 /usr/bin/gpg

ADD assets/ /opt/resource/
RUN chmod +x /opt/resource/*

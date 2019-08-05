FROM debian

RUN apt-get update && \
    apt-get install -y \
      curl \
      libasound2-dev \
      bc && \
    apt-get clean && \
    rm -r /var/lib/apt/lists/*

RUN curl https://raw.githubusercontent.com/sqshq/sampler/master/example.yml > /etc/sampler.yml
RUN curl https://github.com/sqshq/sampler/releases/download/v1.0.1/sampler-1.0.1-linux-amd64 > /usr/local/bin/sampler && \
  chmod 755 /usr/local/bin/sampler


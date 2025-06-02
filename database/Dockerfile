FROM postgres:17

# Dependencies for pg_ivm
RUN apt-get update && apt-get install -y \
    build-essential \
    postgresql-server-dev-17 \
    git \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Install pg_ivm
RUN cd /tmp && \
    git clone https://github.com/sraoss/pg_ivm.git && \
    cd pg_ivm && \
    make && \
    make install


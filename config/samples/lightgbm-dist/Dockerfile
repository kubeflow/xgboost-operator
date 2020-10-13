FROM python:3.7.9-slim
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    ca-certificates \
    cmake \
    build-essential \
    gcc \
    g++ \
    git && \
    rm -rf /var/lib/apt/lists/*

RUN git clone --recursive --branch stable --depth 1 https://github.com/Microsoft/LightGBM && \
    mkdir LightGBM/build && \
    cd LightGBM/build && \
    cmake .. && \
    make -j4 && \
    make install && \
    cd ../python-package && \
    python setup.py install_lib && \
    cd ../.. && \
    rm -rf LightGBM

WORKDIR /app
ADD requirements.txt .
RUN pip install -r /app/requirements.txt

COPY data ./data/.
COPY *.py ./

ENTRYPOINT [ "python", "/app/main.py" ]
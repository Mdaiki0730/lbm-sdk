FROM golang:1.19-alpine as build

RUN apk add --no-cache tar

# prepare node data
WORKDIR /node
COPY ./contrib/rosetta/node/data.tar.gz data.tar.gz
RUN tar -zxvf data.tar.gz -C .

# build rollupd
WORKDIR /rollupd
COPY . ./
RUN go build -o rollupd ./l2app/rollupd/

FROM alpine
RUN apk add gcc libc-dev python3 --no-cache

ENV PATH=$PATH:/bin

COPY --from=build /rollupd/rollupd /bin/rollupd

WORKDIR /rosetta
COPY ./contrib/rosetta/configuration ./
RUN chmod +x run_tests.sh
RUN chmod +x send_funds.sh
RUN chmod +x faucet.py

COPY --from=build /node/root /root/
WORKDIR /root/.l2app

RUN chmod -R 0777 ./

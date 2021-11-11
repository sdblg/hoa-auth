# build stage
FROM golang:latest as builder
WORKDIR /app/
ENV USER=sod
ENV UID=1000
ENV USER_GROUP=sod
ENV GID=2000

RUN addgroup --gid $GID $USER_GROUP
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home /home/$USER \
    --ingroup $USER_GROUP \
    --uid $UID \
    $USER
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main src/main.go
RUN chgrp -R +$GID ./
RUN chown -R +$UID ./
RUN chmod -R 0770 ./

# final stage
FROM alpine:latest as final

ENV USER=sod
ENV UID=1000
ENV USER_GROUP=sod
ENV GID=2000
RUN addgroup --gid $GID $USER_GROUP
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home /home/$USER \
    --ingroup $USER_GROUP \
    --uid $UID \
    $USER

WORKDIR /home/$USER

ARG env
ENV CURRENT_ENV=$env

RUN echo "Env: " $env
RUN chown -R $USER:$USER_GROUP .

USER $USER
COPY --from=builder /app/main .
COPY --from=builder /app/scripts/entrypoint.sh ./entrypoint.sh
COPY --from=builder /app/config/ ./config/

ENTRYPOINT ["sh", "./entrypoint.sh"]
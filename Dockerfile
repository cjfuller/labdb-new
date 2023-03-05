FROM rust:1.67-slim-bullseye

ADD . /app
WORKDIR /app
RUN cargo build --release && mv target/release/labdb /app/labdb && rm -fr target
CMD ["/app/labdb"]

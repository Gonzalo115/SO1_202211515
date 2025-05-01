# --- ETAPA DE CONSTRUCCIÓN ---
    FROM rust:latest as builder

    WORKDIR /usr/src/app
    COPY . .
    
    RUN cargo build --release
    
    # --- ETAPA DE EJECUCIÓN ---
    FROM debian:bookworm-slim


    RUN apt-get update && \
    apt-get install -y --no-install-recommends \
    libssl3 \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

    
    WORKDIR /app
    COPY --from=builder /usr/src/app/target/release/deploymentRust .
    
    EXPOSE 8080
    CMD ["./deploymentRust"]

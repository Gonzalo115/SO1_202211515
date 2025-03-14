#!/bin/bash
IMAGE="containerstack/alpine-stress"


# Tipos de contenedores
TYPES=("cpu" "ram" "io" "disk")

# Valores posibles para cpus
CPUS_VALUES=(0.2 0.3 0.4 0.5)

# Crear 10 contenedores aleatorios
for i in {1..10}; do
    TYPE=${TYPES[$RANDOM % ${#TYPES[@]}]}

    case $TYPE in
        "cpu")
            name_cpu="cpu_$(date +%s%N)_$((RANDOM % 1000))"
            cpus=${CPUS_VALUES[$RANDOM % ${#CPUS_VALUES[@]}]}
            docker run -d -ti \
                --cpus="$cpus" \
                --memory="50m" \
                --name="$name_cpu" "$IMAGE" \
                stress \
                --cpu 1 
            ;;
        "ram")
            name_ram="ram_$(date +%s%N)_$((RANDOM % 1000))"
            docker run -d -ti \
                --memory="256m" \
                --name="$name_ram" "$IMAGE" \
                stress \
                --vm 1 --vm-bytes 256M
            ;;
        "io")
            name_io="io_$(date +%s%N)_$((RANDOM % 1000))"
            docker run -d -ti \
                --name="$name_io" "$IMAGE" \
                stress \
                --io 1
            ;;
        "disk")
            name_disk="disk_$(date +%s%N)_$((RANDOM % 1000))"
            docker run -d -ti \
                --name="$name_disk" "$IMAGE" \
                stress \
                --hdd 1 --hdd-bytes 256M
            ;;
    esac
done

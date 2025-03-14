<div style="text-align: center;">

# Universidad San Carlos de Guatemala

## Facultad de ingenieria

### Proyecto 1: Gestor de contenedores

SISTEMAS OPERATIVOS 1

Gonzalo Fenando Pérez Cazún
202211515

</div>

---

# Arquitectura

Enlace para acceder al recurso: [eraser](https://app.eraser.io/workspace/xoEcZUDSyrCmcsEJidLd?origin=share)

---

# BASH

Aplicando los conocimientos adquiridos sobre la creación de imágenes y contenedores de Docker mediante el uso de la CLI y la creación de scripts de bash. Se tendrá un cronjob el cual se estará ejecutando en lapsos de 30 segundos con la siguiente funcionalidad:

1. De manera aleatoria deberá generar 10 contenedores de las imágenes explicadas en la sección de notas.

**Notas:**
Se deberá generar los contenedores de manera aleatoria de la siguiente imagen de docker (https://github.com/containerstack/alpine-stress)

- Contenedor que consuma Ram
- Contenedor que consuma CPU
- Contenedor que consuma I/O
- Contenedor que consuma Disco

-> create_containers.sh

```bash
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
```

Luego configurar el cronjob Para que se ejecute cada 30 segundos:

Abrir:

```bash
crontab -e
```

Agrega esta línea al final:

```bash
* * * * * cd /ruta ; create_containers.sh
* * * * * sleep 30 && cd /ruta ; create_containers.sh
```

---

![]()

---

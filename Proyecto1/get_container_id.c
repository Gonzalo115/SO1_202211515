#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUFFER_SIZE 256

// Función para obtener el Container ID de un proceso
char* get_container_id(int pid, char *container_id, size_t size) {
    char path[64];
    FILE *file;
    char buffer[BUFFER_SIZE];

    // Construir la ruta al archivo cgroup del proceso
    snprintf(path, sizeof(path), "/proc/%d/cgroup", pid);

    // Abrir el archivo cgroup
    file = fopen(path, "r");
    if (!file) {
        perror("Error al abrir el archivo");
        strncpy(container_id, "N/A", size);
        return container_id;
    }

    // Leer el contenido del archivo
    while (fgets(buffer, BUFFER_SIZE, file) != NULL) {
        // Split para encontrar el ID del contenedor
        const char *prefix = "docker-";
        char *start = strstr(buffer, prefix);

        if (start) {
            // Mover el puntero después de "docker-"
            start += strlen(prefix);
            
            // Encontrar el siguiente punto (".")
            char *end = strchr(start, '.');
            
            if (end) {
                // Reemplazar el punto con un carácter nulo para terminar la cadena
                *end = '\0';
                strncpy(container_id, start, size);
                container_id[size - 1] = '\0'; // Asegurarse de que la cadena esté terminada
                fclose(file);
                return container_id;
            }
        }
    }

    fclose(file);
    strncpy(container_id, "N/A", size);
    return container_id;
}

int main(int argc, char *argv[]) {
    if (argc != 2) {
        printf("Uso: %s <PID>\n", argv[0]);
        return 1;
    }

    int pid = atoi(argv[1]);
    char container_id[16]; // Suficiente para 12 caracteres + '\0'

    // Obtener el ID del contenedor
    get_container_id(pid, container_id, sizeof(container_id));

    // Mostrar el resultado
    printf("Container ID para el proceso %d: %s\n", pid, container_id);

    return 0;
}

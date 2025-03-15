#include <linux/module.h>
#include <linux/kernel.h>
#include <linux/proc_fs.h>
#include <linux/seq_file.h>
#include <linux/sched.h>
#include <linux/sched/signal.h>
#include <linux/uaccess.h>
#include <linux/mm.h>
#include <linux/kernel_stat.h>
#include <linux/fs.h>
#include <linux/slab.h>
#include <linux/string.h>

#define PROC_FILENAME "sysinfo_202211515"
#define BUFFER_SIZE 256

// Prototipos de funciones
unsigned int get_cpu_usage(void);
char* get_container_id(int pid, char *container_id, size_t size);

char* get_container_id(int pid, char *container_id, size_t size) {
    char path[128];
    struct file *file;
    ssize_t bytes_read;
    loff_t pos = 0;

    // Construir la ruta al archivo cgroup del proceso
    snprintf(path, sizeof(path), "/home/fernando/Escritorio/a.txt");

    // Abrir el archivo cgroup
    file = filp_open(path, O_RDONLY, 0);
    if (IS_ERR(file)) {
        printk(KERN_INFO "Error al abrir el archivo cgroup para el PID %d\n", pid);
        strncpy(container_id, "eabrir", size);
        return container_id;
    }

    // Leer el contenido del archivo
    bytes_read = kernel_read(file, container_id, size - 1, &pos);
    if (bytes_read < 0) {
        printk(KERN_INFO "Error al leer el archivo cgroup para el PID %d\n", pid);
        strncpy(container_id, "eleer", size);
        filp_close(file, NULL);
        return container_id;
    }

    // Asegurarse de que la cadena esté terminada
    container_id[bytes_read] = '\0';

    // Cerrar el archivo
    filp_close(file, NULL);

    strncpy("aaaa", container_id, size);

    return container_id;
}

static int show_info(struct seq_file *m, void *v) {
    struct task_struct *task;

    struct sysinfo si;
    si_meminfo(&si);
    unsigned long total_mem = si.totalram * si.mem_unit / 1024;
    unsigned long free_mem = si.freeram * si.mem_unit / 1024;
    unsigned long used_mem = total_mem - free_mem;

    seq_printf(m, "{\n  \"Memoria_total\": %lu,\n  \"Memoria_libre\": %lu,\n  \"Memoria_uso\": %lu,\n",
               total_mem, free_mem, used_mem);

    unsigned int cpu_percent = get_cpu_usage();
    seq_printf(m, "  \"Uso_CPU\": \"%u%%\",\n", cpu_percent);

    seq_printf(m, "  \"procesos\": [\n");
    int first = 1;

    for_each_process(task) {
        if (strstr(task->comm, "stress")) {
            char container_id[42];
            get_container_id(task->pid, container_id, sizeof(container_id));

            if (!first) {
                seq_printf(m, ",\n");
            }
            first = 0;

            unsigned long mem_usage = task->mm ? get_mm_rss(task->mm) * 4 : 0;
            unsigned long cpu_time = task->se.sum_exec_runtime;

            seq_printf(m, "    { \"PID\": %d, \"Nombre\": \"%s\", \"Container_ID\": \"%s\", \"Porcentaje_Memoria\": \"%lu KB\", \"Porcentaje_CPU\": \"%lu ns\" }",
                task->pid, task->comm, container_id, mem_usage, cpu_time);
        }
    }

    seq_printf(m, "\n  ]\n}\n");
    return 0;
}

// Función para calcular el uso de CPU en porcentaje
unsigned int get_cpu_usage(void) {
    unsigned long total_jiffies = 0, idle_jiffies = 0, used_jiffies;
    unsigned int cpu_usage = 0;

    int i;
    for (i = 0; i < num_online_cpus(); i++) {
        total_jiffies += kcpustat_cpu(i).cpustat[CPUTIME_USER] +
                         kcpustat_cpu(i).cpustat[CPUTIME_SYSTEM] +
                         kcpustat_cpu(i).cpustat[CPUTIME_IDLE] +
                         kcpustat_cpu(i).cpustat[CPUTIME_IOWAIT];

        idle_jiffies += kcpustat_cpu(i).cpustat[CPUTIME_IDLE];
    }

    if (total_jiffies > 0) {
        used_jiffies = total_jiffies - idle_jiffies;
        cpu_usage = (used_jiffies * 100) / total_jiffies;
    }

    return cpu_usage;
}

static int open_proc(struct inode *inode, struct file *file) {
    return single_open(file, show_info, NULL);
}

static const struct proc_ops proc_fops = {
    .proc_open = open_proc,
    .proc_read = seq_read,
    .proc_lseek = seq_lseek,
    .proc_release = single_release,
};

// Cargar módulo
static int __init kernel_mod_init(void) {
    proc_create(PROC_FILENAME, 0, NULL, &proc_fops);
    printk(KERN_INFO "Módulo cargado: /proc/%s\n", PROC_FILENAME);
    return 0;
}

// Descargar módulo
static void __exit kernel_mod_exit(void) {
    remove_proc_entry(PROC_FILENAME, NULL);
    printk(KERN_INFO "Módulo descargado: /proc/%s\n", PROC_FILENAME);
}

module_init(kernel_mod_init);
module_exit(kernel_mod_exit);
MODULE_LICENSE("GPL");
MODULE_AUTHOR("Gonzalo");
MODULE_DESCRIPTION("Módulo de kernel para mostrar información del sistema");
MODULE_VERSION("1.1");
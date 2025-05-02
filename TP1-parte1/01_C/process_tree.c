#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <sys/wait.h>
#include <string.h>
#include <linux/prctl.h>
#include <sys/prctl.h>

#define PROCESS_NAME_MAX_LENGTH 16
#define SLEEP_DURATION_SECONDS 30

// Number of children per parent process
#define NUM_CHILDREN_A 1  // B
#define NUM_CHILDREN_B 2  // C and D
#define NUM_CHILDREN_C 1  // E
#define NUM_CHILDREN_D 2  // F and G
#define NUM_CHILDREN_E 2  // H and I
#define NUM_CHILDREN_LEAF 0  // Leaf nodes: F, G, H, I

void handleChildError(const char *childName);
void setProcessName(const char *name);
void sleepAndWait(int numChildren);
void createLeaf(const char *leafName);

int main() {
  pid_t pid = fork();

  if (pid < 0) {
    handleChildError("B");
    return EXIT_FAILURE;
  }

  if (pid > 0) {
    setProcessName("A");
    sleepAndWait(NUM_CHILDREN_A);
    return EXIT_SUCCESS;
  }

  // Process B
  setProcessName("B");

  pid = fork();
  if (pid < 0) {
    handleChildError("C");
    return EXIT_FAILURE;
  }

  if (pid == 0) {
    // Process C
    setProcessName("C");

    pid = fork();
    if (pid < 0) {
      handleChildError("E");
      return EXIT_FAILURE;
    }

    if (pid == 0) {
      // Process E
      setProcessName("E");

      createLeaf("H");
      createLeaf("I");

      sleepAndWait(NUM_CHILDREN_E);
      return EXIT_SUCCESS;
    }

    sleepAndWait(NUM_CHILDREN_C);
    return EXIT_SUCCESS;
  }

  pid = fork();
  if (pid < 0) {
    handleChildError("D");
    return EXIT_FAILURE;
  }

  if (pid == 0) {
    // Process D
    setProcessName("D");

    createLeaf("F");
    createLeaf("G");

    sleepAndWait(NUM_CHILDREN_D);
    return EXIT_SUCCESS;
  }

  sleepAndWait(NUM_CHILDREN_B);
  return EXIT_SUCCESS;
}

void handleChildError(const char *childName) {
  fprintf(stderr, "Error creating process %s\n", childName);
}

void setProcessName(const char *name) {
  char buffer[PROCESS_NAME_MAX_LENGTH];
  strncpy(buffer, name, PROCESS_NAME_MAX_LENGTH);
  prctl(PR_SET_NAME, buffer, 0, 0, 0);
}

void sleepAndWait(int numChildren) {
  sleep(SLEEP_DURATION_SECONDS);
  for (int i = 0; i < numChildren; ++i) {
    wait(NULL);
  }
}

void createLeaf(const char *leafName) {
  pid_t pid = fork();

  if (pid < 0) {
    handleChildError(leafName);
    exit(EXIT_FAILURE);
  }

  if (pid == 0) {
    setProcessName(leafName);
    sleepAndWait(NUM_CHILDREN_LEAF);
    exit(EXIT_SUCCESS);
  }
}

#ifndef PIGLET_H
#define PIGLET_H


int CreateContext(void);
int DestroyContext(void);
void MakeCurrent(void);
void SwapBuffers(void);

int GetDisplayWidth(void);
int GetDisplayHeight(void);

void* GetProcAddress(const char *name);

#endif //PIGLET_H


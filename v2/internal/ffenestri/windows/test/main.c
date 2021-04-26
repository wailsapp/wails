#include <stdio.h>
#include "../../ffenestri_windows.h"

int main() {
    struct Application *app = NewApplication("Wails ❤️ Unicode", 800, 600, 1, 1, 0, 0, 1, 0);
    SetMinWindowSize(app, 100, 100);
    SetMaxWindowSize(app, 800, 800);
    Run(app,0, NULL);
    return 0;
}

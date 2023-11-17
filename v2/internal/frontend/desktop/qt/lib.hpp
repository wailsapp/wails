#pragma once

#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif
extern void appExited(int code);

/* Application */
void *Application_new(char *app_name);
void Application_exec(void *app_ptr);
void Application_quit(void *app_ptr);
char *Application_get_screens(void *app_ptr); // Returns a json blob
/* End Application */

/* Window */
typedef struct Window {
  void *window;
  void *window_layout;
  void *web_engine_view;
} Window;

typedef struct Point {
  int x;
  int y;
} Point;

typedef struct RGBA {
  int r;
  int g;
  int b;
  int a;
} RGBA;

Window *Window_new(void *app_ptr, char *start_url);
void Window_set_title(void *win_ptr, char *title);
void Window_resize(void *win_ptr, int width, int height);
void Window_set_minimum_size(void *win_ptr, int width, int height);
void Window_set_maximum_size(void *win_ptr, int width, int height);
void Window_set_background_color(void *win_ptr, RGBA color);
Point Window_get_size(void *win_ptr);
void Window_hide(void *win_ptr);
int Window_get_flags(void *win_ptr);
void Window_set_flag(void *win_ptr, int flag, bool on);
void Window_show(void *win_ptr);
int Window_get_state(void *win_ptr);
void Window_fullscreen(void *win_ptr);
void Window_maximize(void *win_ptr);
void Window_close(void *win_ptr);
void Window_center(void *win_ptr);
void Window_center(void *win_ptr);
void Window_unminimize(void *win_ptr);
Point Window_get_position(void *win_ptr);
void Window_set_position(void *win_ptr, Point position);
const char *Window_run_message_dialog(void *win_ptr, int dialog_type, char *title, char *message);
const char *Window_open_file_dialog(void *win_ptr, bool isDirectory, bool isMultiple, bool isSave, char *dialog_options);
/* End Window */

/* WebEngineView */
void WebEngineView_load_url(void *web_engine_ptr, char *url);
void WebEngineView_reload(void *web_engine_ptr);
void WebEngineView_run_js(void *web_engine_ptr, char *script);
void WebEngineView_print_page(void *web_engine_ptr);
/* End WebEngineView */

/* Clipboard */
const char *Clipboard_get_text(void *app_ptr);
void Clipboard_set_text(void *app_ptr, char *text);
/* End Clipboard */

void fix_signal(int signum);
void install_signal_handlers();

void cfree(void* ptr);

#ifdef __cplusplus
}
#endif

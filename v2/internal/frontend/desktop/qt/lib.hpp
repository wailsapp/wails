#pragma once

#ifdef __cplusplus
extern "C" {
#endif
extern void appExited(int code);

/* Application */
void *Application_new(char *app_name);
void Application_exec(void *app_ptr);
/* End Application */

/* Window */
typedef struct Window {
  void *window;
  void *window_layout;
  void *web_engine_view;
} Window;

Window *Window_new(void *app_ptr);
void Window_set_title(void *win_ptr, char *title);
void Window_resize(void *win_ptr, int width, int height);
void Window_set_minimum_size(void *win_ptr, int width, int height);
/* End Window */

/* WebEngineView */
void WebEngineView_load_url(void *web_engine_ptr, char *url);
/* End WebEngineView */

void fix_signal(int signum);
void install_signal_handlers();

#ifdef __cplusplus
}
#endif

#include "lib.hpp"

#include <QApplication>
#include <QLabel>
#include <QMetaObject>
#include <QTimer>
#include <QUrl>
#include <QVBoxLayout>
#include <QWebEngineView>
#include <QtWidgets>
#include <condition_variable>
#include <iostream>
#include <memory>
#include <mutex>
#include <thread>

#include <errno.h>
#include <math.h>
#include <signal.h>
#include <stdio.h>
#include <string.h>

#include "util.hpp"

/* Application */

void *Application_new(char *app_name) {
  QCoreApplication::setAttribute(Qt::AA_EnableHighDpiScaling);

  auto queue = std::make_unique<SafeQueue<QApplication *>>();

  // The QApplication must be started on its own (non QThread) thread,
  // See
  // https://forum.qt.io/topic/124878/running-qapplication-exec-from-another-thread-qcoreapplication-qguiapplication
  auto appThread = new std::thread([&]() {
    int numArgs = 1;
    char *args[] = {app_name};
    auto app = new QApplication(numArgs, args);
    queue->enqueue(app);
    appExited(app->exec());
  });

  auto qtApp = queue->dequeue();

  // Ensure that app has started executing before returning.
  QTimer::singleShot(0, qtApp, [&]() { queue->enqueue(nullptr); });
  queue->dequeue();
  return qtApp;
}

void Application_quit(void *app_ptr) {
  auto app = static_cast<QApplication *>(app_ptr);
  runOnAppThread(app, [=]() { app->quit(); });
}

/* End Application */

/* Window */

Window *Window_new(void *app_ptr, char *start_url) {
  auto app = static_cast<QApplication *>(app_ptr);

  Window *win;
  runOnAppThread(
      app,
      [=]() -> Window * {
        auto w = new QWidget();
        w->resize(800, 600);
        w->setMinimumSize(320, 240);

        auto layout = new QVBoxLayout(w);
        layout->setContentsMargins(0, 0, 0, 0);
        layout->setSpacing(0);

        auto view = new QWebEngineView(w);
        layout->addWidget(view);
        view->load(QUrl(start_url));

        w->show();

        return new Window{
            .window = w,
            .window_layout = layout,
            .web_engine_view = view,
        };
      },
      &win);

  return win;
}

void Window_set_title(void *win_ptr, char *title) {
  auto win = static_cast<QWidget *>(win_ptr);
  QString qtitle(title);
  runOnAppThread(win, [=]() { win->setWindowTitle(qtitle); });
}

void Window_set_minimum_size(void *win_ptr, int width, int height) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->setMinimumSize(width, height); });
}

void Window_resize(void *win_ptr, int width, int height) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->resize(width, height); });
}

void Window_hide(void *win_ptr) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->showMinimized(); });
}

void Window_show(void *win_ptr) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->showNormal(); });
}

void Window_fullscreen(void *win_ptr) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->setWindowState(win->windowState() ^ Qt::WindowFullScreen); });
}

void Window_maximize(void *win_ptr) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->setWindowState(win->windowState() ^ Qt::WindowMaximized); });
}

void Window_close(void *win_ptr) {
  auto win = static_cast<QWidget *>(win_ptr);
  runOnAppThread(win, [=]() { win->close(); });
}

/* End Window */

/* WebEngineView */

void WebEngineView_load_url(void *web_engine_ptr, char *url) {
  auto eng = static_cast<QWebEngineView *>(web_engine_ptr);
  runOnAppThread(eng, [=]() { eng->load(QUrl(url)); });
}

void WebEngineView_reload(void *web_engine_ptr) {
  auto eng = static_cast<QWebEngineView *>(web_engine_ptr);
  runOnAppThread(eng, [=]() { eng->reload(); });
}

/* End WebEngineView */

/* Misc */

// CREDIT: https://github.com/rainycape/magick
void fix_signal(int signum) {
  struct sigaction st;

  if (sigaction(signum, NULL, &st) < 0) {
    goto fix_signal_error;
  }
  st.sa_flags |= SA_ONSTACK;
  if (sigaction(signum, &st, NULL) < 0) {
    goto fix_signal_error;
  }
  return;
fix_signal_error:
  fprintf(stderr,
          "error fixing handler for signal %d, please "
          "report this issue to "
          "https://github.com/wailsapp/wails: %s\n",
          signum, strerror(errno));
}

void install_signal_handlers() {
#if defined(SIGCHLD)
  fix_signal(SIGCHLD);
#endif
#if defined(SIGHUP)
  fix_signal(SIGHUP);
#endif
#if defined(SIGINT)
  fix_signal(SIGINT);
#endif
#if defined(SIGQUIT)
  fix_signal(SIGQUIT);
#endif
#if defined(SIGABRT)
  fix_signal(SIGABRT);
#endif
#if defined(SIGFPE)
  fix_signal(SIGFPE);
#endif
#if defined(SIGTERM)
  fix_signal(SIGTERM);
#endif
#if defined(SIGBUS)
  fix_signal(SIGBUS);
#endif
#if defined(SIGSEGV)
  fix_signal(SIGSEGV);
#endif
#if defined(SIGXCPU)
  fix_signal(SIGXCPU);
#endif
#if defined(SIGXFSZ)
  fix_signal(SIGXFSZ);
#endif
}

void bye(void* ptr) {
    free(ptr);
}

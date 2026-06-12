#ifndef WAILS_MOBILE_FEATURES_IOS_H
#define WAILS_MOBILE_FEATURES_IOS_H

#include <stdbool.h>

// Phase A — one-way actions
void ios_share(const char* json);
void ios_open_url(const char* url);
void ios_set_keep_awake(bool enabled);
void ios_set_torch(bool enabled);

// Phase B — state / query
const char* ios_safe_area_json(void);
void ios_set_brightness(double value);     // 0.0 - 1.0
double ios_get_brightness(void);
const char* ios_app_info_json(void);
void ios_set_orientation(const char* mode); // "portrait" | "landscape" | "auto"
const char* ios_get_orientation(void);
void ios_set_status_bar(const char* json);  // {"style":"light|dark|default","hidden":bool}

// Phase C — async results / permissions
void ios_biometric_authenticate(const char* reason);
void ios_post_notification(const char* json); // {"title":"","body":"","delay":seconds}
void ios_secure_set(const char* key, const char* value);
const char* ios_secure_get(const char* key);
void ios_secure_delete(const char* key);

#endif

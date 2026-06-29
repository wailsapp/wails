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
void ios_notifications_init(void);             // register the UNUserNotificationCenter delegate (call once, at launch)
void ios_secure_set(const char* key, const char* value);
const char* ios_secure_get(const char* key);
void ios_secure_delete(const char* key);

// Phase D — sensors & hardware
void ios_haptic(const char* type);          // impact-light|impact-medium|impact-heavy|success|warning|error|selection
void ios_get_location(void);                // async → "common:location" {lat,lng,accuracy,error}
void ios_set_motion(bool enabled);          // accelerometer stream → "common:motion" {x,y,z}
void ios_set_proximity(bool enabled);       // proximity sensor → "common:proximity" {near}
void ios_speak(const char* text);           // text-to-speech
void ios_stop_speak(void);
const char* ios_storage_json(void);         // {"free":bytes,"total":bytes}
const char* ios_storage_path(void);         // absolute path to the app's Application Support directory
const char* ios_power_json(void);           // {"level":0-1,"charging":bool,"lowPower":bool}
const char* ios_network_json(void);         // {"connected":bool,"type":"wifi|cellular|none"}
void ios_set_keyboard_watch(bool enabled);  // keyboard insets → "common:keyboard" {visible,height}
void ios_set_screen_protect(bool enabled);  // screenshot/recording detection → "common:screenCapture"

// Phase E — camera & background
void ios_capture_photo(void);               // camera → "common:capture" {type:"photo",path,size,thumb}
void ios_capture_video(void);               // camera → "common:capture" {type:"video",path,size}
void ios_begin_background_task(int seconds); // open a background-task window → "ios:backgroundTask"
void ios_end_background_task(void);

#endif

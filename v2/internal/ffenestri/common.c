//
// Created by Lea Anthony on 6/1/21.
//

#include "common.h"

// Credit: https://stackoverflow.com/a/8465083
char* concat(const char *string1, const char *string2)
{
    const size_t len1 = strlen(string1);
    const size_t len2 = strlen(string2);
    char *result = malloc(len1 + len2 + 1);
    strcpy(result, string1);
    memcpy(result + len1, string2, len2 + 1);
    return result;
}

// 10k is more than enough for a log message
#define MAXMESSAGE 1024*10
char abortbuffer[MAXMESSAGE];

void ABORT(const char *message, ...) {
    const char *temp = concat("FATAL: ", message);
    va_list args;
    va_start(args, message);
    vsnprintf(abortbuffer, MAXMESSAGE, temp, args);
    printf("%s\n", &abortbuffer[0]);
    MEMFREE(temp);
    va_end(args);
    exit(1);
}

int freeHashmapItem(void *const context, struct hashmap_element_s *const e) {
    free(e->data);
    return -1;
}

const char* getJSONString(JsonNode *item, const char* key) {
    // Get key
    JsonNode *node = json_find_member(item, key);
    const char *result = "";
    if ( node != NULL && node->tag == JSON_STRING) {
        result = node->string_;
    }
    return result;
}

void ABORT_JSON(JsonNode *node, const char* key) {
    ABORT("Unable to read required key '%s' from JSON: %s\n", key, json_encode(node));
}

const char* mustJSONString(JsonNode *node, const char* key) {
    const char* result = getJSONString(node, key);
    if ( result == NULL ) {
        ABORT_JSON(node, key);
    }
    return result;
}
JsonNode* mustJSONObject(JsonNode *node, const char* key) {
    struct JsonNode* result = getJSONObject(node, key);
    if ( result == NULL ) {
        ABORT_JSON(node, key);
    }
    return result;
}

JsonNode* getJSONObject(JsonNode* node, const char* key) {
    return json_find_member(node, key);
}

bool getJSONBool(JsonNode *item, const char* key, bool *result) {
    JsonNode *node = json_find_member(item, key);
    if ( node != NULL && node->tag == JSON_BOOL) {
        *result = node->bool_;
        return true;
    }
    return false;
}

bool getJSONInt(JsonNode *item, const char* key, int *result) {
    JsonNode *node = json_find_member(item, key);
    if ( node != NULL && node->tag == JSON_NUMBER) {
        *result = (int) node->number_;
        return true;
    }
    return false;
}

JsonNode* mustParseJSON(const char* JSON) {
    JsonNode* parsedUpdate = json_decode(JSON);
    if ( parsedUpdate == NULL ) {
        ABORT("Unable to decode JSON: %s\n", JSON);
    }
    return parsedUpdate;
}
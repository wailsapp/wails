//
// Created by Lea Anthony on 6/1/21.
//

#ifndef COMMON_H
#define COMMON_H

#include <stdio.h>
#include <stdarg.h>
#include "string.h"
#include "hashmap.h"
#include "vec.h"
#include "json.h"

#define STREQ(a,b) strcmp(a, b) == 0
#define STREMPTY(string) strlen(string) == 0
#define STRCOPY(a) concat(a, "")
#define STR_HAS_CHARS(input) input != NULL && strlen(input) > 0
#define MEMFREE(input) free((void*)input); input = NULL;
#define FREE_AND_SET(variable, value) if( variable != NULL ) { MEMFREE(variable); } variable = value

// Credit: https://stackoverflow.com/a/8465083
char* concat(const char *string1, const char *string2);
void ABORT(const char *message, ...);
int freeHashmapItem(void *const context, struct hashmap_element_s *const e);
const char* getJSONString(JsonNode *item, const char* key);
const char* mustJSONString(JsonNode *node, const char* key);
JsonNode* getJSONObject(JsonNode* node, const char* key);
JsonNode* mustJSONObject(JsonNode *node, const char* key);

bool getJSONBool(JsonNode *item, const char* key, bool *result);
bool getJSONInt(JsonNode *item, const char* key, int *result);

JsonNode* mustParseJSON(const char* JSON);

#endif //ASSETS_C_COMMON_H

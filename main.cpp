#include <iostream>
#include <regex>
#include <curl/curl.h>

#include "json.hpp"

#define URL "http://reddit.com/r/hockey.json?limit=100"
#define UA "hockeygfy v2.0 by u/aggrolite; twitter=@hockeygfy; website=hockeygfy.com; src=github.com/aggrolite/hockeygfy"


using namespace::std;
using namespace::nlohmann;

//const bool debug = true;

// Write response body to our provided pointer.
size_t writeBody(char* content, size_t size, size_t nmemb, void* body) {
    size_t realsize = size * nmemb;
    string* b = static_cast<string*>(body);
    char* c = static_cast<char*>(content);

    // writeCallback() could be called more than once.
    // Append response body and parse as JSON later.
    b->append(c, realsize);

    return realsize;
}

void setOptions(CURL* c, string* b) {
    // Set URL.
    curl_easy_setopt(c, CURLOPT_URL, URL);

    // Set User Agent.
    curl_easy_setopt(c, CURLOPT_USERAGENT, UA);

    // Follow redirect.
    curl_easy_setopt(c, CURLOPT_FOLLOWLOCATION, 1L);

    // Write content to given pointer.
    curl_easy_setopt(c, CURLOPT_WRITEDATA, b);

    // Write content using defined callback.
    curl_easy_setopt(c, CURLOPT_WRITEFUNCTION, writeBody);
}

// Ignore non-gfycat links.
bool filterLinks(int depth, json::parse_event_t event, json& parsed) {
    bool isObject = (event == json::parse_event_t::object_end);
    bool hasData = parsed.count("data") > 0;

    if (isObject and hasData and parsed["data"].count("is_self")) {
        auto d = parsed["data"];
        string url = d["url"];
        regex re("^https?://(?:www\\.)?gfycat\\.com", regex::icase);

        if (!d["is_self"] and regex_search(url, re)) {
            return true;
        }
        return false;
    }
    return true;
}

int main(int argc, char **argv) {

    // Create curl handle.
    CURL* conn = curl_easy_init();
    if (!conn) {
        cerr << "Failed to create curl handle.\n";
        return 1;
    }

    // Create buffer for JSON content.
    string buffer;

    // Set curl handle options.
    setOptions(conn, &buffer);

    CURLcode res = curl_easy_perform(conn);

    if (res != CURLE_OK) {
        cerr << "Error: " << curl_easy_strerror(res) << "\n";
        return 1;
    }
    curl_easy_cleanup(conn);

    // Parse JSON and use callback to filter out non-gfycat links.
    json j = json::parse(buffer, static_cast<json::parser_callback_t>(filterLinks));

    // Print URLs collected.
    for (auto& el : j["data"]["children"]) {
        cout << el["data"]["url"] << "\n";
    }
    return 0;
}

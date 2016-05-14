#define CURL_MAX_WRITE_SIZE 5000000

#include <iostream>
#include <curl/curl.h>

#include "json.hpp"

#define URL "http://reddit.com/r/hockey.json"
#define UA "hockeygfy v2.0 by u/aggrolite; twitter=@hockeygfy; website=hockeygfy.com; src=github.com/aggrolite/hockeygfy"


using namespace::std;
using namespace::nlohmann;

size_t writeCallback(char* content, size_t size, size_t nmemb, void* userp) {
    size_t realsize = size * nmemb;
    ((std::string*)userp)->append((char*)content, size * nmemb);
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
    curl_easy_setopt(c, CURLOPT_WRITEFUNCTION, writeCallback);
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

    json j = json::parse(buffer);

    cout << j["kind"] << "\n";

    return 0;
}

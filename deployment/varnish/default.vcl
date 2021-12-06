vcl 4.1;

import directors;
import std;

backend default {
    .host = "varnish-cache:80";
}

sub vcl_init {
    # Called when VCL is loaded, before any requests pass through it.
    # Typically used to initialize VMODs.

    new vdir = directors.round_robin();
    vdir.add_backend(default);
}

sub vcl_pipe {
    # Called upon entering pipe mode.
    # In this mode, the request is passed on to the backend, and any further data from both the client
    # and backend is passed on unaltered until either end closes the connection. Basically, Varnish will
    # degrade into a simple TCP proxy, shuffling bytes back and forth. For a connection in pipe mode,
    # no other VCL subroutine will ever get called after vcl_pipe.

    # Note that only the first request to the backend will have
    # X-Forwarded-For set.  If you use X-Forwarded-For and want to
    # have it set for all requests, make sure to have:
    # set bereq.http.connection = "close";
    # here.  It is not set by default as it might break some broken web
    # applications, like IIS with NTLM authentication.

    # set bereq.http.Connection = "Close";

    # Implementing websocket support (https://www.varnish-cache.org/docs/4.0/users-guide/vcl-example-websockets.html)
    if (req.http.upgrade) {
        set bereq.http.upgrade = req.http.upgrade;
    }

    return (pipe);
}

sub vcl_pass {
    # Called upon entering pass mode. In this mode, the request is passed on to the backend, and the
    # backend's response is passed on to the client, but is not entered into the cache. Subsequent
    # requests submitted over the same client connection are handled normally.

    # return (pass);
}

# The data on which the hashing will take place
sub vcl_hash {
    # Called after vcl_recv to create a hash value for the request. This is used as a key
    # to look up the object in Varnish.

    hash_data(req.url);

    if (req.http.host) {
        hash_data(req.http.host);
    } else {
        hash_data(server.ip);
    }

    # hash cookies for requests that have them
    if (req.http.Cookie) {
        hash_data(req.http.Cookie);
    }

    # Cache the HTTP vs HTTPs separately
    if (req.http.X-Forwarded-Proto) {
        hash_data(req.http.X-Forwarded-Proto);
    }
}

sub vcl_hit {
    # Called when a cache lookup is successful.

    if (obj.ttl >= 0s) {
        # A pure unadultered hit, deliver it
        return (deliver);
    }

    # https://www.varnish-cache.org/docs/trunk/users-guide/vcl-grace.html
    # When several clients are requesting the same page Varnish will send one request to the backend and place the others
    # on hold while fetching one copy from the backend. In some products this is called request coalescing and Varnish does
    # this automatically.
    # If you are serving thousands of hits per second the queue of waiting requests can get huge. There are two potential
    # problems - one is a thundering herd problem - suddenly releasing a thousand threads to serve content might send the
    # load sky high. Secondly - nobody likes to wait. To deal with this we can instruct Varnish to keep the objects in cache
    # beyond their TTL and to serve the waiting requests somewhat stale content.

    # if (!std.healthy(req.backend_hint) && (obj.ttl + obj.grace > 0s)) {
    #   return (deliver);
    # } else {
    #   return (miss);
    # }

    # We have no fresh fish. Lets look at the stale ones.
    if (std.healthy(req.backend_hint)) {
        # Backend is healthy. Limit age to 10s.
        if (obj.ttl + 10s > 0s) {
            # set req.http.grace = "normal(limited)";
            return (deliver);
        }
    } else {
        # backend is sick - use full grace
        if (obj.ttl + obj.grace > 0s) {
            # set req.http.grace = "full";
            return (deliver);
        }
    }
}

sub vcl_miss {
    # Called after a cache lookup if the requested document was not found in the cache. Its purpose
    # is to decide whether or not to attempt to retrieve the document from the backend, and which
    # backend to use.

    return (fetch);
}

# The routine when we deliver the HTTP request to the user
# Last chance to modify headers that are sent to the client
sub vcl_deliver {
    # Called before a cached object is delivered to the client.

    # Add debug header to see if it's a HIT/MISS and the number of hits, disable when not needed
    if (obj.hits > 0) {
        set resp.http.X-Cache = "HIT";
    } else {
        set resp.http.X-Cache = "MISS";
    }

    # Please note that obj.hits behaviour changed in 4.0, now it counts per objecthead, not per object
    # and obj.hits may not be reset in some cases where bans are in use. See bug 1492 for details.
    # So take hits with a grain of salt
    set resp.http.X-Cache-Hits = obj.hits;

    # Remove some headers: PHP version
    unset resp.http.X-Powered-By;

    # Remove some headers: Apache version & OS
    unset resp.http.Server;
    unset resp.http.X-Drupal-Cache;
    unset resp.http.X-Varnish;
    unset resp.http.Via;
    unset resp.http.Link;
    unset resp.http.X-Generator;

    return (deliver);
}

sub vcl_purge {
    # Only handle actual PURGE HTTP methods, everything else is discarded
    if (req.method == "PURGE") {
    # restart request
        set req.http.X-Purge = "Yes";
        return(restart);
    }
}

sub vcl_synth {
    if (resp.status == 720) {
        # We use this special error status 720 to force redirects with 301 (permanent) redirects
        # To use this, call the following from anywhere in vcl_recv: return (synth(720, "http://host/new.html"));
        set resp.http.Location = resp.reason;
        set resp.status = 301;
        return (deliver);
    } elseif (resp.status == 721) {
        # And we use error status 721 to force redirects with a 302 (temporary) redirect
        # To use this, call the following from anywhere in vcl_recv: return (synth(720, "http://host/new.html"));
        set resp.http.Location = resp.reason;
        set resp.status = 302;
        return (deliver);
    }
    return (deliver);
}

sub vcl_fini {
    # Called when VCL is discarded only after all requests have exited the VCL.
    # Typically used to clean up VMODs.

    return (ok);
}

sub vcl_recv {
    # Called at the beginning of a request, after the complete request has been received and parsed.
    # Its purpose is to decide whether or not to serve the request, how to do it, and, if applicable,
    # which backend to use.
    # also used to modify the request

    set req.backend_hint = vdir.backend(); # send all traffic to the vdir director

    # Normalize the header if it exists, remove the port (in case you're testing this on various TCP ports)
    if (req.http.Host) {
        set req.http.Host = regsub(req.http.Host, ":[0-9]+", "");
    }

    # Remove the proxy header (see https://httpoxy.org/#mitigate-varnish)
    unset req.http.proxy;

    # Normalize the query arguments
    set req.url = std.querysort(req.url);

    if (req.method == "PURGE") {
        return (purge);
    }

    # Only deal with "normal" types
    if (req.method != "GET" && req.method != "HEAD" && req.method != "PUT" &&
            req.method != "POST" && req.method != "TRACE" && req.method != "OPTIONS" &&
            req.method != "PATCH" && req.method != "DELETE") {
        return (pipe);
    }

    # Implementing websocket support (https://www.varnish-cache.org/docs/4.0/users-guide/vcl-example-websockets.html)
    if (req.http.Upgrade ~ "(?i)websocket") {
        return (pipe);
    }

    # Only cache GET or HEAD requests. This makes sure the POST requests are always passed.
    if (req.method != "GET" && req.method != "HEAD") {
        return (pass);
    }


    # Backend declaration
    if (req.http.host == "varnish-cache") {
        set req.backend_hint = default;
        if (! req.url ~ "/") {
            return(pass);
        }
    }
}

sub vcl_backend_response {
    # Happens after we have read the response headers from the backend.
    # Here you clean the response headers, removing silly Set-Cookie headers
    # and other mistakes your backend does.

    set beresp.ttl = 5m;
    set beresp.grace = 30m;

    # Don't cache 50x responses
    if (beresp.status == 500 || beresp.status == 502 || beresp.status == 503 || beresp.status == 504) {
        return (abandon);
    }

    return (deliver);
}

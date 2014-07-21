# Hecate 
### WORK IN PROGRESS THAT README MIGHT NOT REFLECT THE WAY IT'S WORKING FOR NOW

A small script to spit out nginx configuration from mesos or kubernetes


# Idea

Hecate needs to update a *nginx* configuration, by interrogating
different api :

    1 - Mesos based on app name
    2 - (NOT YET) Kubernetes api based on label

# Command line


### Get configuration endpoints

Get json endpoints for a named application

    $ hecate --scheduler=kubernetes \
             --host=http://0.0.0.0:8080 \
             --name=envspitter \
             --endpoints

Response :

    200 OK
    {
        "name": "envspitter",
        "endpoints": 
        [
            {
                "host": "host1.com",
                "ports":
                [ 
                    "52345", 
                    "23452" 
                ]
            },
            {
                "host": "host2.com",
                "ports":
                [ 
                    "52335", 
                    "22252" 
                ]
            }
        ]
    }

Get nginx endpoints for a named application

    $ hecate --adapter=kubernetes \
             --host=http://0.0.0.0:8080 \
             --name=envspitter \
             --endpoints=nginx

Response :

    200 OK
    upstream batch_backends {
        server host1.com:52345 weight=100 max_fails=3 fail_timeout=30s ;
        server host1.com:23452 weight=100 max_fails=3 fail_timeout=30s ;
        server host2.com:52335 weight=100 max_fails=3 fail_timeout=30s ;
        server host2.com:22252 weight=100 max_fails=3 fail_timeout=30s ;
    }
    
### Get configuration endpoints

# NOT YET IMPLEMENTED
Create vhost for an application

    $ hecate --adapter=mesos \
             --host=http://mesoshost.com \
             --port=8080 \
             --name=envspitter \
             --hostname=envspitter.com \
             --host_port=8080 \
             --ssl=true \ 
             --ssl_certificate=/path/host.com_concat.crt \
             --ssl_key=/path/host.com.key \

Response

    server {
        listen 443;
        server_name envspitter.com;

        ssl on;
        ssl_certificate     /path/host.com_concat.crt;
        ssl_certificate_key /path/host.com.key;
        ssl_session_cache shared:SSL:20m;
        ssl_session_timeout 10m;
        ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
        ssl_prefer_server_ciphers on;
        ssl_ciphers ECDH+AESGCM:ECDH+AES256:ECDH+AES128:DH+3DES:!ADH:!AECDH:!MD5;
        resolver 8.8.8.8 8.8.4.4;

        access_log /var/log/nginx/envspitter.access.log ;

        error_page 404 500 503 400 /notfound.html ;
        location / {
            proxy_pass http://envspitter/;
            proxy_intercept_errors on;
        }
        location = /notfound.html {
           root /var/www/shared/off/ ;
        }
    }
